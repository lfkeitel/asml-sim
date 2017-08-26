package main

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"
)

const (
	numOfMemoryCells = 256
	numOfRegisters   = 16
)

type vm struct {
	registers  []uint8
	memory     []uint8
	pc         uint8
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
}

func newVM(code []uint8, disableState bool) *vm {
	if len(code) > 255 { // 256 - 1 for printer cell
		panic("Program too big")
	}

	newvm := &vm{
		registers:  make([]uint8, numOfRegisters),
		memory:     make([]uint8, numOfMemoryCells),
		pc:         0,
		printState: !disableState,
	}

	for i, c := range code {
		newvm.memory[i] = c
	}

	return newvm
}

func (vm *vm) run(out io.Writer) error {
mainLoop:
	for {
		instruction := vm.memory[vm.pc]
		opcode := instruction >> 4
		operand1 := instruction & 15
		operand2 := vm.memory[vm.pc+1] >> 4
		operand3 := vm.memory[vm.pc+1] & 15

		if vm.printState {
			vm.printVMState()
		}

		switch opcode {
		case 1:
			vm.writeStateMessage("LOAD\n")
			vm.loadFromMem(operand1, operand2, operand3)
		case 2:
			vm.writeStateMessage("LOAD\n")
			vm.loadIntoReg(operand1, operand2, operand3)
		case 3:
			vm.writeStateMessage("STORE\n")
			vm.storeRegInMemory(operand1, operand2, operand3)
		case 4:
			vm.writeStateMessage("MOVE\n")
			vm.moveRegisters(operand2, operand3)
		case 5:
			vm.writeStateMessage("ADD\n")
			vm.addCompliment(operand1, operand2, operand3)
		case 6:
			vm.writeStateMessage("ADD\n")
			vm.addFloats(operand1, operand2, operand3)
		case 7:
			vm.writeStateMessage("OR\n")
			vm.orRegisters(operand1, operand2, operand3)
		case 8:
			vm.writeStateMessage("AND\n")
			vm.andRegisters(operand1, operand2, operand3)
		case 9:
			vm.writeStateMessage("XOR\n")
			vm.xorRegisters(operand1, operand2, operand3)
		case 10:
			vm.writeStateMessage("ROTATE\n")
			vm.rotateRegister(operand1, operand3)
		case 11:
			vm.writeStateMessage("JUMP\n")
			vm.jumpEq(operand1, operand2, operand3)
		case 12:
			vm.writeStateMessage("HALT\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop
		default:
			vm.writeString("INVALID OPCODE\n")
		}

		if vm.printState {
			vm.writeString("\n")
		}

		vm.pc += 2

		if vm.memory[255] > 0 {
			vm.printer.WriteByte(byte(vm.memory[255]))
			vm.memory[255] = 0
		}
	}

	out.Write(vm.output.Bytes())

	return nil
}

func (vm *vm) writeStateMessage(s string) {
	if vm.printState {
		vm.writeString(s)
	}
}

// printVMState prints all values in the registers and memory cells
func (vm *vm) printVMState() {
	vm.writeString("Registers  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n")
	vm.writeString("           ")
	for _, val := range vm.registers {
		vm.writeString(formatHex(val) + " ")
	}

	vm.writeString("\n\nMemory     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n")
	for i := 0; i < numOfMemoryCells; i = i + 16 {
		vm.writeString(formatHex(uint8(i)))
		vm.writeString("         ")

		for j := 0; j < 16; j++ {
			vm.writeString(formatHex(vm.memory[i+j]))
			vm.writeString(" ")
		}

		vm.writeString("\n")
	}

	vm.writeString("\nProgram Counter      = ")
	vm.writeString(formatHex(vm.pc))
	vm.writeString("\nInstruction Register = ")
	vm.writeString(formatHex(vm.memory[vm.pc]))
	vm.writeString(formatHex(vm.memory[vm.pc+1]))
	vm.writeString(" ")
}

func (vm *vm) writeString(s string) {
	vm.output.WriteString(s)
}

func (vm *vm) writePrinter() {
	vm.output.Write(vm.printer.Bytes())
}

// Format a uint8 as a hex number with leading zeros
func formatHex(num uint8) string {
	return fmt.Sprintf("%02X", num)
}

// Opcode definitions

func (vm *vm) loadFromMem(r, x, y uint8) {
	vm.registers[r] = vm.memory[(x<<4)+y]
}

func (vm *vm) loadIntoReg(r, x, y uint8) {
	vm.registers[r] = (x << 4) + y
}

func (vm *vm) storeRegInMemory(r, x, y uint8) {
	vm.memory[(x<<4)+y] = vm.registers[r]
}

func (vm *vm) moveRegisters(r, s uint8) {
	vm.registers[s] = vm.registers[r]
}

func (vm *vm) addCompliment(r, s, t uint8) {
	vm.registers[r] = uint8(int8(vm.registers[s]) + int8(vm.registers[t]))
}

func (vm *vm) addFloats(r, s, t uint8) {}

func (vm *vm) orRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] | vm.registers[t]
}

func (vm *vm) andRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] & vm.registers[t]
}

func (vm *vm) xorRegisters(r, s, t uint8) {
	vm.registers[r] = vm.registers[s] ^ vm.registers[t]
}

func (vm *vm) rotateRegister(r, x uint8) {
	vm.registers[r] = bits.RotateLeft8(vm.registers[r], int(-x))
}

func (vm *vm) jumpEq(r, x, y uint8) {
	if vm.registers[r] == vm.registers[0] {
		vm.pc = ((x << 4) + y) - 2 // The main loop adds 2, compensate
	}
}
