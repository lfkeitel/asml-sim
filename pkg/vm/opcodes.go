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

func (vm *VM) storeRegInMemory(r uint8, x uint16) {
	switch {
	case isDoubleReg(r):
		vm.writeMem(x, 2, vm.readDoubleReg(r))
	default:
		vm.writeMem(x, 1, uint16(vm.readReg(r)))
	}
}

func (vm *VM) moveRegisters(r, s uint8) {
	vm.writeAnyReg(r, vm.readAnyReg(s))
}

func (vm *VM) addCompliment(r, s, t uint8) {
	sv := vm.readAnyReg2Comp(s)
	tv := vm.readAnyReg2Comp(t)
	vm.writeAnyReg(r, uint16(int16(sv)+int16(tv)))
}

func (vm *VM) addImmCompliment(r, s, x uint8) {
	sv := vm.readAnyReg2Comp(s)
	// converting to int8 then int16 preserves the signed value of uint8
	vm.writeAnyReg(r, uint16(int16(sv)+int16(int8(x))))
}

func (vm *VM) orRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint16(int16(sv)|int16(tv)))
}

func (vm *VM) andRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint16(int16(sv)&int16(tv)))
}

func (vm *VM) xorRegisters(r, s, t uint8) {
	sv := vm.readAnyReg(s)
	tv := vm.readAnyReg(t)
	vm.writeAnyReg(r, uint16(int16(sv)^int16(tv)))
}

func (vm *VM) rotateRegister(r, x uint8) {
	switch {
	case isDoubleReg(r):
		vm.writeDoubleReg(r, bits.RotateLeft16(vm.readDoubleReg(r), int(-x)))
	default:
		vm.writeReg(r, bits.RotateLeft8(vm.readReg(r), int(-x)))
	}
}

func (vm *VM) jumpEq(r uint8, d uint16) {
	if vm.readAnyReg(r) == uint16(vm.readReg(0)) {
		vm.pc = d
	}
}

func (vm *VM) storeRegInMemoryAddr(d, s uint8) {
	addr := vm.readAnyReg(d)
	vm.storeRegInMemory(s, addr)
}

func (vm *VM) loadRegInMemoryAddr(d, s uint8) {
	addr := vm.readAnyReg(s)
	vm.loadFromMem(d, addr)
}

func (vm *VM) jumpAbs(d uint16) {
	vm.pc = d
}

func (vm *VM) loadSP(d uint16) {
	vm.sp = vm.readMem16(d)
}

func (vm *VM) loadSPImm(d uint16) {
	vm.sp = d
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

func (vm *VM) call(pc uint16) {
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
