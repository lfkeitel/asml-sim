package vm

func isDoubleReg(r uint8) bool {
	return r >= regA && r <= regD
}

func regWidth(r uint8) uint8 {
	if r <= 9 {
		return 1
	} else if isDoubleReg(r) {
		return 2
	}
	return 0
}

func checkRegWidth(regs ...uint8) {
	w := regWidth(regs[0])
	for _, r := range regs {
		if regWidth(r) != w {
			panic("register widths don't match")
		}
	}
}

// width can be 1 or 2
// return value will be 8 or 16-bit depending on width
func (vm *VM) readMem(addr uint16, width int) uint16 {
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
	return uint8(vm.readMem(addr, 1))
}

func (vm *VM) readMem16(addr uint16) uint16 {
	return vm.readMem(addr, 2)
}

// width can be 1 or 2
func (vm *VM) writeMem(addr uint16, width int, val uint16) {
	switch width {
	case 1:
		vm.memory[addr] = uint8(val)
	case 2:
		vm.memory[addr] = uint8(val >> 8)
		vm.memory[addr+1] = uint8(val)
	}
}

func (vm *VM) writeMem8(addr uint16, val uint8) {
	vm.writeMem(addr, 1, uint16(val))
}

func (vm *VM) writeMem16(addr uint16, val uint16) {
	vm.writeMem(addr, 2, uint16(val))
}

func (vm *VM) writeAnyReg(r uint8, v uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, v)
	default:
		vm.writeReg(r, uint8(v))
	}
}

func (vm *VM) readAnyReg(r uint8) uint16 {
	switch {
	case isDoubleReg(r):
		return vm.readDoubleReg(r)
	default:
		return uint16(vm.readReg(r))
	}
}

func (vm *VM) readAnyReg2Comp(r uint8) uint16 {
	switch {
	case isDoubleReg(r):
		return vm.readDoubleReg(r)
	default:
		return uint16(int8(vm.readReg(r)))
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
