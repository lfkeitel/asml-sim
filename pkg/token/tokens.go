package token

// Language opcodes
const (
	NOOP  byte = 0x0
	LOADA byte = 0x1
	LOADI byte = 0x2
	STRA  byte = 0x3
	MOVR  byte = 0x4
	ADD   byte = 0x5
	ADDI  byte = 0x6
	OR    byte = 0x7
	AND   byte = 0x8
	XOR   byte = 0x9
	ROT   byte = 0xA
	JMP   byte = 0xB
	HALT  byte = 0xC
	STRR  byte = 0xD
	LOADR byte = 0xE
	JMPA  byte = 0xF
)

var Opcodes = map[string]byte{
	"NOOP":  NOOP,
	"LOADA": LOADA,
	"LOADI": LOADI,
	"STRA":  STRA,
	"MOVR":  MOVR,
	"ADD":   ADD,
	"ADDI":  ADDI,
	"OR":    OR,
	"AND":   AND,
	"XOR":   XOR,
	"ROT":   ROT,
	"JMP":   JMP,
	"HALT":  HALT,
	"STRR":  STRR,
	"LOADR": LOADR,
	"JMPA":  JMPA,
}
