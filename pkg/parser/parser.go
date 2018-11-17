package parser

import (
	"fmt"

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
		case token.FCB:
			p.rawDataFCB()
		case token.FDB:
			p.rawDataFDB()

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

func (p *Parser) rawDataFCB() {
	p.readToken()

	for {
		if !p.curTokenIs(token.NUMBER) && !p.curTokenIs(token.STRING) {
			p.parseErr("Constant byte must be a number or string")
			return
		}

		if p.curTokenIs(token.STRING) {
			p.p.appendCode([]byte(p.ct.Literal)...)
		} else {
			val, err := parseUint16(p.ct.Literal)
			if err != nil || val > 255 {
				p.parseErr("Invalid bytes")
				return
			}
			p.p.appendCode(uint8(val))
		}

		p.readToken()
		if p.curTokenIs(token.END_INST) {
			break
		}

		if !p.curTokenIs(token.COMMA) {
			p.tokenErr(token.COMMA)
			return
		}
		p.readToken()
	}
}

func (p *Parser) rawDataFDB() {
	p.readToken()

	for {
		if p.curTokenIs(token.STRING) {
			panic("NO!")
		}

		val, ok := p.parseAddress(0)
		if !ok {
			p.parseErr("Invalid bytes")
			return
		}
		p.p.appendCode(uint8(val>>8), uint8(val))

		p.readToken()
		if p.curTokenIs(token.END_INST) {
			break
		}

		if !p.curTokenIs(token.COMMA) {
			p.tokenErr(token.COMMA)
			return
		}
		p.readToken()
	}
}
