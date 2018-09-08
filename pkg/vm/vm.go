package vm

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"os"

	"github.com/lfkeitel/asml-sim/pkg/token"
)

const (
	numOfMemoryCells = 65536
	numOfRegisters   = 10

	// Double width
	regA = 0xA
	regB = 0xB
	regC = 0xC
	regD = 0xD

	// Quad width
	regE = 0xE
	regF = 0xF
)

func isDoubleReg(r uint8) bool {
	return r >= regA && r <= regD
}

func isQuadReg(r uint8) bool {
	return r == regE || r == regF
}

func regWidth(r uint8) uint8 {
	if r <= 9 {
		return 1
	} else if r >= regA && r <= regD {
		return 2
	}
	return 4
}

func checkRegWidth(regs ...uint8) {
	w := regWidth(regs[0])
	for _, r := range regs {
		if regWidth(r) != w {
			panic("register widths don't match")
		}
	}
}

type VM struct {
	registers  []uint8
	memory     []uint8
	pc         uint16
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
	flags      *token.Flags
}

func New(code []uint8, printState bool) *VM {
	if len(code) < 2 {
		fmt.Println("No code given")
		os.Exit(1)
	}

	flags, valid := token.FlagsFromBytes([]byte{code[0], code[1]})
	if valid {
		code = code[2:]
	}

	if len(code) > numOfMemoryCells-1 { // Reserve printer cell
		fmt.Println("Program too big")
		os.Exit(1)
	}

	newvm := &VM{
		registers:  make([]uint8, numOfRegisters),
		memory:     make([]uint8, numOfMemoryCells),
		pc:         0,
		printState: printState,
		flags:      flags,
	}

	for i, c := range code {
		newvm.memory[i] = c
	}

	return newvm
}

func (vm *VM) Output() []byte {
	return vm.output.Bytes()
}

