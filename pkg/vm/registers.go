package vm

type Register uint8

// Register names
const (
	Register0 Register = iota
	Register1
	Register2
	Register3
	Register4
	Register5
	Register6
	Register7
	Register8
	Register9
	RegisterA
	RegisterB
	RegisterC
	RegisterD
)

func IsDoubleReg(r Register) bool {
	return r >= regA && r <= regD
}

func regWidth(r Register) uint8 {
	if r <= 9 {
		return 1
	} else if IsDoubleReg(r) {
		return 2
	}
	return 0
}

func checkRegWidth(regs ...Register) {
	w := regWidth(regs[0])
	for _, r := range regs {
		if regWidth(r) != w {
			panic("register widths don't match")
		}
	}
}

// width can be 1 or 2
// return value will be 8 or 16-bit depending on width
func (vm *VM) ReadMem(addr uint16, width int) uint16 {
	switch width {
	case 1:
		return uint16(vm.memory[addr])
	case 2:
		b1 := uint16(vm.memory[addr])
		b2 := uint16(vm.memory[addr+1])
		return (b1 << 8) + b2
	}
	return 0
}

func (vm *VM) readMem8(addr uint16) uint8 {
	return uint8(vm.ReadMem(addr, 1))
}

func (vm *VM) readMem16(addr uint16) uint16 {
	return vm.ReadMem(addr, 2)
}

// width can be 1 or 2
func (vm *VM) WriteMem(addr uint16, width int, val uint16) {
	switch width {
	case 1:
		vm.memory[addr] = uint8(val)
	case 2:
		vm.memory[addr] = uint8(val >> 8)
		vm.memory[addr+1] = uint8(val)
	}
}

func (vm *VM) writeMem8(addr uint16, val uint8) {
	vm.WriteMem(addr, 1, uint16(val))
}

func (vm *VM) writeMem16(addr uint16, val uint16) {
	vm.WriteMem(addr, 2, uint16(val))
}

func (vm *VM) WriteReg(r Register, v uint16) {
	switch {
	case IsDoubleReg(r):
		vm.writeDoubleReg(r, v)
	default:
		vm.writeSingleReg(r, uint8(v))
	}
}

func (vm *VM) ReadReg(r Register) uint16 {
	switch {
	case IsDoubleReg(r):
		return vm.readDoubleReg(r)
	default:
		return uint16(vm.readSingleReg(r))
	}
}

func (vm *VM) readAnyReg2Comp(r Register) uint16 {
	switch {
	case IsDoubleReg(r):
		return vm.readDoubleReg(r)
	default:
		return uint16(int8(vm.readSingleReg(r)))
	}
}

func (vm *VM) writeSingleReg(r Register, v uint8) {
	vm.registers[r] = v
}

func (vm *VM) readSingleReg(r Register) uint8 {
	return vm.registers[r]
}

func (vm *VM) writeDoubleReg(r Register, v uint16) {
	switch r {
	case regA:
		vm.registers[2] = uint8(v >> 8)
		vm.registers[3] = uint8(v)
	case regB:
		vm.registers[4] = uint8(v >> 8)
		vm.registers[5] = uint8(v)
	case regC:
		vm.registers[6] = uint8(v >> 8)
		vm.registers[7] = uint8(v)
	case regD:
		vm.registers[8] = uint8(v >> 8)
		vm.registers[9] = uint8(v)
	}
}

func (vm *VM) readDoubleReg(r Register) uint16 {
	switch r {
	case regA:
		return (uint16(vm.registers[2]) << 8) + uint16(vm.registers[3])
	case regB:
		return (uint16(vm.registers[4]) << 8) + uint16(vm.registers[5])
	case regC:
		return (uint16(vm.registers[6]) << 8) + uint16(vm.registers[7])
	case regD:
		return (uint16(vm.registers[8]) << 8) + uint16(vm.registers[9])
	}
	return 0
}
