package main

import (
//"fmt"
)

func (cpu *gameboyCPU) setFlag(flag string, value bool) {
	switch flag {
	case "Z":
		if value {
			cpu.AF |= 128
		} else {
			cpu.AF &^= 128
		}
	case "N":
		if value {
			cpu.AF |= 64
		} else {
			cpu.AF &^= 64
		}
	case "H":
		if value {
			cpu.AF |= 32
		} else {
			cpu.AF &^= 32
		}
	case "C":
		if value {
			cpu.AF |= 16
		} else {
			cpu.AF &^= 16
		}
	}
}

func (cpu *gameboyCPU) getFlag(flag string) uint16 {
	var flagbit uint16
	switch flag {
	case "Z":
		flagbit = (cpu.AF >> 7) & 1
	case "N":
		flagbit = (cpu.AF >> 6) & 1
	case "H":
		flagbit = (cpu.AF >> 5) & 1
	case "C":
		flagbit = (cpu.AF >> 4) & 1
	}
	return flagbit
}

//Get 8 bit registers
func (cpu *gameboyCPU) A() uint8 {
	return uint8(cpu.AF >> 8)
}

func (cpu *gameboyCPU) B() uint8 {
	return uint8(cpu.BC >> 8)
}

func (cpu *gameboyCPU) C() uint8 {
	return uint8(cpu.BC & 0xFF)
}

func (cpu *gameboyCPU) D() uint8 {
	return uint8(cpu.DE >> 8)
}

func (cpu *gameboyCPU) E() uint8 {
	return uint8(cpu.DE & 0xFF)
}

func (cpu *gameboyCPU) H() uint8 {
	return uint8(cpu.HL >> 8)
}

func (cpu *gameboyCPU) L() uint8 {
	return uint8(cpu.HL & 0xFF)
}

//Set 8 bit registers
func (cpu *gameboyCPU) setA(data uint8) {
	cpu.AF = uint16(data)<<8 | cpu.AF&0xFF
}

func (cpu *gameboyCPU) setB(data uint8) {
	cpu.BC = uint16(data)<<8 | cpu.BC&0xFF
}

func (cpu *gameboyCPU) setC(data uint8) {
	cpu.BC = cpu.BC<<8 | uint16(data)
}

func (cpu *gameboyCPU) setD(data uint8) {
	cpu.DE = uint16(data)<<8 | cpu.DE&0xFF
}

func (cpu *gameboyCPU) setE(data uint8) {
	cpu.DE = cpu.DE<<8 | uint16(data)
}

func (cpu *gameboyCPU) setH(data uint8) {
	cpu.HL = uint16(data)<<8 | cpu.HL&0xFF
}

func (cpu *gameboyCPU) setL(data uint8) {
	cpu.HL = cpu.HL<<8 | uint16(data)
}

//Addressing modes
func (cpu *gameboyCPU) d8() uint8 {
	immediateData := cpu.gb.mmu.readbyte(cpu.PC)
	cpu.PC++
	return immediateData
}

func (cpu *gameboyCPU) d16() uint16 {
	hi := uint16(cpu.gb.mmu.readbyte(cpu.PC + 1))
	low := uint16(cpu.gb.mmu.readbyte(cpu.PC))
	cpu.PC += 2

	return hi<<8 | low
}

//OPCODES
//Wonderful explanation for half carry flags at https://robdor.com/2016/08/10/gameboy-emulator-half-carry-flag/

func (cpu *gameboyCPU) INCR8(opcode uint8) {
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(1)&0x10 == 0x10))

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) DECR8(opcode uint8) {
	//REMEMBER H flag

	cpu.r8Write[opcode](cpu.r8Read[opcode]() - 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) ADDR8(opcode uint8, value uint8) {
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(value&0x0F)&0x10 == 0x10))
	cpu.setFlag("C", uint32(cpu.r8Read[opcode]())+uint32(value) > 0xFFFF)

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + value)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) SUBR8(opcode uint8, value uint8) {
	//REMEMBER H flag
	
	cpu.setFlag("C", (uint32(cpu.r8Read[opcode]())-uint32(value) > 0xFFFF))


	cpu.r8Write[opcode](cpu.r8Read[opcode]() - value)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)
}
