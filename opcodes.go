package main

import (
	//"fmt"
	//"strings"
)

func opcodeFormat(patternArray [8]uint8,opcode uint8) bool {
	//Takes an input in the form of a string such as 
	//"11220011" and return true if the opcode matches 
	//the pattern (2 are ignored bits)

	match := true
	for i := 0; i < 8; i++ {
		if patternArray[i] != 2{
			if patternArray[i] == 1{
				if (opcode & (1<<(7-i))) == 0{ //Checks if (7-ith) bit is not set
					match = false
				}
			} else if patternArray[i] == 0{
				if (opcode & (1<<(7-i))) > 0{ //Checks if (7-ith) bit is not set
					match = false
				}
			}
		}
	}

	return match
}

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

func (cpu *gameboyCPU) getFlag(flag string) bool {
	flagSet := false
	switch flag {
	case "Z":
		if (cpu.AF >> 7) & 1 == 1{
			flagSet = true
		}
	case "N":
		if (cpu.AF >> 6) & 1 == 1{
			flagSet = true
		}
	case "H":
		if (cpu.AF >> 5) & 1 == 1{
			flagSet = true
		}
	case "C":
		if (cpu.AF >> 4) & 1 == 1{
			flagSet = true
		}
	}
	return flagSet
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
	opcode &= 0x7
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(1)&0x10 == 0x10))

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) DECR8(opcode uint8) {
	opcode &= 0x07
	//REMEMBER H flag

	cpu.r8Write[opcode](cpu.r8Read[opcode]() - 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) ADDR8(opcode uint8, value uint8) {
	opcode &= 0x07
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(value&0x0F)&0x10 == 0x10))
	cpu.setFlag("C", uint32(cpu.r8Read[opcode]())+uint32(value) > 0xFFFF)

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + value)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) SUBR8(opcode uint8, value uint8) {
	opcode &= 0x07
	//REMEMBER H flag
	
	cpu.setFlag("C", (uint32(cpu.r8Read[opcode]())-uint32(value) > 0xFFFF))


	cpu.r8Write[opcode](cpu.r8Read[opcode]() - value)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)
}
