package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/lfkeitel/asml-sim/pkg/token"
)

var (
	ASMLHeader = []byte("ASML")
)

type bitParts uint8

const (
	fullBits   bitParts = 0
	higherBits bitParts = 1
	lowerBits  bitParts = 2
)

type LabelReplace struct {
	l      string
	offset uint16
	part   bitParts
}

type Lexer struct {
	in              io.ReadSeeker
	labels          map[string]uint16       // Label definitions
	labelPlaces     map[uint16]LabelReplace // Memory locations that need labels
	currMemLocation uint16
	linenum         int
}

func New(in io.ReadSeeker) *Lexer {
	return &Lexer{
		in: in,
	}
}

func (l *Lexer) checkBinaryFile() []byte {
	// Read in a compiled ASML file
	header := make([]byte, 4)
	n, err := l.in.Read(header)
	if err != nil {
		fmt.Printf("Error reading file header: %s\n", err)
		os.Exit(1)
	}
	if n < 4 {
		fmt.Println("Invalid file")
		os.Exit(1)
	}

	if bytes.Equal(header, ASMLHeader) {
		var buf bytes.Buffer
		io.Copy(&buf, l.in)
		return buf.Bytes()
	}
	return nil
}

func (l *Lexer) Lex() []uint8 {
	bin := l.checkBinaryFile()
	if bin != nil {
		return bin
	}

	// Rewind file to read in as source
	l.in.Seek(0, 0)
	reader := bufio.NewReader(l.in)
	code := l.readCode(reader)
	l.linkCode(code)

	return code
}

type LabelMap map[string]uint16
type LabelLinkMap map[uint16]LabelReplace

func (m LabelLinkMap) FindOffsets(label string) []uint16 {
	links := make([]uint16, 0, 5)
	for l, lr := range m {
		if lr.l == label {
			links = append(links, l)
		}
	}
	return links
}

func (l *Lexer) LexNoLink() ([]uint8, LabelMap, LabelLinkMap) {
	bin := l.checkBinaryFile()
	if bin != nil {
		return bin, nil, nil
	}

	// Rewind file to read in as source
	l.in.Seek(0, 0)
	reader := bufio.NewReader(l.in)
	return l.readCode(reader), l.labels, l.labelPlaces
}

func (l *Lexer) readCode(reader *bufio.Reader) []uint8 {
	var code []uint8
	l.linenum = 0
	l.labels = make(LabelMap)          // Label definitions
	l.labelPlaces = make(LabelLinkMap) // Memory locations that need labels
	l.currMemLocation = 0

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		l.linenum++

		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == ';' { // comment/blank line
			continue
		}

		if line[0] == ':' { // label definition
			l.labels[string(line[1:])] = l.currMemLocation
			continue
		}

		if line[0] == '"' { // String, converted to bytes
			code = append(code, line[1:len(line)-1]...)
			continue
		}

		if line[0] == '\'' { // Single byte in ASCII
			code = append(code, line[1])
			continue
		}

		if line[0] == '@' { // Compiler directive
			if len(code) != 0 {
				fmt.Printf("Compiler directives must be before any code. Error in line %d\n", l.linenum)
				os.Exit(1)
			}

			if bytes.Equal([]byte("@runtime"), line) {
				code = append(code, runtime...)
				for n, o := range runtimeLabels {
					l.labels[n] = o
				}
				l.currMemLocation += uint16(len(code))
				l.labelPlaces[mainLabelLoc] = LabelReplace{
					l: "main",
				}
			}
			continue
		}

		instruction := bytes.Split(line, []byte{' '})
		opcode, valid := token.Opcodes[string(instruction[0])]
		if !valid { // Literal bytes
			for _, rawbyte := range instruction {
				raw, err := strconv.ParseUint(string(rawbyte), 0, 8)
				if err != nil {
					fmt.Printf("Invalid byte sequence on line %d: %v\n", l.linenum, err)
					os.Exit(1)
				}

				code = append(code, uint8(raw))
			}
			continue
		}

		switch opcode {
		case token.HALT, token.NOOP:
			code = append(code, opcode)
			l.currMemLocation++
		case token.PUSH, token.POP:
			reg := l.oneReg(instruction[1:])
			code = append(code, opcode, reg)
			l.currMemLocation += 2
		case token.ROT:
			reg, size := l.oneRegOneDigit(instruction[1:])
			code = append(code, opcode, reg, size)
			l.currMemLocation += 3
		case token.MOVR, token.STRR, token.LOADR:
			reg1, reg2 := l.twoRegisters(instruction[1:])
			code = append(code, opcode, reg1, reg2)
			l.currMemLocation += 3
		case token.ADD, token.OR, token.AND, token.XOR:
			reg1, reg2, reg3 := l.threeRegisters(instruction[1:])
			code = append(code, opcode, reg1, reg2, reg3)
			l.currMemLocation += 4
		case token.LOADI, token.LOADA, token.STRA, token.JMP:
			reg, b := l.oneRegTwoByte(instruction[1:])
			code = append(code, opcode, reg, uint8(b>>8), uint8(b))
			l.currMemLocation += 4
		case token.ADDI:
			reg1, reg2, b := l.twoRegOneByte(instruction[1:])
			code = append(code, opcode, reg1, reg2, b)
			l.currMemLocation += 4
		case token.JMPA, token.LDSP, token.LDSPI:
			b := l.twoByte(instruction[1:])
			code = append(code, opcode, uint8(b>>8), uint8(b))
			l.currMemLocation += 3
		default:
			fmt.Printf("Invalid opcode on line %d\n", l.linenum)
			os.Exit(1)
		}
	}

	return code
}

