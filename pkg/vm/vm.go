package vm

import (
	"bytes"
	"fmt"
	"io"
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

type VM struct {
	registers  []uint8
	memory     []uint8
	pc, sp     uint16
	output     bytes.Buffer
	printer    bytes.Buffer
	printState bool
}

func New(code []uint8, printState bool) *VM {
	if len(code) < 2 {
		fmt.Println("No code given")
		os.Exit(1)
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
		case token.ADDI:
			vm.writeStateMessage("Instr: ADDI\n")
			vm.addImmCompliment(vm.fetchByte(), vm.fetchByte(), vm.fetchByte())
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
		case token.JMPA:
			vm.writeStateMessage("Instr: JUMPA\n")
			vm.jumpAbs(vm.fetchUint16())
		case token.LDSP:
			vm.writeStateMessage("Instr: LDSP\n")
			vm.loadSP(vm.fetchUint16())
		case token.LDSPI:
			vm.writeStateMessage("Instr: LDSPI\n")
			vm.loadSPImm(vm.fetchUint16())
		case token.PUSH:
			vm.writeStateMessage("Instr: PUSH\n")
			vm.push(vm.fetchByte())
		case token.POP:
			vm.writeStateMessage("Instr: POP\n")
			vm.pop(vm.fetchByte())
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
