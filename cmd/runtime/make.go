package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/lfkeitel/asml-sim/pkg/lexer"
)

var (
	infile, outfile string
)

func init() {
	flag.StringVar(&infile, "in", "", "Input file")
	flag.StringVar(&outfile, "out", "", "Output file")
}

func main() {
	flag.Parse()
	if infile == "" || outfile == "" {
		fmt.Println("-infile and -outfile flags required")
		os.Exit(1)
	}

	file, err := os.Open(infile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	source, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	file.Seek(0, 0)

	lex := lexer.New(file)
	code, labels, linklocs := lex.LexNoLink()

	mainLinks := linklocs.FindOffsets("main")
	if len(mainLinks) == 0 {
		fmt.Println("No jmp to main found")
		os.Exit(1)
	}

	data := map[string]interface{}{
		"source": string(source),
		"code":   code,
		"labels": labels,
		"main":   mainLinks[0],
	}

	file2, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file2.Close()
	if err := outtmpl.Execute(file2, data); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var outtmpl = template.Must(template.New("").Parse(`// Code generated by runtime gen. DO NOT EDIT.

package lexer

/*
Runtime Source:

-----------------------
{{.source}}
-----------------------

The main label is supplied by user code and resolved at link time.
*/

var runtime = {{printf "%#v" .code}}
var runtimeLabels = map[string]uint16{ {{range $k, $v := .labels}}
	"{{$k}}": {{printf "%#v" $v}},{{end}}
}
var mainLabelLoc = uint16({{printf "%#v" .main}})
`))