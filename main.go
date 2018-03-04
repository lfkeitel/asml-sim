package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	infile      string
	outfile     string
	showState   bool
	printMem    bool
	printLegacy bool
	compile     bool
)

func init() {
	flag.StringVar(&infile, "in", "MachineIn.txt", "Input code")
	flag.StringVar(&outfile, "out", "stdout", "Output file")
	flag.BoolVar(&showState, "state", false, "Disable writting state every cycle")
	flag.BoolVar(&printMem, "printmem", false, "Print the initial memory layout and exit")
	flag.BoolVar(&printLegacy, "legacy", false, "Print source code converted for original implementation")
	flag.BoolVar(&compile, "compile", false, "Compile file to ASML program")
}

func main() {
	flag.Parse()

	code := loadCode()

	if printLegacy {
		for i, b := range code {
			if i&1 != 1 { // Check even indexes
				if b>>4 > HALT { // Only opcodes 1-12 were in the original implementation
					fmt.Println("ERROR: Opcodes D and E are not available in the legacy implementation")
					return
				}
			}
			fmt.Printf("%02X\n", b)
		}
		return
	}

	if compile {
		var out io.WriteCloser
		if outfile == "stdout" {
			out = os.Stdout
		} else {
			var err error
			out, err = os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		out.Write(asmlHeader)
		out.Write(code)
		out.Close()
		return
	}

	sim := newVM(code, showState)

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
