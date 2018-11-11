package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lfkeitel/asml-sim/pkg/lexer"
	"github.com/lfkeitel/asml-sim/pkg/token"
)

type Parser struct {
	l        *lexer.Lexer
	ct, peek token.Token
	p        *Program
	err      error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		p: &Program{
			Code:    make([]uint8, 0, 100),
			Labels:  make(LabelMap),
			LinkMap: make(LabelLinkMap),
		},
	}
	p.readToken()
	p.readToken()
	return p
}

func (p *Parser) Parse() (*Program, error) {
	for p.ct.Type != token.EOF {
		if p.err != nil {
			break
		}

		switch p.ct.Type {
		case token.END_INST:
			break
		case token.LABEL:
			p.makeLabel()
		case token.NUMBER, token.STRING:
			p.rawData()

		case token.LOADI:
			p.insLoadi()
		case token.LOADA:
			p.insLoada()
		case token.LOADR:
			p.insLoadr()

		case token.STRA:
			p.insStorea()
		case token.STRR:
			p.insStorer()

		case token.MOVR:
			p.insMovr()

		case token.ADD:
			p.insAdd()
		case token.ADDI:
			p.insAddi()

		case token.OR:
			p.insOr()
		case token.AND:
			p.insAnd()
		case token.XOR:
			p.insXor()

		case token.ROT:
			p.insRot()

		case token.PUSH:
			p.insPush()
		case token.POP:
			p.insPop()

		case token.CALLR:
			p.insCallr()
		case token.CALLI:
			p.insCalli()

		case token.JMP:
			p.insJmp()
		case token.JMPA:
			p.insJmpa()

		case token.LDSP:
			p.insLdsp()
		case token.LDSPI:
			p.insLdspi()

		case token.HALT:
			p.insHalt()
		case token.NOOP:
			p.insNoop()
		case token.RTN:
			p.insRtn()

		default:
			p.err = fmt.Errorf("line %d, col %d Unknown token %q", p.ct.Line, p.ct.Column, p.ct.Type.String())
		}

		p.readToken()
	}

	return p.p, p.err
}

func (p *Parser) readToken() {
	p.ct = p.peek
	p.peek = p.l.NextToken()

	for p.peek.Type == token.COMMENT {
		p.peek = p.l.NextToken()
	}
}

func (p *Parser) curTokenIs(t token.Type) bool  { return p.ct.Type == t }
func (p *Parser) peekTokenIs(t token.Type) bool { return p.peek.Type == t }

func (p *Parser) tokenErr(t ...token.Type) {
	p.err = fmt.Errorf("expected %q on line %d, got %q", t, p.ct.Line, p.ct.Type)
}

func (p *Parser) parseErr(msg string) {
	p.err = fmt.Errorf("%s on line %d", msg, p.ct.Line)
}

func (p *Parser) makeLabel() {
	p.p.addLabel(p.ct.Literal)
}

func (p *Parser) rawData() {
	switch p.ct.Type {
	case token.STRING:
		p.p.Code = append(p.p.Code, p.ct.Literal...)
		p.p.pc += uint16(len(p.ct.Literal))
	case token.NUMBER:
		var (
			raw uint64
			err error
		)

		if p.ct.Literal[0] == '!' {
			raw, err = strconv.ParseUint(string(p.ct.Literal[1:]), 2, 8)
		} else {
			raw, err = strconv.ParseUint(string(p.ct.Literal), 0, 8)
		}

		if err != nil {
			fmt.Printf("Invalid byte sequence on line %d: %v\n", p.ct.Line, err)
			os.Exit(1)
		}

		p.p.Code = append(p.p.Code, uint8(raw))
		p.p.incPC()
	}
}
