package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	infile       string
	outfile      string
	disableState bool
	printMem     bool
	printLegacy  bool
)

func init() {
	flag.StringVar(&infile, "in", "MachineIn.txt", "Input code")
	flag.StringVar(&outfile, "out", "MachineOut.txt", "Output file")
	flag.BoolVar(&disableState, "nostate", false, "Disable writting state every cycle")
	flag.BoolVar(&printMem, "printmem", false, "Print the initial memory layout and exit")
	flag.BoolVar(&printLegacy, "legacy", false, "Print source code converted for original implementation")
}

func main() {
	flag.Parse()

	code := loadCode()

	if printLegacy {
		for _, b := range code {
			fmt.Printf("%02X\n", b)
		}
		return
	}

	sim := newVM(code, disableState)

	if printMem {
		sim.printVMState()
		os.Stdout.Write(sim.output.Bytes())
		os.Stdout.Write([]byte{'\n'})
		return
	}

	var output io.Writer
	if outfile == "stdout" {
		output = os.Stdout
	} else {
		file, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		output = file
		defer file.Close()
	}

	if err := sim.run(output); err != nil {
		fmt.Println(err.Error())
	}
}

type labelReplace struct {
	l      string
	offset uint8
}

func loadCode() []uint8 {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var code []uint8
	linenum := 0
	labels := make(map[string]uint8)            // Label definitions
	labelPlaces := make(map[uint8]labelReplace) // Memory locations that need labels
	cmemloc := uint8(0)

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
			labels[string(line[1:])] = cmemloc
			continue
		}

		instruction := bytes.SplitN(line, []byte{' '}, 2)
		byte1, err := strconv.ParseUint(string(instruction[0]), 16, 8)
		if err != nil {
			fmt.Printf("Error on line %d\n", line)
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
			labelPlaces[cmemloc+1] = labelReplace{
				l:      string(label),
				offset: offset,
			}
			code = append(code, uint8(byte1), 0)
		} else {
			byte2, err := strconv.ParseUint(string(instruction[1]), 16, 8)
			if err != nil {
				fmt.Printf("Error on line %d\n", linenum)
				os.Exit(1)
			}

			code = append(code, uint8(byte1), uint8(byte2))
		}
		cmemloc += 2
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
