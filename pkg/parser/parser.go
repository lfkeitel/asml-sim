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
			Parts:  []CodePart{newCodePart(0)},
			Labels: make(LabelMap),
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

		case token.LOAD:
			p.insLoad()

		case token.STR:
			p.insStore()

		case token.XFER:
			p.insMovr()

		case token.ADD:
			p.insAdd()

		case token.OR:
			p.insOr()
		case token.AND:
			p.insAnd()
		case token.XOR:
			p.insXor()

		case token.ROTR:
			p.insRotr()
		case token.ROTL:
			p.insRotl()

		case token.PUSH:
			p.insPush()
		case token.POP:
			p.insPop()

		case token.CALL:
			p.insCall()

		case token.JMP:
			p.insJmp()
		case token.JMPA:
			p.insJmpa()

		case token.LDSP:
			p.insLdsp()

		case token.HALT:
			p.insHalt()
		case token.NOOP:
			p.insNoop()
		case token.RTN:
			p.insRtn()

		case token.RMB:
			p.insRmb()
		case token.ORG:
			p.insOrg()

		default:
			p.err = fmt.Errorf("line %d, col %d Unknown token %v", p.ct.Line, p.ct.Column, p.ct.Type.String())
		}

		p.readToken()
	}

	if p.err != nil {
		return p.p, p.err
	}

	p.err = p.p.validate()
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
	p.err = fmt.Errorf("expected %v on line %d, got %v", t, p.ct.Line, p.ct.Type)
}

func (p *Parser) expectToken(t token.Type) {
	p.readToken()
	if !p.curTokenIs(t) {
		p.tokenErr(t)
	}
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
		p.p.appendCode([]byte(p.ct.Literal)...)
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

		p.p.appendCode(uint8(raw))
	}
}
