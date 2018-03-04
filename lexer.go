package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	asmlHeader = []byte("ASML")
)

type labelReplace struct {
	l      string
	offset uint8
}

type lexer struct {
	in              io.ReadSeeker
	labels          map[string]uint8       // Label definitions
	labelPlaces     map[uint8]labelReplace // Memory locations that need labels
	currMemLocation uint8
	linenum         int
}

func loadCode() []uint8 {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	lex := &lexer{
		in: file,
	}
	return lex.lex()
}

func (l *lexer) lex() []uint8 {
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

	if bytes.Equal(header, asmlHeader) {
		var buf bytes.Buffer
		io.Copy(&buf, l.in)
		return buf.Bytes()
	}

	// Rewind file to read in as source
	l.in.Seek(0, 0)
	reader := bufio.NewReader(l.in)
	var code []uint8
	l.linenum = 0
	l.labels = make(map[string]uint8)            // Label definitions
	l.labelPlaces = make(map[uint8]labelReplace) // Memory locations that need labels
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

		instruction := bytes.Split(line, []byte{' '})
		opcode, valid := opcodes[string(instruction[0])]
		if !valid { // Literal bytes
			for _, rawbyte := range instruction {
				raw, err := strconv.ParseUint(string(rawbyte), 16, 8)
				if err != err {
					fmt.Printf("Invalid byte sequence on line %d\n", l.linenum)
					os.Exit(1)
				}

				code = append(code, uint8(raw))
			}
			continue
		}

		switch opcode {
		case HALT, NOOP:
			code = append(code, opcode<<4, 0)
		case ROT:
			reg, size := l.oneRegOneDigit(instruction[1:])
			code = append(code, opcode<<4+reg, size)
		case MOVR, STRR, LOADR:
			reg1, reg2 := l.twoRegisters(instruction[1:])
			code = append(code, opcode<<4, reg1<<4+reg2)
		case ADD, ADDF, OR, AND, XOR:
			reg1, reg2, reg3 := l.threeRegisters(instruction[1:])
			code = append(code, opcode<<4+reg1, reg2<<4+reg3)
		case LOADA, LOADI, STRA, JMP:
			reg, b := l.oneRegOneByte(instruction[1:])
			code = append(code, opcode<<4+reg, b)
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
		code[loc] = memloc + label.offset
	}

	return code
}

func (l *lexer) oneRegOneDigit(instruction [][]byte) (uint8, uint8) {
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

func (l *lexer) twoRegisters(instruction [][]byte) (uint8, uint8) {
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

func (l *lexer) threeRegisters(instruction [][]byte) (uint8, uint8, uint8) {
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

func (l *lexer) oneRegOneByte(instruction [][]byte) (uint8, uint8) {
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
		label := instruction[1][1:]
		var offset uint8
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
			offset = uint8(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}
		l.labelPlaces[l.currMemLocation+1] = labelReplace{
			l:      string(label),
			offset: offset,
		}
	} else if instruction[1][0] == '\'' { // Literal byte character
		digit = uint64(instruction[1][1])
	} else { // Byte hex
		digit, err = strconv.ParseUint(string(instruction[1]), 16, 8)
		if err != nil {
			return 0, 0
		}
	}
	return uint8(reg), uint8(digit)
}

// Old lexer for original syntax, kept for maybe a compatibility mode
func (l *lexer) lexold() []uint8 {
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

	if bytes.Equal(header, asmlHeader) {
		var buf bytes.Buffer
		io.Copy(&buf, l.in)
		return buf.Bytes()
	}

	// Rewind file to read in as source
	l.in.Seek(0, 0)
	reader := bufio.NewReader(l.in)
	var code []uint8
	linenum := 0
	labels := make(map[string]uint8)            // Label definitions
	labelPlaces := make(map[uint8]labelReplace) // Memory locations that need labels
	currMemLocation := uint8(0)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		linenum++

		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == ';' { // comment/blank line
			continue
		}

		if line[0] == ':' { // label definition
			labels[string(line[1:])] = currMemLocation
			continue
		}

		instruction := bytes.SplitN(line, []byte{' '}, 2)
		byte1, err := strconv.ParseUint(string(instruction[0]), 16, 8)
		if err != nil {
			fmt.Printf("Error on line %d, invalid byte\n", linenum)
			os.Exit(1)
		}

		if len(instruction) != 2 {
			fmt.Printf("Error on line %d, expected two bytes got 1\n", linenum)
			os.Exit(1)
		}

		if instruction[1][0] == '~' {
			label := instruction[1][1:]
			var offset uint8
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
					fmt.Printf("Invalid offset on line %d\n", linenum)
					os.Exit(1)
				}
				offset = uint8(offset64)
				if subIndex > 0 {
					offset = -offset
				}
			}
			labelPlaces[currMemLocation+1] = labelReplace{
				l:      string(label),
				offset: offset,
			}
			code = append(code, uint8(byte1), 0)
		} else {
			byte2, err := strconv.ParseUint(string(instruction[1]), 16, 8)
			if err != nil {
				fmt.Printf("Error on line %d, invalid byte\n", linenum)
				os.Exit(1)
			}

			code = append(code, uint8(byte1), uint8(byte2))
		}
		currMemLocation += 2
	}

	// Replace labels
	for loc, label := range labelPlaces {
		memloc, exists := labels[label.l]
		if !exists {
			fmt.Printf("Label %s not defined\n", label.l)
			os.Exit(1)
		}
		code[loc] = memloc + label.offset
	}

	return code
}
