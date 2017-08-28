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
)

func init() {
	flag.StringVar(&infile, "in", "MachineIn.txt", "Input code")
	flag.StringVar(&outfile, "out", "MachineOut.txt", "Output file")
	flag.BoolVar(&disableState, "nostate", false, "Disable writting state every cycle")
	flag.BoolVar(&printMem, "printmem", false, "Print the initial memory layout and exit")
}

func main() {
	flag.Parse()

	code := loadCode()

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

		instruction := bytes.SplitN(line, []byte{' '}, 2)
		byte1, err := strconv.ParseUint(string(instruction[0]), 16, 8)
		if err != nil {
			fmt.Printf("Error on line %d\n", line)
			os.Exit(1)
		}

		byte2, err := strconv.ParseUint(string(instruction[1]), 16, 8)
		if err != nil {
			fmt.Printf("Error on line %d\n", linenum)
			os.Exit(1)
		}

		code = append(code, uint8(byte1), uint8(byte2))
	}

	return code
}
