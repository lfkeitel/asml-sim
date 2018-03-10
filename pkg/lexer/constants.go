package lexer

// Language opcodes
const (
	NOOP  byte = 0
	LOADA byte = 1
	LOADI byte = 2
	STRA  byte = 3
	MOVR  byte = 4
	ADD   byte = 5
	// ADDF  byte = 6
	OR    byte = 7
	AND   byte = 8
	XOR   byte = 9
	ROT   byte = 10
	JMP   byte = 11
	HALT  byte = 12
	STRR  byte = 13
	LOADR byte = 14
)

var opcodes = map[string]byte{
	"NOOP":  NOOP,
	"LOADA": LOADA,
	"LOADI": LOADI,
	"STRA":  STRA,
	"MOVR":  MOVR,
	"ADD":   ADD,
	// "ADDF":  ADDF,
	"OR":    OR,
	"AND":   AND,
	"XOR":   XOR,
	"ROT":   ROT,
	"JMP":   JMP,
	"HALT":  HALT,
	"STRR":  STRR,
	"LOADR": LOADR,
}

/*
Runtime:
-----------------------
LOADI %C 0xFF
LOADI %D 1
LOADI %E 0xFE
LOADI %F ~exit

JMP %0 ~main

:exit
    HALT

:return
    STRA %F ~$+3
	JMP %0 ~exit
-----------------------

The main label is supplied by user code and resolved at link time.
*/
var runtime = []byte{0x2C, 0xFF, 0x2D, 0x01, 0x2E, 0xFE, 0x2F, 0x0A, 0xB0, 0x10, 0xC0, 0x00, 0x3F, 0x0F, 0xB0, 0x0A}
var runtimeLabels = map[string]uint8{
	"exit":   0x0A,
	"return": 0x0C,
}
var mainLabelLoc = uint8(9)
