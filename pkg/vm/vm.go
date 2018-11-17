package vm

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/lfkeitel/asml-sim/pkg/opcodes"
	"github.com/lfkeitel/asml-sim/pkg/parser"
)

const (
	numOfMemoryCells = 65536
	numOfRegisters   = 10

	// Double width
	regA = 0xA
	regB = 0xB
	regC = 0xC
	regD = 0xD
)

type VM struct {
	registers  []uint8
	memory     []uint8
	pc, sp     uint16
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
}

func New(code []parser.CodePart, printState bool) *VM {
	if len(code) == 0 {
		fmt.Println("No code given")
		os.Exit(1)
	}

	newvm := &VM{
		registers:  make([]uint8, numOfRegisters),
		memory:     make([]uint8, numOfMemoryCells),
		printState: printState,
	}

	for _, c := range code {
		pc := c.StartPC

		overflow := false
		for i, b := range c.Bytes {
			if overflow {
				fmt.Println("Code overflowed past address 0xFFFF")
				os.Exit(1)
			}

			loc := uint16(i) + pc
			newvm.memory[loc] = b
			if loc == numOfMemoryCells-1 {
				overflow = true
			}
		}
	}

	newvm.Reset()

	return newvm
}

func (vm *VM) Reset() {
	vm.pc = (uint16(vm.memory[0xFFFE]) << 8) | uint16(vm.memory[0xFFFF])
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
		}

		switch opcode {
		case opcodes.NOOP:
			// noop

		case opcodes.LOADI:
			vm.writeStateMessage("Instr: LOADI\n")
			vm.loadIntoReg(vm.fetchByte(), vm.fetchUint16())
		case opcodes.LOADA:
			vm.writeStateMessage("Instr: LOADA\n")
			vm.loadFromMem(vm.fetchByte(), vm.fetchUint16())
		case opcodes.LOADR:
			vm.writeStateMessage("Instr: LOADR\n")
			vm.loadRegInMemoryAddr(vm.fetchByte(), vm.fetchByte())

		case opcodes.STRA:
			vm.writeStateMessage("Instr: STRA\n")
			vm.storeRegToMemory(vm.fetchByte(), vm.fetchUint16())
		case opcodes.STRR:
			vm.writeStateMessage("Instr: STRR\n")
			vm.storeRegToRegAddr(vm.fetchByte(), vm.fetchByte())

		case opcodes.XFER:
			vm.writeStateMessage("Instr: XFER\n")
			vm.xferRegisters(vm.fetchByte(), vm.fetchByte())

		case opcodes.ADDA:
			vm.writeStateMessage("Instr: ADDA\n")
			vm.addAddr(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ADDI:
			vm.writeStateMessage("Instr: ADDI\n")
			vm.addImm(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ADDR:
			vm.writeStateMessage("Instr: ADDR\n")
			vm.addReg(vm.fetchByte(), vm.fetchByte())

		case opcodes.ORA:
			vm.writeStateMessage("Instr: ORA\n")
			vm.orAddr(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ORI:
			vm.writeStateMessage("Instr: ORI\n")
			vm.orImm(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ORR:
			vm.writeStateMessage("Instr: ORR\n")
			vm.orReg(vm.fetchByte(), vm.fetchByte())

		case opcodes.ANDA:
			vm.writeStateMessage("Instr: ANDA\n")
			vm.andAddr(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ANDI:
			vm.writeStateMessage("Instr: ANDI\n")
			vm.andImm(vm.fetchByte(), vm.fetchUint16())
		case opcodes.ANDR:
			vm.writeStateMessage("Instr: ANDR\n")
			vm.andReg(vm.fetchByte(), vm.fetchByte())

		case opcodes.XORA:
			vm.writeStateMessage("Instr: XORA\n")
			vm.xorAddr(vm.fetchByte(), vm.fetchUint16())
		case opcodes.XORI:
			vm.writeStateMessage("Instr: XORI\n")
			vm.xorImm(vm.fetchByte(), vm.fetchUint16())
		case opcodes.XORR:
			vm.writeStateMessage("Instr: XORR\n")
			vm.xorReg(vm.fetchByte(), vm.fetchByte())

		case opcodes.ROTR:
			vm.writeStateMessage("Instr: ROTR\n")
			vm.rotrRegister(vm.fetchByte(), vm.fetchByte())
		case opcodes.ROTL:
			vm.writeStateMessage("Instr: ROTL\n")
			vm.rotlRegister(vm.fetchByte(), vm.fetchByte())

		case opcodes.JMP:
			vm.writeStateMessage("Instr: JMP\n")
			vm.jumpEq(vm.fetchByte(), vm.fetchUint16())
		case opcodes.JMPA:
			vm.writeStateMessage("Instr: JMPA\n")
			vm.jumpAbs(vm.fetchUint16())

		case opcodes.HALT:
			vm.writeStateMessage("Instr: HALT\n")
			if vm.printState {
				vm.writeString("\nPrinter: ")
			}
			vm.writePrinter()
			vm.writeString("\n")
			break mainLoop

		case opcodes.LDSPA:
			vm.writeStateMessage("Instr: LDSPA\n")
			vm.loadSPAddr(vm.fetchUint16())
		case opcodes.LDSPI:
			vm.writeStateMessage("Instr: LDSPI\n")
			vm.loadSPImm(vm.fetchUint16())
		case opcodes.LDSPR:
			vm.writeStateMessage("Instr: LDSPR\n")
			vm.loadSPReg(vm.fetchByte())

		case opcodes.PUSH:
			vm.writeStateMessage("Instr: PUSH\n")
			vm.push(vm.fetchByte())
		case opcodes.POP:
			vm.writeStateMessage("Instr: POP\n")
			vm.pop(vm.fetchByte())

		case opcodes.CALLA:
			vm.writeStateMessage("Instr: CALLA\n")
			vm.calla(vm.fetchUint16())
		case opcodes.CALLR:
			vm.writeStateMessage("Instr: CALLR\n")
			vm.callr(vm.fetchByte())

		case opcodes.RTN:
			vm.writeStateMessage("Instr: RTN\n")
			vm.rtn()

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
		if vm.memory[numOfMemoryCells-3] > 0 {
			vm.printer.WriteByte(byte(vm.memory[numOfMemoryCells-3]))
			vm.memory[numOfMemoryCells-3] = 0
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
