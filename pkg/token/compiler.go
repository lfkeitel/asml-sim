package token

// Flags contain compiler directives from a source file
type Flags struct{}

func FlagsFromBytes(code []byte) (*Flags, bool) {
	f := &Flags{}
	return f, true
}

func (f *Flags) Bytes() []byte {
	flags := []byte{FLAGS << 4, 0}
	return flags
}
