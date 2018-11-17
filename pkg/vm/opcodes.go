package vm

import (
	"math/bits"
)

// Opcode definitions

func (vm *VM) loadFromMem(r uint8, x uint16) {
	switch {
	case IsDoubleReg(Register(r)):
		vm.writeDoubleReg(Register(r), uint16(vm.ReadMem(x, 2)))
	default:
		vm.writeSingleReg(Register(r), uint8(vm.ReadMem(x, 1)))
	}
}

func (vm *VM) loadIntoReg(r uint8, x uint16) {
	switch {
	case IsDoubleReg(Register(r)):
		vm.writeDoubleReg(Register(r), x)
	default:
		vm.writeSingleReg(Register(r), uint8(x))
	}
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	addr := vm.ReadReg(Register(s))
	vm.loadFromMem(d, addr)
}

func (vm *VM) storeRegToMemory(r uint8, x uint16) {
	switch {
	case IsDoubleReg(Register(r)):
		vm.WriteMem(x, 2, vm.readDoubleReg(Register(r)))
	default:
		vm.WriteMem(x, 1, uint16(vm.readSingleReg(Register(r))))
	}
}

func (vm *VM) storeRegToRegAddr(s, d uint8) {
	addr := vm.ReadReg(Register(d))
	vm.storeRegToMemory(s, addr)
}

func (vm *VM) xferRegisters(r, s uint8) {
	vm.WriteReg(Register(r), vm.ReadReg(Register(s)))
}

func (vm *VM) addAddr(d uint8, s uint16) {
	dr := Register(d)
	switch {
	case IsDoubleReg(dr):
		vm.writeDoubleReg(dr, vm.ReadReg(dr)+vm.readMem16(s))
	default:
		vm.writeSingleReg(dr, uint8(vm.ReadReg(dr))+vm.readMem8(s))
	}
}

func (vm *VM) addImm(d uint8, s uint16) {
	dr := Register(d)
	vm.WriteReg(dr, vm.ReadReg(dr)+s)
}

func (vm *VM) addReg(d uint8, s uint8) {
	vm.WriteReg(Register(d), vm.ReadReg(Register(s))+vm.ReadReg(Register(d)))
}

func (vm *VM) orAddr(d uint8, s uint16) {
	dr := Register(d)
	switch {
	case IsDoubleReg(dr):
		vm.writeDoubleReg(dr, vm.ReadReg(dr)|vm.readMem16(s))
	default:
		vm.writeSingleReg(dr, uint8(vm.ReadReg(dr))|vm.readMem8(s))
	}
}

func (vm *VM) orImm(d uint8, s uint16) {
	dr := Register(d)
	vm.WriteReg(dr, vm.ReadReg(dr)|s)
}

func (vm *VM) orReg(d uint8, s uint8) {
	vm.WriteReg(Register(d), vm.ReadReg(Register(s))|vm.ReadReg(Register(d)))
}

func (vm *VM) andAddr(d uint8, s uint16) {
	dr := Register(d)
	switch {
	case IsDoubleReg(dr):
		vm.writeDoubleReg(dr, vm.ReadReg(dr)&vm.readMem16(s))
	default:
		vm.writeSingleReg(dr, uint8(vm.ReadReg(dr))&vm.readMem8(s))
	}
}

func (vm *VM) andImm(d uint8, s uint16) {
	dr := Register(d)
	vm.WriteReg(dr, vm.ReadReg(dr)&s)
}

func (vm *VM) andReg(d uint8, s uint8) {
	vm.WriteReg(Register(d), vm.ReadReg(Register(s))&vm.ReadReg(Register(d)))
}

func (vm *VM) xorAddr(d uint8, s uint16) {
	dr := Register(d)
	switch {
	case IsDoubleReg(dr):
		vm.writeDoubleReg(dr, vm.ReadReg(dr)^vm.readMem16(s))
	default:
		vm.writeSingleReg(dr, uint8(vm.ReadReg(dr))^vm.readMem8(s))
	}
}

func (vm *VM) xorImm(d uint8, s uint16) {
	dr := Register(d)
	vm.WriteReg(dr, vm.ReadReg(dr)^s)
}

func (vm *VM) xorReg(d uint8, s uint8) {
	vm.WriteReg(Register(d), vm.ReadReg(Register(s))^vm.ReadReg(Register(d)))
}

func (vm *VM) rotrRegister(r, x uint8) {
	rr := Register(r)
	switch {
	case IsDoubleReg(rr):
		vm.writeDoubleReg(rr, bits.RotateLeft16(vm.readDoubleReg(rr), int(-x)))
	default:
		vm.writeSingleReg(rr, bits.RotateLeft8(vm.readSingleReg(rr), int(-x)))
	}
}

func (vm *VM) rotlRegister(r, x uint8) {
	rr := Register(r)
	switch {
	case IsDoubleReg(rr):
		vm.writeDoubleReg(rr, bits.RotateLeft16(vm.readDoubleReg(rr), int(x)))
	default:
		vm.writeSingleReg(rr, bits.RotateLeft8(vm.readSingleReg(rr), int(x)))
	}
}

func (vm *VM) jumpEq(r uint8, d uint16) {
	if vm.ReadReg(Register(r)) == uint16(vm.readSingleReg(0)) {
		vm.pc = d
	}
}

func (vm *VM) jumpAbs(d uint16) {
	vm.pc = d
}

func (vm *VM) loadSPAddr(d uint16) {
	vm.sp = vm.readMem16(d)
}

func (vm *VM) loadSPImm(d uint16) {
	vm.sp = d
}

func (vm *VM) loadSPReg(d uint8) {
	vm.sp = vm.ReadReg(Register(d))
}

func (vm *VM) push(r uint8) {
	rr := Register(r)
	switch {
	case IsDoubleReg(rr):
		vm.push16(vm.readDoubleReg(rr))
	default:
		vm.sp--
		vm.writeMem8(vm.sp, vm.readSingleReg(rr))
	}
}

func (vm *VM) pop(r uint8) {
	rr := Register(r)
	switch {
	case IsDoubleReg(rr):
		vm.writeDoubleReg(rr, vm.pop16())
	default:
		vm.writeSingleReg(rr, vm.readMem8(vm.sp))
		vm.sp++
	}
}

func (vm *VM) push16(v uint16) {
	vm.sp -= 2
	vm.writeMem16(vm.sp, v)
}

func (vm *VM) pop16() uint16 {
	v := vm.readMem16(vm.sp)
	vm.sp += 2
	return v
}

func (vm *VM) calla(pc uint16) {
	vm.push16(vm.pc)
	vm.pc = pc
}

func (vm *VM) callr(r uint8) {
	vm.push16(vm.pc)
	vm.pc = vm.ReadReg(Register(r))
}

func (vm *VM) rtn() {
	vm.pc = vm.pop16()
}
