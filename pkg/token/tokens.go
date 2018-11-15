package token

import "strconv"

type Type int

type Token struct {
	Type         Type
	Literal      string
	Line, Column int
}

// Language opcodes
const (
	ILLEGAL Type = iota
	EOF
	COMMENT
	END_INST
	COMMA
	IMMEDIATE

	IDENT
	LABEL
	NUMBER
	STRING
	REGISTER

	keyword_beg
	NOOP
	LOAD
	STR
	XFER
	ADD
	OR
	AND
	XOR
	ROTR
	ROTL
	JMP
	HALT
	JMPA
	LDSP
	PUSH
	POP
	CALL
	RTN
	RMB
	ORG
	keyword_end
)

var tokens = [...]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	COMMENT:   "COMMENT",
	END_INST:  "END_INST",
	COMMA:     ",",
	IMMEDIATE: "#",

	// Identifiers & literals
	IDENT:    "IDENT",
	LABEL:    "LABEL",
	NUMBER:   "NUMBER",
	STRING:   "STRING",
	REGISTER: "REGISTER",

	// Keywords
	NOOP: "NOOP",
	LOAD: "LOAD",
	STR:  "STR",
	XFER: "XFER",
	ADD:  "ADD",
	OR:   "OR",
	AND:  "AND",
	XOR:  "XOR",
	ROTR: "ROTR",
	ROTL: "ROTL",
	JMP:  "JMP",
	HALT: "HALT",
	JMPA: "JMPA",
	LDSP: "LDSP",
	PUSH: "PUSH",
	POP:  "POP",
	CALL: "CALL",
	RTN:  "RTN",
	RMB:  "RMB",
	ORG:  "ORG",
}

// Opcodes maps strings to an opcode byte value
var Opcodes map[string]Type

func init() {
	Opcodes = make(map[string]Type, keyword_end-keyword_beg)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		Opcodes[tokens[i]] = i
	}
}

// String returns the string representation of a Type.
// If the Type is an operator or keyword, it will return
// the literal representation of the Type. For example,
// ADD will return "+". Other types will return their
// constant name. IDENT returns "IDENT".
func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// LookupIdent checks if the string is a keyword. It it is,
// the corresponding Type will be returned. Otherwise,
// the Type IDENT will be returned.
func LookupIdent(ident string) Type {
	if tok, ok := Opcodes[ident]; ok {
		return tok
	}
	return IDENT
}

// IsKeyword returns if the given Type is a keyword Token.
func IsKeyword(t Type) bool {
	return keyword_beg < t && t < keyword_end
}

// NewSimpleToken returns a Token with no literal representation
// beyond what can be obtained with Type.String().
func NewSimpleToken(tokType Type, line, col int) Token {
	return NewToken(tokType, "", line, col)
}

// NewToken creates a Token with the type tokType and literal
// representation lit.
func NewToken(tokType Type, lit string, line, col int) Token {
	return Token{
		Type:    tokType,
		Literal: lit,
		Line:    line,
		Column:  col,
	}
}
