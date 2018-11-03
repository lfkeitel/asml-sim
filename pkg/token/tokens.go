package token

// Language opcodes
const (
	NOOP byte = iota
	LOADA
	LOADI
	STRA
	MOVR
	ADD
	ADDI
	OR
	AND
	XOR
	ROT
	JMP
	HALT
	STRR
	LOADR
	JMPA
	LDSP
	LDSPI
	PUSH
	POP
)

// Opcodes maps strings to an opcode byte value
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
	"LDSP":  LDSP,
	"LDSPI": LDSPI,
	"PUSH":  PUSH,
	"POP":   POP,
}