func (l *Lexer) linkCode(code []uint8) {
	for loc, label := range l.labelPlaces {
		memloc, exists := l.labels[label.l]
		if !exists {
			fmt.Printf("Label %s not defined\n", label.l)
			os.Exit(1)
		}

		switch label.part {
		case higherBits:
			code[loc] = uint8(memloc>>8) + uint8(label.offset>>8)
		case lowerBits:
			code[loc] = uint8(memloc) + uint8(label.offset)
		case fullBits:
			fallthrough
		default:
			code[loc] = uint8(memloc>>8) + uint8(label.offset>>8)
			code[loc+1] = uint8(memloc) + uint8(label.offset)
		}
	}
}

func (l *Lexer) oneReg(instruction [][]byte) uint8 {
	if len(instruction) < 1 {
		return 0
	}

	if instruction[0][0] != '%' {
		return 0
	}

	reg, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0
	}
	if reg > 15 {
		reg = 0
	}

	return uint8(reg)
}

func (l *Lexer) oneRegOneDigit(instruction [][]byte) (uint8, uint8) {
	if len(instruction) < 2 {
		return 0, 0
	}

	if instruction[0][0] != '%' {
		return 0, 0
	}

	reg, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0, 0
	}
	if reg > 15 {
		reg = 0
	}

	digit, err := strconv.ParseUint(string(instruction[1]), 16, 8)
	if err != nil {
		return 0, 0
	}
	if digit > 15 {
		digit = 0
	}
	return uint8(reg), uint8(digit)
}

func (l *Lexer) twoRegisters(instruction [][]byte) (uint8, uint8) {
	if len(instruction) < 2 {
		return 0, 0
	}

	if instruction[0][0] != '%' || instruction[1][0] != '%' {
		return 0, 0
	}

	reg1, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0, 0
	}
	if reg1 > 15 {
		reg1 = 0
	}

	reg2, err := strconv.ParseUint(string(instruction[1][1:]), 16, 8)
	if err != nil {
		return 0, 0
	}
	if reg2 > 15 {
		reg2 = 0
	}
	return uint8(reg1), uint8(reg2)
}

func (l *Lexer) threeRegisters(instruction [][]byte) (uint8, uint8, uint8) {
	if len(instruction) < 3 {
		return 0, 0, 0
	}

	if instruction[0][0] != '%' || instruction[1][0] != '%' || instruction[2][0] != '%' {
		return 0, 0, 0
	}

	reg1, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0, 0, 0
	}
	if reg1 > 15 {
		reg1 = 0
	}

	reg2, err := strconv.ParseUint(string(instruction[1][1:]), 16, 8)
	if err != nil {
		return 0, 0, 0
	}
	if reg2 > 15 {
		reg2 = 0
	}

	reg3, err := strconv.ParseUint(string(instruction[2][1:]), 16, 8)
	if err != nil {
		return 0, 0, 0
	}
	if reg3 > 15 {
		reg3 = 0
	}
	return uint8(reg1), uint8(reg2), uint8(reg3)
}

