package vm

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

func (vm *VM) readMem8(addr uint16) uint8 {
	return uint8(vm.readMem(addr, 1))
}

func (vm *VM) readMem16(addr uint16) uint16 {
	return uint16(vm.readMem(addr, 2))
}

func (vm *VM) readMem32(addr uint16) uint32 {
	return vm.readMem(addr, 4)
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

func (vm *VM) writeMem8(addr uint16, val uint8) {
	vm.writeMem(addr, 1, uint32(val))
}

func (vm *VM) writeMem16(addr uint16, val uint16) {
	vm.writeMem(addr, 2, uint32(val))
}

func (vm *VM) writeMem32(addr uint16, val uint32) {
	vm.writeMem(addr, 4, val)
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

func (vm *VM) readDoubleReg(r uint8) uint16 {
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
