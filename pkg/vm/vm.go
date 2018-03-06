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
	numOfRegisters = 16
)

var (
	memoryCellCount = [...]int{
		256,
		65536,
	}
)

type VM struct {
	registers  []uint8
	memory     []uint8
	pc         uint16
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
	flags      *token.Flags
	memorySize int
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

	numOfMemoryCells := memoryCellCount[flags.Size]
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
		memorySize: numOfMemoryCells,
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
		instruction := vm.memory[vm.pc]
		opcode := instruction >> 4          // First 4 bits of byte 1
		operand1 := instruction & 15        // Second 4 bits of byte 1
		operand2 := vm.memory[vm.pc+1] >> 4 // First 4 bits of byte 2
		operand3 := vm.memory[vm.pc+1] & 15 // Second 4 bits of byte 2

		if vm.printState {
			vm.PrintState()
		}

		switch opcode {
		case token.NOOP:
			// noop
		case token.LOADA:
			vm.writeStateMessage("LOADA\n")
			d1 := (operand2 << 4) + operand3
			var d2 uint8
			if vm.flags.Size == token.SixteenBit {
				d2 = vm.memory[vm.pc+2]
				vm.pc++
			}
			vm.loadFromMem(operand1, d1, d2)
		case token.LOADI:
			vm.writeStateMessage("LOADI\n")
			vm.loadIntoReg(operand1, operand2, operand3)
		case token.STRA:
			d1 := (operand2 << 4) + operand3
			var d2 uint8
			if vm.flags.Size == token.SixteenBit {
				d2 = vm.memory[vm.pc+2]
				vm.pc++
			}

			vm.writeStateMessage("STOREA\n")
			vm.storeRegInMemory(operand1, d1, d2)
		case token.MOVR:
			vm.writeStateMessage("MOVE\n")
			vm.moveRegisters(operand2, operand3)
		case token.ADD:
			vm.writeStateMessage("ADD\n")
			vm.addCompliment(operand1, operand2, operand3)
		case token.OR:
			vm.writeStateMessage("OR\n")
			vm.orRegisters(operand1, operand2, operand3)
		case token.AND:
			vm.writeStateMessage("AND\n")
			vm.andRegisters(operand1, operand2, operand3)
		case token.XOR:
			vm.writeStateMessage("XOR\n")
			vm.xorRegisters(operand1, operand2, operand3)
		case token.ROT:
			vm.writeStateMessage("ROTATE\n")
			vm.rotateRegister(operand1, operand3)
		case token.JMP:
			vm.writeStateMessage("JUMP\n")
			d1 := (operand2 << 4) + operand3
			var d2 uint8
			if vm.flags.Size == token.SixteenBit {
				d2 = vm.memory[vm.pc+2]
				vm.pc++
			}
			vm.jumpEq(operand1, d1, d2)
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
			vm.storeRegInMemoryAddr(operand2, operand3)
		case token.LOADR:
			vm.writeStateMessage("LOADR\n")
			vm.loadRegInMemoryAddr(operand2, operand3)
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

		vm.pc += 2

		// Print character in memory address FF and reset it to 0
		if vm.memory[vm.memorySize-1] > 0 {
			vm.printer.WriteByte(byte(vm.memory[vm.memorySize-1]))
			vm.memory[vm.memorySize-1] = 0
		}
	}

	out.Write(vm.output.Bytes())

	return nil
}

func (vm *VM) writeStateMessage(s string) {
	if vm.printState {
		vm.writeString(s)
	}
}

// PrintState prints all values in the registers and memory cells
func (vm *VM) PrintState() {
	vm.writeString("Registers  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n")
	vm.writeString("           ")
	for _, val := range vm.registers {
		vm.writeString(formatHex(val) + " ")
	}

	if vm.flags.Size == token.EightBit {
		vm.printMemory8Bit()
	} else if vm.flags.Size == token.SixteenBit {
		vm.printMemory16Bit()
	}

	vm.writeString("\nProgram Counter  = ")
	vm.writeString(formatHex16(vm.pc))

	vm.writeString("\nNext Instruction = ")
	vm.writeString(formatHex(vm.memory[vm.pc]))
	vm.writeString(" ")
	vm.writeString(formatHex(vm.memory[vm.pc+1]))
	vm.writeString(" ")
}

func (vm *VM) printMemory8Bit() {
	vm.writeString("\n\nMemory     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n\n")
	for i := 0; i < vm.memorySize; i = i + 16 {
		vm.writeString(formatHex(uint8(i)))
		vm.writeString("         ")

		for j := 0; j < 16; j++ {
			vm.writeString(formatHex(vm.memory[i+j]))
			vm.writeString(" ")
		}

		vm.writeString("\n")
	}
}

func (vm *VM) printMemory16Bit() {
	vm.writeString("\n\nMemory     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F  10 11 12 13 14 15 16 17 18 19 1A 1B 1C 1D 1E 1F\n\n")
	for i := 0; i < vm.memorySize; i = i + 32 {
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

func (vm *VM) loadFromMem(r, x, y uint8) {
	if vm.flags.Size == token.EightBit {
		vm.registers[r] = vm.memory[x]
	} else {
		vm.registers[r] = vm.memory[(uint16(x)<<8)+uint16(y)]
	}
}

func (vm *VM) loadIntoReg(r, x, y uint8) {
	vm.registers[r] = (x << 4) + y
}

func (vm *VM) storeRegInMemory(r, x, y uint8) {
	if vm.flags.Size == token.EightBit {
		vm.memory[x] = vm.registers[r]
	} else {
		vm.memory[(uint16(x)<<8)+uint16(y)] = vm.registers[r]
	}
}

func (vm *VM) moveRegisters(r, s uint8) {
	vm.registers[r] = vm.registers[s]
}

func (vm *VM) addCompliment(r, s, t uint8) {
	vm.registers[r] = uint8(int8(vm.registers[s]) + int8(vm.registers[t]))
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
	vm.registers[r] = bits.RotateLeft8(vm.registers[r], int(-x))
}

func (vm *VM) jumpEq(r, d1, d2 uint8) {
	if vm.registers[r] == vm.registers[0] {
		if vm.flags.Size == token.EightBit {
			vm.pc = uint16(d1) - 2 // The main loop adds 2, compensate
		} else {
			vm.pc = (uint16(d1) << 8) + uint16(d2) - 2 // The main loop adds 2, compensate
		}
	}
}

func (vm *VM) storeRegInMemoryAddr(d, s uint8) {
	if vm.flags.Size == token.EightBit {
		vm.memory[vm.registers[d]] = vm.registers[s]
	} else {
		vm.memory[uint16(vm.registers[d])<<8+uint16(vm.registers[d+1])] = vm.registers[s]
	}
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	if vm.flags.Size == token.EightBit {
		vm.registers[d] = vm.memory[vm.registers[s]]
	} else {
		vm.registers[d] = vm.memory[uint16(vm.registers[s])<<8+uint16(vm.registers[s+1])]
	}
}
