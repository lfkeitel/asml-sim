package vm

import (
	"fmt"
)

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
	vm.writeString("\nStack Pointer  = ")
	vm.writeString(formatHex16(vm.sp))
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
