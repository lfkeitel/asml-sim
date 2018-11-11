package opcodes

// Language opcodes
const (
	NOOP byte = iota
	ADD
	ADDI

	AND
	OR
	ROT
	XOR

	CALLI
	CALLR
	RTN

	HALT

	JMP
	JMPA

	LDSP
	LDSPI

	LOADA
	LOADI
	LOADR

	STRA
	STRR

	MOVR

	POP
	PUSH
)
