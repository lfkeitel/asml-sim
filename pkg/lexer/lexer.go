package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/lfkeitel/asml-sim/pkg/token"
)

var (
	ASMLHeader = []byte("ASML")
)

type Lexer struct {
	input  *bufio.Reader
	curCh  byte // current char under examination
	peekCh byte // peek character
	line   int  // line in source file
	column int  // column in current line
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{
		input:  bufio.NewReader(reader),
		line:   1,
		column: 1,
	}
	// Populate both current and peek char
	l.readChar()
	l.readChar()
	return l
}

func NewString(input string) *Lexer {
	return New(strings.NewReader(input))
}

func (l *Lexer) readChar() {
	l.curCh = l.peekCh

	var err error
	l.peekCh, err = l.input.ReadByte()
	if err != nil {
		l.peekCh = 0
	}

	if l.curCh == '\r' {
		l.readChar()
	}
	l.column++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.curCh == '\n' {
		tok = token.NewSimpleToken(token.END_INST, l.line, l.column)
		l.resetPos()
		l.readChar()
		return tok
	}

	l.devourWhitespace()

	switch l.curCh {
	case ':':
		l.readChar()
		tok = token.NewToken(token.LABEL, l.readIdentifier(), l.line, l.column)
	case '#':
		tok = token.NewSimpleToken(token.IMMEDIATE, l.line, l.column)
	case '"':
		tok = token.NewToken(token.STRING, l.readString(), l.line, l.column)
	case ';':
		col := l.column
		tok = token.NewToken(token.COMMENT, l.readSingleLineComment(), l.line, col)
		l.resetPos()
	case '%':
		l.readChar()
		if l.curCh == 'S' && l.peekCh == 'P' {
			l.readChar()
			tok = token.NewToken(token.REGISTER, "SP", l.line, l.column)
		} else {
			tok = token.NewToken(token.REGISTER, string(l.curCh), l.line, l.column)
		}
	case 0:
		tok = token.NewSimpleToken(token.EOF, l.line, l.column)

	default:
		if isLetter(l.curCh) {
			lit := l.readIdentifier()
			tokType := token.LookupIdent(lit)
			tok = token.NewToken(tokType, lit, l.line, l.column)
			return tok
		} else if isDigit(l.curCh) || l.curCh == '!' {
			tok = l.readNumber()
			return tok
		}

		tok = token.NewSimpleToken(token.ILLEGAL, l.line, l.column)
	}

	l.readChar()
	return tok
}

func (l *Lexer) resetPos() {
	l.line++
	l.column = 0
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isIdent(l.curCh) {
		ident.WriteByte(l.curCh)
		l.readChar()
	}
	return ident.String()
}

// TODO: Support escape sequences, standard Go should be fine, or PHP.
func (l *Lexer) readString() string {
	var ident bytes.Buffer
	l.readChar() // Go past the starting double quote

	for l.curCh != '"' {
		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return ident.String()
}

func (l *Lexer) readNumber() token.Token {
	var ident bytes.Buffer

	for isDigit(l.curCh) || isHexDigit(l.curCh) || l.curCh == '!' {
		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return token.NewToken(token.NUMBER, ident.String(), l.line, l.column)
}

func (l *Lexer) readSingleLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over semicolon

	for l.curCh != '\n' && l.curCh != 0 {
		com.WriteByte(l.curCh)
		l.readChar()
	}
	return strings.TrimSpace(com.String())
}

func (l *Lexer) devourWhitespace() {
	for isWhitespace(l.curCh) {
		l.readChar()
	}
}

func isIdent(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '-' || ch == '+'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

func isHexDigit(ch byte) bool {
	return 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F' || ch == 'x'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
