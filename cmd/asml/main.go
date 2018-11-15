package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lfkeitel/asml-sim/pkg/lexer"
	"github.com/lfkeitel/asml-sim/pkg/linker"
	"github.com/lfkeitel/asml-sim/pkg/parser"
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

	if code == nil {
		os.Exit(1)
	}

	if compile {
		writeCompiledCode(code)
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

func loadCode(infile string) []parser.CodePart {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	if code := checkBinaryFile(file); code != nil {
		return code
	}
	file.Seek(0, 0)

	lex := lexer.New(file)
	p := parser.New(lex)
	program, err := p.Parse()
	if err != nil {
		fmt.Printf("Parsing failed: %q\n", err)
		return nil
	}

	if err := linker.Link(program); err != nil {
		fmt.Printf("Linking failed: %q\n", err)
		return nil
	}

	return program.Parts
}

func checkBinaryFile(file *os.File) []parser.CodePart {
	// Read in a compiled ASML file
	header := make([]byte, 18)
	n, err := file.Read(header)
	if err != nil {
		fmt.Printf("Error reading file header: %s\n", err)
		os.Exit(1)
	}
	if n < 4 {
		fmt.Println("Invalid file")
		os.Exit(1)
	}

	if bytes.Equal(header, []byte("S007000041534D4CCB")) {
		fmt.Println("SRECORDS cannot be read yet")
		os.Exit(1)
	}
	return nil
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

func writeCompiledCode(code []uint8) {
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
}
