// Code generated by runtime gen. DO NOT EDIT.

package lexer

/*
Runtime Source:

-----------------------
LOADI %C 0xFFFE
LOADI %D ~exit

JMPA ~main

:exit
    HALT

:return
    STRA %D ~$+5
    JMPA ~exit

-----------------------

The main label is supplied by user code and resolved at link time.
*/

var runtime = []byte{0x2, 0xc, 0xff, 0xfe, 0x2, 0xd, 0x0, 0x0, 0xf, 0x0, 0x0, 0xc, 0x3, 0xd, 0x0, 0x11, 0xf, 0x0, 0x0}
var runtimeLabels = map[string]uint16{ 
	"exit": 0xb,
	"return": 0xc,
}
var mainLabelLoc = uint16(0x9)
