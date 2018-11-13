package vm

import (
	"math/bits"
)

// Opcode definitions

func (vm *VM) loadFromMem(r uint8, x uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, uint16(vm.readMem(x, 2)))
	default:
		vm.writeReg(r, uint8(vm.readMem(x, 1)))
	}
}

func (vm *VM) loadIntoReg(r uint8, x uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, x)
	default:
		vm.writeReg(r, uint8(x))
	}
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	addr := vm.readAnyReg(s)
	vm.loadFromMem(d, addr)
}

func (vm *VM) storeRegToMemory(r uint8, x uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeMem(x, 2, vm.readDoubleReg(r))
	default:
		vm.writeMem(x, 1, uint16(vm.readReg(r)))
	}
}

func (vm *VM) storeRegToRegAddr(s, d uint8) {
	addr := vm.readAnyReg(d)
	vm.storeRegToMemory(s, addr)
}

func (vm *VM) xferRegisters(r, s uint8) {
	vm.writeAnyReg(r, vm.readAnyReg(s))
}

func (vm *VM) addAddr(d uint8, s uint16) {
	switch {
	case isDoubleReg(d):
		vm.writeDoubleReg(d, vm.readAnyReg(d)+vm.readMem16(s))
	default:
		vm.writeReg(d, uint8(vm.readAnyReg(d))+vm.readMem8(s))
	}
}

func (vm *VM) addImm(d uint8, s uint16) {
	vm.writeAnyReg(d, vm.readAnyReg(d)+s)
}

func (vm *VM) addReg(d uint8, s uint8) {
	vm.writeAnyReg(d, vm.readAnyReg(s)+vm.readAnyReg(d))
}

func (vm *VM) orAddr(d uint8, s uint16) {
	switch {
	case isDoubleReg(d):
		vm.writeDoubleReg(d, vm.readAnyReg(d)|vm.readMem16(s))
	default:
		vm.writeReg(d, uint8(vm.readAnyReg(d))|vm.readMem8(s))
	}
}

func (vm *VM) orImm(d uint8, s uint16) {
	vm.writeAnyReg(d, vm.readAnyReg(d)|s)
}

func (vm *VM) orReg(d uint8, s uint8) {
	vm.writeAnyReg(d, vm.readAnyReg(s)|vm.readAnyReg(d))
}

func (vm *VM) andAddr(d uint8, s uint16) {
	switch {
	case isDoubleReg(d):
		vm.writeDoubleReg(d, vm.readAnyReg(d)&vm.readMem16(s))
	default:
		vm.writeReg(d, uint8(vm.readAnyReg(d))&vm.readMem8(s))
	}
}

func (vm *VM) andImm(d uint8, s uint16) {
	vm.writeAnyReg(d, vm.readAnyReg(d)&s)
}

func (vm *VM) andReg(d uint8, s uint8) {
	vm.writeAnyReg(d, vm.readAnyReg(s)&vm.readAnyReg(d))
}

func (vm *VM) xorAddr(d uint8, s uint16) {
	switch {
	case isDoubleReg(d):
		vm.writeDoubleReg(d, vm.readAnyReg(d)^vm.readMem16(s))
	default:
		vm.writeReg(d, uint8(vm.readAnyReg(d))^vm.readMem8(s))
	}
}

func (vm *VM) xorImm(d uint8, s uint16) {
	vm.writeAnyReg(d, vm.readAnyReg(d)^s)
}

func (vm *VM) xorReg(d uint8, s uint8) {
	vm.writeAnyReg(d, vm.readAnyReg(s)^vm.readAnyReg(d))
}

func (vm *VM) rotrRegister(r, x uint8) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, bits.RotateLeft16(vm.readDoubleReg(r), int(-x)))
	default:
		vm.writeReg(r, bits.RotateLeft8(vm.readReg(r), int(-x)))
	}
}

func (vm *VM) rotlRegister(r, x uint8) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, bits.RotateLeft16(vm.readDoubleReg(r), int(x)))
	default:
		vm.writeReg(r, bits.RotateLeft8(vm.readReg(r), int(x)))
	}
}

func (vm *VM) jumpEq(r uint8, d uint16) {
	if vm.readAnyReg(r) == uint16(vm.readReg(0)) {
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
	vm.sp = vm.readAnyReg(d)
}

func (vm *VM) push(r uint8) {
	switch {
	case isDoubleReg(r):
		vm.push16(vm.readDoubleReg(r))
	default:
		vm.sp--
		vm.writeMem8(vm.sp, vm.readReg(r))
	}
}

func (vm *VM) pop(r uint8) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, vm.pop16())
	default:
		vm.writeReg(r, vm.readMem8(vm.sp))
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
	vm.pc = vm.readAnyReg(r)
}

func (vm *VM) rtn() {
	vm.pc = vm.pop16()
}
