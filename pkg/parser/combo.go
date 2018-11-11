package parser

import "github.com/lfkeitel/asml-sim/pkg/token"

func (p *Parser) parseNormalInst(imm, addr, reg byte) {
	p.readToken()
	dest, ok := p.parseRegister()
	if !ok {
		return
	}

	p.readToken()
	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.IDENT) {
		val, ok := p.parseAddress(2)
		if !ok {
			return
		}

		p.p.appendCode(addr, dest, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.IMMEDIATE) {
		p.readToken()
		val, ok := p.parseAddress(2)
		if !ok {
			return
		}

		p.p.appendCode(imm, dest, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.REGISTER) {
		src, ok := p.parseRegister()
		if !ok {
			return
		}

		p.p.appendCode(reg, dest, src)
	} else {
		p.tokenErr(token.NUMBER, token.IMMEDIATE, token.REGISTER)
	}
}

func (p *Parser) parseInstAddrReg(addr, reg byte) {
	p.readToken()
	dest, ok := p.parseRegister()
	if !ok {
		return
	}

	p.readToken()
	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.IDENT) {
		val, ok := p.parseAddress(2)
		if !ok {
			return
		}

		p.p.appendCode(addr, dest, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.REGISTER) {
		src, ok := p.parseRegister()
		if !ok {
			return
		}

		p.p.appendCode(reg, dest, src)
	} else {
		p.tokenErr(token.NUMBER, token.REGISTER)
	}
}

func (p *Parser) parseNormalInstNoDest(imm, addr, reg byte) {
	p.readToken()
	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.IDENT) {
		val, ok := p.parseAddress(1)
		if !ok {
			return
		}

		p.p.appendCode(addr, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.IMMEDIATE) {
		p.readToken()
		val, ok := p.parseAddress(1)
		if !ok {
			return
		}

		p.p.appendCode(imm, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.REGISTER) {
		src, ok := p.parseRegister()
		if !ok {
			return
		}

		p.p.appendCode(reg, src)
	} else {
		p.tokenErr(token.NUMBER, token.IMMEDIATE, token.REGISTER)
	}
}

func (p *Parser) parseInstAddrRegNoDest(addr, reg byte) {
	p.readToken()
	if p.curTokenIs(token.NUMBER) || p.curTokenIs(token.IDENT) {
		val, ok := p.parseAddress(1)
		if !ok {
			return
		}

		p.p.appendCode(addr, uint8(val>>8), uint8(val))
	} else if p.curTokenIs(token.REGISTER) {
		src, ok := p.parseRegister()
		if !ok {
			return
		}

		p.p.appendCode(reg, src)
	} else {
		p.tokenErr(token.NUMBER, token.REGISTER)
	}
}
