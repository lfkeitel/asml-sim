package token

// AddressSize is the size of a memory address
type AddressSize uint8

// Available address sizes, these do not denote the CPU bit size
const (
	EightBit   AddressSize = 0
	SixteenBit AddressSize = 1
)

// Flags contain compiler directives from a source file
type Flags struct {
	Size AddressSize
}

func FlagsFromBytes(code []byte) (*Flags, bool) {
	f := &Flags{}
	if len(code) < 2 {
		return f, false
	}

	opcode := code[0] >> 4 // First 4 bits of byte 1
	//operand1 := code[0] & 15 // Second 4 bits of byte 1
	operand2 := code[1] >> 4 // First 4 bits of byte 2
	//operand3 := code[1] & 15 // Second 4 bits of byte 2

	if opcode != FLAGS {
		return f, false
	}

	if operand2 == 8 {
		f.Size = SixteenBit
	}

	return f, true
}

func (f *Flags) Bytes() []byte {
	flags := []byte{FLAGS << 4, 0}

	if f.Size == SixteenBit {
		flags[1] = flags[1] | 128
	}

	return flags
}
