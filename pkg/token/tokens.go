package token

// Language opcodes
const (
	NOOP  byte = 0
	LOADA byte = 1
	LOADI byte = 2
	STRA  byte = 3
	MOVR  byte = 4
	ADD   byte = 5
	FLAGS byte = 6
	OR    byte = 7
	AND   byte = 8
	XOR   byte = 9
	ROT   byte = 10
	JMP   byte = 11
	HALT  byte = 12
	STRR  byte = 13
	LOADR byte = 14
)

var Opcodes = map[string]byte{
	"NOOP":  NOOP,
	"LOADA": LOADA,
	"LOADI": LOADI,
	"STRA":  STRA,
	"MOVR":  MOVR,
	"ADD":   ADD,
	"ADDF":  FLAGS,
	"OR":    OR,
	"AND":   AND,
	"XOR":   XOR,
	"ROT":   ROT,
	"JMP":   JMP,
	"HALT":  HALT,
	"STRR":  STRR,
	"LOADR": LOADR,
}
