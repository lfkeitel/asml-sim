package parser

import (
	"strconv"
	"strings"

	"github.com/lfkeitel/asml-sim/pkg/token"
)

func parseRegisterLiteral(s string) (uint8, error) {
	reg, err := strconv.ParseUint(s, 16, 8)
	return uint8(reg), err
}

func parseUint16(s string) (uint16, error) {
	var (
		val uint64
		err error
	)

	if s[0] == '!' {
		val, err = strconv.ParseUint(s[1:], 2, 16)
	} else {
		val, err = strconv.ParseUint(s, 0, 16)
	}

	return uint16(val), err
}

func (p *Parser) parseRegister() (uint8, bool) {
	if !p.curTokenIs(token.REGISTER) {
		p.tokenErr(token.REGISTER)
		return 0, false
	}

	regT := p.ct

	reg, err := parseRegisterLiteral(regT.Literal)
	if err != nil || reg > 13 {
		p.parseErr("invalid register")
		return 0, false
	}

	return reg, true
}

func (p *Parser) parseUint16() (uint16, bool) {
	if !p.curTokenIs(token.NUMBER) {
		p.tokenErr(token.NUMBER)
		return 0, false
	}

	valT := p.ct

	val, err := parseUint16(valT.Literal)
	if err != nil {
		p.parseErr("invalid immediate value")
		return 0, false
	}

	return val, true
}

func (p *Parser) parseAddress(pcoffset uint16) (uint16, bool) {
	var val uint16
	var err error

	if p.curTokenIs(token.NUMBER) {
		valT := p.ct

		val, err = parseUint16(valT.Literal)
		if err != nil {
			p.parseErr("invalid address")
			return 0, false
		}
	} else if p.curTokenIs(token.IDENT) {
		label := p.ct.Literal
		lit := label
		var offset uint16

		addIndex := strings.Index(lit, "+")
		subIndex := strings.Index(lit, "-")

		if addIndex > 0 || subIndex > 0 {
			ind := addIndex
			if subIndex > 0 {
				ind = subIndex
			}
			label = lit[:ind]
			offset64, err := strconv.ParseInt(string(lit[ind+1:]), 0, 16)
			if err != nil {
				p.parseErr("invalid address offset")
				return 0, false
			}

			offset = uint16(offset64)
			if subIndex > 0 {
				offset = -offset
			}
		}

		if label == "$" {
			val = p.p.pc() + offset
		} else {
			p.p.addLink(pcoffset, label, int16(offset))
		}
	} else if p.curTokenIs(token.STRING) {
		seq := p.ct.Literal

		switch len(seq) {
		case 0:
			return 0, true
		case 1:
			return uint16(seq[0]), true
		case 2:
			return uint16(seq[0])<<8 + uint16(seq[1]), true
		default:
			p.parseErr("string too long")
			return 0, false
		}
	} else {
		p.tokenErr(token.LABEL, token.NUMBER)
		return 0, false
	}

	return val, true
}