func (vm *VM) Run(out io.Writer) error {
mainLoop:
	for {
		opcode := vm.fetchByte()

		if vm.printState {
			vm.PrintState()
			fmt.Println(vm.output.String())
			vm.output.Reset()
			// time.Sleep(1 * time.Second)
		}

		switch opcode {
		case token.NOOP:
			// noop
		case token.LOADA:
			vm.writeStateMessage("Instr: LOADA\n")
			vm.loadFromMem(vm.fetchByte(), vm.fetchUint16())
		case token.LOADI:
			vm.writeStateMessage("Instr: LOADI\n")
			vm.loadIntoReg(vm.fetchByte(), vm.fetchUint16())
		case token.STRA:
			vm.writeStateMessage("Instr: STRA\n")
			vm.storeRegInMemory(vm.fetchByte(), vm.fetchUint16())
		case token.MOVR:
			vm.writeStateMessage("Instr: MOVE\n")
			vm.moveRegisters(vm.fetchByte(), vm.fetchByte())
		case token.ADD:
			vm.writeStateMessage("Instr: ADD\n")
			vm.addCompliment(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.OR:
			vm.writeStateMessage("Instr: OR\n")
			vm.orRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.AND:
			vm.writeStateMessage("Instr: AND\n")
			vm.andRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.XOR:
			vm.writeStateMessage("Instr: XOR\n")
			vm.xorRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.ROT:
			vm.writeStateMessage("Instr: ROTATE\n")
			vm.rotateRegister(vm.fetchByte(), vm.fetchByte())
		case token.JMP:
			vm.writeStateMessage("Instr: JUMP\n")
			vm.jumpEq(vm.fetchByte(), vm.fetchUint16())
		case token.HALT:
			vm.writeStateMessage("Instr: HALT\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop
		case token.STRR:
			vm.writeStateMessage("Instr: STORER\n")
			vm.storeRegInMemoryAddr(vm.fetchByte(), vm.fetchByte())
		case token.LOADR:
			vm.writeStateMessage("Instr: LOADR\n")
			vm.loadRegInMemoryAddr(vm.fetchByte(), vm.fetchByte())
		case token.BREAK:
			// NOOP
		default:
			vm.writeString("INVALID OPCODE\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop
		}

		// Print character in memory address FF and reset it to 0
		if vm.memory[numOfMemoryCells-1] > 0 {
			vm.printer.WriteByte(byte(vm.memory[numOfMemoryCells-1]))
			vm.memory[numOfMemoryCells-1] = 0
		}
	}

	out.Write(vm.output.Bytes())

	return nil
}

func (vm *VM) fetchByte() byte {
	b1 := vm.memory[vm.pc]
	vm.pc++
	return b1
}

func (vm *VM) fetchUint16() uint16 {
	b1 := uint16(vm.fetchByte())
	b2 := uint16(vm.fetchByte())
	return (b1 << 8) + b2
}

func (vm *VM) writeStateMessage(s string) {
	if vm.printState {
		vm.writeString(s)
	}
}

// PrintState prints all values in the registers and memory cells
func (vm *VM) PrintState() {
	vm.writeString("Registers   0  1  2  3  4  5  6  7  8  9\n")
	vm.writeString("           ")
	for _, val := range vm.registers {
		vm.writeString(formatHex(val) + " ")
	}

	vm.printMemory16Bit()

	vm.writeString("\nProgram Counter  = ")
	vm.writeString(formatHex16(vm.pc - 1))
	vm.writeString("\n\n")
}

func (vm *VM) printMemory16Bit() {
	vm.writeString("\n\nMemory     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F  10 11 12 13 14 15 16 17 18 19 1A 1B 1C 1D 1E 1F\n\n")
	for i := 0; i < 256; i = i + 32 {
		vm.writeString(formatHex16(uint16(i)))
		vm.writeString("       ")

		for j := 0; j < 16; j++ {
			vm.writeString(formatHex(vm.memory[i+j]))
			vm.writeString(" ")
		}
		vm.writeString(" ")

		for j := 16; j < 32; j++ {
			vm.writeString(formatHex(vm.memory[i+j]))
			vm.writeString(" ")
		}

		vm.writeString("\n")
	}
}

func (vm *VM) writeString(s string) {
	vm.output.WriteString(s)
}

func (vm *VM) writePrinter() {
	vm.output.Write(vm.printer.Bytes())
}

// Format a uint8 as a hex number with leading zeros
func formatHex(num uint8) string {
	return fmt.Sprintf("%02X", num)
}

func formatHex16(num uint16) string {
	return fmt.Sprintf("%04X", num)
}

// width can be 1, 2, or 4
// return value will be 8, 16, or 32-bit depending on width
func (vm *VM) readMem(addr uint16, width int) uint32 {
	switch width {
	case 1:
		return uint32(vm.memory[addr])
	case 2:
		b1 := uint16(vm.memory[addr])
		b2 := uint16(vm.memory[addr+1])
		return uint32((b1 << 8) + b2)
	case 4:
		b1 := uint32(vm.memory[addr])
		b2 := uint32(vm.memory[addr+1])
		b3 := uint32(vm.memory[addr+2])
		b4 := uint32(vm.memory[addr+3])
		return (b1 << 24) + (b2 << 16) + (b3 << 8) + b4
	}
	return 0
}

func (vm *VM) writeMem(addr uint16, width int, val uint32) {
	switch width {
	case 1:
		vm.memory[addr] = uint8(val)
	case 2:
		vm.memory[addr] = uint8(val >> 8)
		vm.memory[addr+1] = uint8(val)
	case 4:
		vm.memory[addr] = uint8(val >> 24)
		vm.memory[addr+1] = uint8(val >> 16)
		vm.memory[addr+2] = uint8(val >> 8)
		vm.memory[addr+3] = uint8(val)
	}
}

func (vm *VM) writeAnyReg(r uint8, v uint32) {
	switch {
	case isQuadReg(r):
		vm.writeQuadReg(r, v)
	case isDoubleReg(r):
		vm.writeDoubleReg(r, uint16(v))
	default:
		vm.writeReg(r, uint8(v))
	}
}

func (vm *VM) readAnyReg(r uint8) uint32 {
	switch {
	case isQuadReg(r):
		return vm.readQuadReg(r)
	case isDoubleReg(r):
		return uint32(vm.readDoubleReg(r))
	default:
		return uint32(vm.readReg(r))
	}
}

func (vm *VM) readAnyReg2Comp(r uint8) uint32 {
	switch {
	case isQuadReg(r):
		return vm.readQuadReg(r)
	case isDoubleReg(r):
		return uint32(int16(vm.readDoubleReg(r)))
	default:
		return uint32(int8(vm.readReg(r)))
	}
}

func (vm *VM) writeReg(r, v uint8) {
	vm.registers[r] = v
}

func (vm *VM) readReg(r uint8) uint8 {
	return vm.registers[r]
}

func (vm *VM) writeDoubleReg(r uint8, v uint16) {
	if r == regA {
		vm.registers[2] = uint8(v >> 8)
		vm.registers[3] = uint8(v)
	} else if r == regB {
		vm.registers[4] = uint8(v >> 8)
		vm.registers[5] = uint8(v)
	} else if r == regC {
		vm.registers[6] = uint8(v >> 8)
		vm.registers[7] = uint8(v)
	} else if r == regD {
		vm.registers[8] = uint8(v >> 8)
		vm.registers[9] = uint8(v)
	}
}

func (vm *VM) readDoubleReg(r uint8) uint16 {
	if r == regA {
		return (uint16(vm.registers[2]) << 8) + uint16(vm.registers[3])
	} else if r == regB {
		return (uint16(vm.registers[4]) << 8) + uint16(vm.registers[5])
	} else if r == regC {
		return (uint16(vm.registers[6]) << 8) + uint16(vm.registers[7])
	} else if r == regD {
		return (uint16(vm.registers[8]) << 8) + uint16(vm.registers[9])
	}
	return 0
}

func (vm *VM) writeQuadReg(r uint8, v uint32) {
	if r == regE {
		vm.registers[2] = uint8(v >> 24)
		vm.registers[3] = uint8(v >> 16)
		vm.registers[4] = uint8(v >> 8)
		vm.registers[5] = uint8(v)
	} else if r == regF {
		vm.registers[6] = uint8(v >> 24)
		vm.registers[7] = uint8(v >> 16)
		vm.registers[8] = uint8(v >> 8)
		vm.registers[9] = uint8(v)
	}
}

func (vm *VM) readQuadReg(r uint8) uint32 {
	if r == regE {
		return (uint32(vm.registers[2]) << 24) + (uint32(vm.registers[3]) << 16) + (uint32(vm.registers[4]) << 8) + uint32(vm.registers[5])
	} else if r == regF {
		return (uint32(vm.registers[6]) << 24) + (uint32(vm.registers[7]) << 16) + (uint32(vm.registers[8]) << 8) + uint32(vm.registers[9])
	}
	return 0
}

// Opcode definitions

func (vm *VM) loadFromMem(r uint8, x uint16) {
	switch {
	case isQuadReg(r):
		vm.writeQuadReg(r, vm.readMem(x, 4))
	case isDoubleReg(r):
		vm.writeDoubleReg(r, uint16(vm.readMem(x, 2)))
	default:
		vm.writeReg(r, uint8(vm.readMem(x, 1)))
	}
}

func (vm *VM) loadIntoReg(r uint8, x uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, x)
	default:
		vm.writeReg(r, uint8(x))
	}
}

func (vm *VM) storeRegInMemory(r uint8, x uint16) {
	switch {
	case isQuadReg(r):
		vm.writeMem(x, 4, vm.readQuadReg(r))
	case isDoubleReg(r):
		vm.writeMem(x, 2, uint32(vm.readDoubleReg(r)))
	default:
		vm.writeMem(x, 1, uint32(vm.readReg(r)))
	}
}

func (vm *VM) moveRegisters(r, s uint8) {
	vm.writeAnyReg(r, vm.readAnyReg(s))
}

func (vm *VM) addCompliment(r, s, t uint8) {
	sv := vm.readAnyReg2Comp(s)
	tv := vm.readAnyReg2Comp(t)
	vm.writeAnyReg(r, uint32(int32(sv)+int32(tv)))
}

func (vm *VM) orRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint32(int32(sv)|int32(tv)))
}

func (vm *VM) andRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint32(int32(sv)&int32(tv)))
}

func (vm *VM) xorRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint32(int32(sv)^int32(tv)))
}

func (vm *VM) rotateRegister(r, x uint8) {
	switch {
	case isQuadReg(r):
		vm.writeQuadReg(r, bits.RotateLeft32(vm.readQuadReg(r), int(-x)))
	case isDoubleReg(r):
		vm.writeDoubleReg(r, bits.RotateLeft16(vm.readDoubleReg(r), int(-x)))
	default:
		vm.writeReg(r, bits.RotateLeft8(vm.readReg(r), int(-x)))
	}
}

func (vm *VM) jumpEq(r uint8, d uint16) {
	if uint8(vm.readAnyReg(r)) == vm.readReg(0) {
		vm.pc = d
	}
}

func (vm *VM) storeRegInMemoryAddr(d, s uint8) {
	addr := uint16(vm.readAnyReg(d))
	vm.storeRegInMemory(s, addr)
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	addr := uint16(vm.readAnyReg(s))
	vm.loadFromMem(d, addr)
}
