package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lfkeitel/asml-sim/pkg/lexer"
	"github.com/lfkeitel/asml-sim/pkg/vm"
)

var (
	outfile      string
	showState    bool
	printMem     bool
	printLegacy  bool
	compile      bool
	printVersion bool

	version   string
	buildTime string
	builder   string
	goversion string
)

func init() {
	flag.StringVar(&outfile, "out", "stdout", "Output file")
	flag.BoolVar(&showState, "state", false, "Write state every cycle")
	flag.BoolVar(&printMem, "printmem", false, "Print the initial memory layout and exit")
	flag.BoolVar(&compile, "compile", false, "Compile file to ASML program")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
}

func main() {
	flag.Parse()

	if printVersion {
		printVersionInfo()
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	infile := flag.Arg(0)
	code := loadCode(infile)

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

		out.Write(lexer.ASMLHeader)
		out.Write(code)
		out.Close()
		return
	}

	sim := vm.New(code, showState)

	if printMem {
		sim.PrintState()
		os.Stdout.Write(sim.Output())
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

	if err := sim.Run(output); err != nil {
		fmt.Println(err.Error())
	}
}

func loadCode(infile string) []uint8 {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	lex := lexer.New(file)
	return lex.Lex()
}

func printVersionInfo() {
	fmt.Printf(`ASML Virtual Machine - (C) 2018 Lee Keitel
Architecture: 8-bit registers, 16-bit addresses
Version:      %s
Built:        %s
Compiled by:  %s
Go version:   %s
`, version, buildTime, builder, goversion)
}
