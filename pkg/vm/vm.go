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
	numOfRegisters   = 16
)

type VM struct {
	registers  []uint16
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
		registers:  make([]uint16, numOfRegisters),
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
		}

		switch opcode {
		case token.NOOP:
			// noop
		case token.LOADA:
			vm.writeStateMessage("LOADA\n")
			vm.loadFromMem(vm.fetchByte(), vm.fetchUint16())
		case token.LOADI:
			vm.writeStateMessage("LOADI\n")
			vm.loadIntoReg(vm.fetchByte(), vm.fetchUint16())
		case token.STRA:
			vm.writeStateMessage("STRA\n")
			vm.storeRegInMemory(vm.fetchByte(), vm.fetchUint16())
		case token.MOVR:
			vm.writeStateMessage("MOVE\n")
			vm.moveRegisters(vm.fetchByte(), vm.fetchByte())
		case token.ADD:
			vm.writeStateMessage("ADD\n")
			vm.addCompliment(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.OR:
			vm.writeStateMessage("OR\n")
			vm.orRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.AND:
			vm.writeStateMessage("AND\n")
			vm.andRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.XOR:
			vm.writeStateMessage("XOR\n")
			vm.xorRegisters(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
		case token.ROT:
			vm.writeStateMessage("ROTATE\n")
			vm.rotateRegister(vm.fetchByte(), vm.fetchByte())
		case token.JMP:
			vm.writeStateMessage("JUMP\n")
			vm.jumpEq(vm.fetchByte(), vm.fetchUint16())
		case token.HALT:
			vm.writeStateMessage("HALT\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop
		case token.STRR:
			vm.writeStateMessage("STORER\n")
			vm.storeRegInMemoryAddr(vm.fetchByte(), vm.fetchByte())
		case token.LOADR:
			vm.writeStateMessage("LOADR\n")
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

		if vm.printState {
			vm.writeString("\n")
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
	vm.writeString("Registers     0    1    2    3    4    5    6    7    8    9    A    B    C    D    E    F\n")
	vm.writeString("           ")
	for _, val := range vm.registers {
		vm.writeString(formatHex16(val) + " ")
	}

	vm.printMemory16Bit()

	vm.writeString("\nProgram Counter  = ")
	vm.writeString(formatHex16(vm.pc - 1))

	vm.writeString("\nNext Instruction = N/A ")
	// vm.writeString(formatHex(vm.memory[vm.pc]))
	// vm.writeString(" ")
	// vm.writeString(formatHex(vm.memory[vm.pc+1]))
	// vm.writeString(" ")
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

// Opcode definitions

func (vm *VM) loadFromMem(r uint8, x uint16) {
	b1 := uint16(vm.memory[x])
	b2 := uint16(vm.memory[x+1])
	vm.registers[r] = (b1 << 8) + b2
}

func (vm *VM) loadIntoReg(r uint8, x uint16) {
	vm.registers[r] = x
}

func (vm *VM) storeRegInMemory(r uint8, x uint16) {
	vm.memory[x] = uint8(vm.registers[r] >> 8)
	vm.memory[x+1] = uint8(vm.registers[r])
}

func (vm *VM) moveRegisters(r, s uint8) {
	vm.registers[r] = vm.registers[s]
}

func (vm *VM) addCompliment(r, s, t uint8) {
	vm.registers[r] = uint16(int16(vm.registers[s]) + int16(vm.registers[t]))
}

func (vm *VM) orRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] | vm.registers[t]
}

func (vm *VM) andRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] & vm.registers[t]
}

func (vm *VM) xorRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] ^ vm.registers[t]
}

func (vm *VM) rotateRegister(r, x uint8) {
	vm.registers[r] = bits.RotateLeft16(vm.registers[r], int(-x))
}

func (vm *VM) jumpEq(r uint8, d uint16) {
	if vm.registers[r] == vm.registers[0] {
		vm.pc = d
	}
}

func (vm *VM) storeRegInMemoryAddr(d, s uint8) {
	vm.storeRegInMemory(s, vm.registers[d])
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	vm.loadFromMem(d, vm.registers[s])
}
