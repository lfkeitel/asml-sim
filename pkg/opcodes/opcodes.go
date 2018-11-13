package opcodes

// Language opcodes
const (
	NOOP byte = iota

	ADDA
	ADDI
	ADDR

	ANDA
	ANDI
	ANDR

	ORA
	ORI
	ORR

	XORA
	XORI
	XORR

	ROTR
	ROTL

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
