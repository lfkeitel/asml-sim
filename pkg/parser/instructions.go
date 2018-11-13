package parser

import (
	"github.com/lfkeitel/asml-sim/pkg/opcodes"
	"github.com/lfkeitel/asml-sim/pkg/token"
)

func (p *Parser) insLoad() { p.parseInst(opcodes.LOADI, opcodes.LOADA, opcodes.LOADR) }

func (p *Parser) insStore() { p.parseInstNoImm(opcodes.STRA, opcodes.STRR) }

func (p *Parser) insMovr() { p.parseRegReg(opcodes.XFER) }

func (p *Parser) insAdd() { p.parseInst(opcodes.ADDI, opcodes.ADDA, opcodes.ADDR) }

func (p *Parser) insOr()  { p.parseInst(opcodes.ORI, opcodes.ORA, opcodes.ORR) }
func (p *Parser) insAnd() { p.parseInst(opcodes.ANDI, opcodes.ANDA, opcodes.ANDR) }
func (p *Parser) insXor() { p.parseInst(opcodes.XORI, opcodes.XORA, opcodes.XORR) }

func (p *Parser) insRotr() { p.parseRegHalfNumber(opcodes.ROTR) }
func (p *Parser) insRotl() { p.parseRegHalfNumber(opcodes.ROTL) }

func (p *Parser) insJmp()  { p.parseRegNumber(opcodes.JMP) }
func (p *Parser) insJmpa() { p.parseNumber(opcodes.JMPA) }

func (p *Parser) insLdsp() { p.parseInstNoDest(opcodes.LDSPI, opcodes.LDSPA, opcodes.LDSPR) }

func (p *Parser) insPush() { p.parseReg(opcodes.PUSH) }
func (p *Parser) insPop()  { p.parseReg(opcodes.POP) }

func (p *Parser) insCall() { p.parseInstNoImmNoDest(opcodes.CALLA, opcodes.CALLR) }

func (p *Parser) insHalt() { p.parseNoArgs(opcodes.HALT) }
func (p *Parser) insNoop() { p.parseNoArgs(opcodes.NOOP) }
func (p *Parser) insRtn()  { p.parseNoArgs(opcodes.RTN) }

// Common argument parsers

func (p *Parser) parseNoArgs(c byte) {
	p.p.appendCode(c)

	p.expectToken(token.END_INST)
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

	p.expectToken(token.END_INST)
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

	p.expectToken(token.END_INST)
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

	p.expectToken(token.END_INST)
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
	if !p.curTokenIs(token.IMMEDIATE) {
		p.tokenErr(token.IMMEDIATE)
		return
	}

	p.readToken()
	val, ok := p.parseAddress(2)
	if !ok {
		return
	}

	if val > 255 {
		p.parseErr("rotation number too large, must be 0-255")
	}

	// Write code
	p.p.appendCode(c, reg, uint8(val))

	p.expectToken(token.END_INST)
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

	p.expectToken(token.END_INST)
}
