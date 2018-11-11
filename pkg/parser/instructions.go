package parser

import (
	"github.com/lfkeitel/asml-sim/pkg/opcodes"
	"github.com/lfkeitel/asml-sim/pkg/token"
)

func (p *Parser) insLoad() { p.parseNormalInst(opcodes.LOADI, opcodes.LOADA, opcodes.LOADR) }

func (p *Parser) insStore() { p.parseInstAddrReg(opcodes.STRA, opcodes.STRR) }

func (p *Parser) insMovr() { p.parseRegReg(opcodes.XFER) }

func (p *Parser) insAdd()  { p.parseRegRegReg(opcodes.ADD) }
func (p *Parser) insAddi() { p.parseRegRegHalfNumber(opcodes.ADDI) }

func (p *Parser) insOr()  { p.parseRegRegReg(opcodes.OR) }
func (p *Parser) insAnd() { p.parseRegRegReg(opcodes.AND) }
func (p *Parser) insXor() { p.parseRegRegReg(opcodes.XOR) }

func (p *Parser) insRot() { p.parseRegHalfNumber(opcodes.ROT) }

func (p *Parser) insJmp()  { p.parseRegNumber(opcodes.JMP) }
func (p *Parser) insJmpa() { p.parseNumber(opcodes.JMPA) }

func (p *Parser) insLdsp() { p.parseNormalInstNoDest(opcodes.LDSPI, opcodes.LDSPA, opcodes.LDSPR) }

func (p *Parser) insPush() { p.parseReg(opcodes.PUSH) }
func (p *Parser) insPop()  { p.parseReg(opcodes.POP) }

func (p *Parser) insCall() { p.parseInstAddrRegNoDest(opcodes.CALLA, opcodes.CALLR) }

func (p *Parser) insHalt() { p.parseNoArgs(opcodes.HALT) }
func (p *Parser) insNoop() { p.parseNoArgs(opcodes.NOOP) }
func (p *Parser) insRtn()  { p.parseNoArgs(opcodes.RTN) }

// Common argument parsers

func (p *Parser) parseNoArgs(c byte) {
	p.p.appendCode(c)

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseReg(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg)

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseRegReg(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 2
	p.readToken()
	reg2, ok := p.parseRegister()
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg, reg2)

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseRegRegReg(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 2
	p.readToken()
	reg2, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 3
	p.readToken()
	reg3, ok := p.parseRegister()
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg, reg2, reg3)

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseRegNumber(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 2
	p.readToken()
	val, ok := p.parseAddress(2)
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg, uint8(val>>8), uint8(val))

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseRegHalfNumber(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 2
	p.readToken()
	val, ok := p.parseAddress(2)
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg, uint8(val))

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseRegRegHalfNumber(c byte) {
	// Arg 1
	p.readToken()
	reg, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 2
	p.readToken()
	reg2, ok := p.parseRegister()
	if !ok {
		return
	}

	// Arg 3
	p.readToken()
	val, ok := p.parseAddress(3)
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, reg, reg2, uint8(val))

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}

func (p *Parser) parseNumber(c byte) {
	// Arg 1
	p.readToken()
	val, ok := p.parseAddress(1)
	if !ok {
		return
	}

	// Write code
	p.p.appendCode(c, uint8(val>>8), uint8(val))

	if p.peekTokenIs(token.END_INST) {
		p.readToken()
	}
}
