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

	CALLA
	CALLR
	RTN

	HALT

	JMP
	JMPA

	LDSPA
	LDSPI
	LDSPR

	LOADA
	LOADI
	LOADR

	STRA
	STRR

	XFER

	POP
	PUSH
)
