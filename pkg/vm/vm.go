package vm

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"os"

	"github.com/lfkeitel/asml-sim/pkg/lexer"
)

const (
	numOfMemoryCells = 256
	numOfRegisters   = 16
)

type VM struct {
	registers  []uint8
	memory     []uint8
	pc         uint8
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
}

func New(code []uint8, printState bool) *VM {
	if len(code) > numOfMemoryCells-1 { // Reserve printer cell
		fmt.Println("Program too big")
		os.Exit(0)
	}

	newvm := &VM{
		registers:  make([]uint8, numOfRegisters),
		memory:     make([]uint8, numOfMemoryCells),
		pc:         0,
		printState: printState,
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
		case lexer.NOOP:
			// noop
		case lexer.LOADA:
			vm.writeStateMessage("LOADA\n")
			vm.loadFromMem(operand1, operand2, operand3)
		case lexer.LOADI:
			vm.writeStateMessage("LOADI\n")
			vm.loadIntoReg(operand1, operand2, operand3)
		case lexer.STRA:
			if operand2 == 15 && operand3 == 15 {
				vm.writeStateMessage("PRINT\n")
			} else {
				vm.writeStateMessage("STOREA\n")
			}
			vm.storeRegInMemory(operand1, operand2, operand3)
		case lexer.MOVR:
			vm.writeStateMessage("MOVE\n")
			vm.moveRegisters(operand2, operand3)
		case lexer.ADD:
			vm.writeStateMessage("ADD\n")
			vm.addCompliment(operand1, operand2, operand3)
		case lexer.OR:
			vm.writeStateMessage("OR\n")
			vm.orRegisters(operand1, operand2, operand3)
		case lexer.AND:
			vm.writeStateMessage("AND\n")
			vm.andRegisters(operand1, operand2, operand3)
		case lexer.XOR:
			vm.writeStateMessage("XOR\n")
			vm.xorRegisters(operand1, operand2, operand3)
		case lexer.ROT:
			vm.writeStateMessage("ROTATE\n")
			vm.rotateRegister(operand1, operand3)
		case lexer.JMP:
			vm.writeStateMessage("JUMP\n")
			vm.jumpEq(operand1, operand2, operand3)
		case lexer.HALT:
			vm.writeStateMessage("HALT\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop
		case lexer.STRR:
			vm.writeStateMessage("STORER\n")
			vm.storeRegInMemoryAddr(operand2, operand3)
		case lexer.LOADR:
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
		if vm.memory[255] > 0 {
			vm.printer.WriteByte(byte(vm.memory[255]))
			vm.memory[255] = 0
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

	vm.writeString("\n\nMemory     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n\n")
	for i := 0; i < numOfMemoryCells; i = i + 16 {
		vm.writeString(formatHex(uint8(i)))
		vm.writeString("         ")

		for j := 0; j < 16; j++ {
			vm.writeString(formatHex(vm.memory[i+j]))
			vm.writeString(" ")
		}

		vm.writeString("\n")
	}

	vm.writeString("\nProgram Counter  = ")
	vm.writeString(formatHex(vm.pc))

	vm.writeString("\nNext Instruction = ")
	vm.writeString(formatHex(vm.memory[vm.pc]))
	vm.writeString(" ")
	vm.writeString(formatHex(vm.memory[vm.pc+1]))
	vm.writeString(" ")
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

// Opcode definitions

func (vm *VM) loadFromMem(r, x, y uint8) {
	vm.registers[r] = vm.memory[(x<<4)+y]
}

func (vm *VM) loadIntoReg(r, x, y uint8) {
	vm.registers[r] = (x << 4) + y
}

func (vm *VM) storeRegInMemory(r, x, y uint8) {
	vm.memory[(x<<4)+y] = vm.registers[r]
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

func (vm *VM) jumpEq(r, x, y uint8) {
	if vm.registers[r] == vm.registers[0] {
		vm.pc = ((x << 4) + y) - 2 // The main loop adds 2, compensate
	}
}

func (vm *VM) storeRegInMemoryAddr(r, s uint8) {
	vm.memory[vm.registers[r]] = vm.registers[s]
}

func (vm *VM) loadRegInMemoryAddr(r, s uint8) {
	vm.registers[r] = vm.memory[vm.registers[s]]
}