func (l *Lexer) oneRegTwoByte(instruction [][]byte) (uint8, uint16) {
	if len(instruction) < 2 {
		return 0, 0
	}

	if instruction[0][0] != '%' {
		return 0, 0
	}

	reg, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0, 0
	}
	if reg > 15 {
		reg = 0
	}

	var digit uint64
	if instruction[1][0] == '~' { // Label
		bits := fullBits
		label := instruction[1][1:]
		if label[0] == '^' {
			bits = higherBits
			label = label[1:]
		} else if label[0] == '`' {
			bits = lowerBits
			label = label[1:]
		}

		var offset uint16
		addIndex := bytes.Index(instruction[1], []byte{'+'})
		subIndex := bytes.Index(instruction[1], []byte{'-'})
		if addIndex > 0 || subIndex > 0 {
			ind := addIndex
			if subIndex > 0 {
				ind = subIndex
			}
			label = instruction[1][1:ind]
			offset64, err := strconv.ParseInt(string(instruction[1][ind+1:]), 0, 16)
			if err != nil {
				fmt.Printf("Invalid offset on line %d\n", l.linenum)
				os.Exit(1)
			}
			offset = uint16(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}
		if label[0] == '$' {
			digit = uint64(l.currMemLocation + offset)
		} else {
			l.labelPlaces[l.currMemLocation+2] = LabelReplace{
				l:      string(label),
				offset: offset,
				part:   bits,
			}
		}
	} else if instruction[1][0] == '\'' { // Literal byte character
		digit = uint64(instruction[1][1])
	} else {
		digit, err = strconv.ParseUint(string(instruction[1]), 0, 16)
		if err != nil {
			return 0, 0
		}
	}
	return uint8(reg), uint16(digit)
}

func (l *Lexer) twoRegOneByte(instruction [][]byte) (uint8, uint8, uint8) {
	if len(instruction) < 3 {
		return 0, 0, 0
	}

	if instruction[0][0] != '%' || instruction[1][0] != '%' {
		return 0, 0, 0
	}

	reg1, err := strconv.ParseUint(string(instruction[0][1:]), 16, 8)
	if err != nil {
		return 0, 0, 0
	}
	if reg1 > 15 {
		reg1 = 0
	}

	reg2, err := strconv.ParseUint(string(instruction[1][1:]), 16, 8)
	if err != nil {
		return 0, 0, 0
	}
	if reg2 > 15 {
		reg2 = 0
	}

	var b uint8
	if instruction[2][0] == '-' {
		digit, _ := strconv.ParseInt(string(instruction[2]), 0, 8)
		b = uint8(digit)
	} else {
		digit, _ := strconv.ParseUint(string(instruction[2]), 0, 8)
		b = uint8(digit)
	}
	return uint8(reg1), uint8(reg2), b
}

func (l *Lexer) twoByte(instruction [][]byte) uint16 {
	if len(instruction) < 1 {
		return 0
	}

	var digit uint64
	var err error
	if instruction[0][0] == '~' { // Label
		bits := fullBits
		label := instruction[0][1:]
		if label[0] == '^' {
			bits = higherBits
			label = label[1:]
		} else if label[0] == '`' {
			bits = lowerBits
			label = label[1:]
		}

		var offset uint16
		addIndex := bytes.Index(instruction[0], []byte{'+'})
		subIndex := bytes.Index(instruction[0], []byte{'-'})
		if addIndex > 0 || subIndex > 0 {
			ind := addIndex
			if subIndex > 0 {
				ind = subIndex
			}
			label = instruction[0][1:ind]
			offset64, err := strconv.ParseInt(string(instruction[0][ind+1:]), 0, 16)
			if err != nil {
				fmt.Printf("Invalid offset on line %d\n", l.linenum)
				os.Exit(1)
			}
			offset = uint16(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}
		if label[0] == '$' {
			digit = uint64(l.currMemLocation + offset)
		} else {
			l.labelPlaces[l.currMemLocation+1] = LabelReplace{
				l:      string(label),
				offset: offset,
				part:   bits,
			}
		}
	} else if instruction[0][0] == '\'' { // Literal byte character
		digit = uint64(instruction[0][1])
	} else {
		digit, err = strconv.ParseUint(string(instruction[0]), 0, 16)
		if err != nil {
			return 0
		}
	}
	return uint16(digit)
}
