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

type labelReplace struct {
	l      string
	offset uint16
	part   bitParts
}

type Lexer struct {
	in              io.ReadSeeker
	labels          map[string]uint16       // Label definitions
	labelPlaces     map[uint16]labelReplace // Memory locations that need labels
	currMemLocation uint16
	linenum         int
}

func New(in io.ReadSeeker) *Lexer {
	return &Lexer{
		in: in,
	}
}

func (l *Lexer) Lex() []uint8 {
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

	// Rewind file to read in as source
	l.in.Seek(0, 0)
	reader := bufio.NewReader(l.in)
	var code []uint8
	l.linenum = 0
	l.labels = make(map[string]uint16)            // Label definitions
	l.labelPlaces = make(map[uint16]labelReplace) // Memory locations that need labels
	l.currMemLocation = 0
	directives := &token.Flags{
		Size: token.EightBit,
	}

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
				fmt.Println("Directives must be before any code")
				os.Exit(1)
			}

			line = line[1:]
			directive := bytes.Split(line, []byte{' '})
			if len(directive) != 2 {
				fmt.Printf("Invalid directive on line %d\n", l.linenum)
				os.Exit(1)
			}

			switch string(directive[0]) {
			case "bits":
				if string(directive[1]) == "16" {
					directives.Size = token.SixteenBit
				} else if string(directive[1]) == "8" {
					directives.Size = token.EightBit
				} else {
					fmt.Printf("Invalid bitsize on line %d\n", l.linenum)
					os.Exit(1)
				}
			}

			continue
		}

		instruction := bytes.Split(line, []byte{' '})
		opcode, valid := token.Opcodes[string(instruction[0])]
		if !valid { // Literal bytes
			for _, rawbyte := range instruction {
				if rawbyte[0] == '0' && len(rawbyte) > 1 && (rawbyte[1] == 'x' || rawbyte[1] == 'X') {
					rawbyte = rawbyte[2:]
				}

				raw, err := strconv.ParseUint(string(rawbyte), 16, 8)
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
			code = append(code, opcode<<4, 0)
		case token.ROT:
			reg, size := l.oneRegOneDigit(instruction[1:])
			code = append(code, opcode<<4+reg, size)
		case token.MOVR, token.STRR, token.LOADR:
			reg1, reg2 := l.twoRegisters(instruction[1:])
			code = append(code, opcode<<4, reg1<<4+reg2)
		case token.ADD, token.FLAGS, token.OR, token.AND, token.XOR:
			reg1, reg2, reg3 := l.threeRegisters(instruction[1:])
			code = append(code, opcode<<4+reg1, reg2<<4+reg3)
		case token.LOADI:
			reg, b := l.oneRegOneByte(instruction[1:])
			code = append(code, opcode<<4+reg, b)
		case token.LOADA, token.STRA, token.JMP:
			reg, b := l.oneRegTwoByte(instruction[1:])
			if directives.Size == token.EightBit && b > 255 {
				fmt.Printf("Address too big for eight bits on line %d\n", l.linenum)
				os.Exit(1)
			}

			if directives.Size == token.EightBit {
				code = append(code, opcode<<4+reg, uint8(b))
			} else {
				code = append(code, opcode<<4+reg, uint8(b>>8), uint8(b))
				l.currMemLocation++
			}
		default:
			fmt.Printf("Invalid opcode on line %d\n", l.linenum)
			os.Exit(1)
		}

		l.currMemLocation += 2
	}

	// Replace labels
	for loc, label := range l.labelPlaces {
		memloc, exists := l.labels[label.l]
		if !exists {
			fmt.Printf("Label %s not defined\n", label.l)
			os.Exit(1)
		}

		if directives.Size == token.EightBit {
			code[loc] = uint8(memloc) + uint8(label.offset)
		} else {
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

	for name, memloc := range l.labels {
		fmt.Printf("%s: %04X\n", name, memloc)
	}

	return append(directives.Bytes(), code...)
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

func (l *Lexer) oneRegOneByte(instruction [][]byte) (uint8, uint8) {
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
			offset64, err := strconv.ParseInt(string(instruction[1][ind+1:]), 16, 8)
			if err != nil {
				fmt.Printf("Invalid offset on line %d\n", l.linenum)
				os.Exit(1)
			}
			offset = uint16(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}
		l.labelPlaces[l.currMemLocation+1] = labelReplace{
			l:      string(label),
			offset: offset,
			part:   bits,
		}
	} else if instruction[1][0] == '\'' { // Literal byte character
		digit = uint64(instruction[1][1])
	} else { // Byte hex
		if instruction[1][0] == '0' && len(instruction[1]) > 1 && (instruction[1][1] == 'x' || instruction[1][1] == 'X') {
			instruction[1] = instruction[1][2:]
		}
		digit, err = strconv.ParseUint(string(instruction[1]), 16, 8)
		if err != nil {
			return 0, 0
		}
	}
	return uint8(reg), uint8(digit)
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
			offset64, err := strconv.ParseInt(string(instruction[1][ind+1:]), 16, 16)
			if err != nil {
				fmt.Printf("Invalid offset on line %d\n", l.linenum)
				os.Exit(1)
			}
			offset = uint16(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}
		l.labelPlaces[l.currMemLocation+1] = labelReplace{
			l:      string(label),
			offset: offset,
			part:   bits,
		}
	} else if instruction[1][0] == '\'' { // Literal byte character
		digit = uint64(instruction[1][1])
	} else { // Byte hex
		if instruction[1][0] == '0' && len(instruction[1]) > 1 && (instruction[1][1] == 'x' || instruction[1][1] == 'X') {
			instruction[1] = instruction[1][2:]
		}
		digit, err = strconv.ParseUint(string(instruction[1]), 16, 16)
		if err != nil {
			return 0, 0
		}
	}
	return uint8(reg), uint16(digit)
}
