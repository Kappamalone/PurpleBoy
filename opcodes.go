package main

import (
//"fmt"
//"strings"
)

func opcodeFormat(patternArray [8]uint8, opcode uint8) bool {
	//Takes an input in the form of a string such as
	//"11220011" and return true if the opcode matches
	//the pattern (2 are ignored bits)

	match := true
	for i := 0; i < 8; i++ {
		if patternArray[i] != 2 {
			if patternArray[i] == 1 {
				if (opcode & (1 << (7 - i))) == 0 { //Checks if (7-ith) bit is not set
					match = false
				}
			} else if patternArray[i] == 0 {
				if (opcode & (1 << (7 - i))) > 0 { //Checks if (7-ith) bit is not set
					match = false
				}
			}
		}
	}

	return match
}

func (cpu *gameboyCPU) setFlag(flag string, operand bool) {
	switch flag {
	case "Z":
		if operand {
			cpu.AF |= 128
		} else {
			cpu.AF &^= 128
		}
	case "N":
		if operand {
			cpu.AF |= 64
		} else {
			cpu.AF &^= 64
		}
	case "H":
		if operand {
			cpu.AF |= 32
		} else {
			cpu.AF &^= 32
		}
	case "C":
		if operand {
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
		if (cpu.AF>>7)&1 == 1 {
			flagSet = true
		}
	case "N":
		if (cpu.AF>>6)&1 == 1 {
			flagSet = true
		}
	case "H":
		if (cpu.AF>>5)&1 == 1 {
			flagSet = true
		}
	case "C":
		if (cpu.AF>>4)&1 == 1 {
			flagSet = true
		}
	}
	return flagSet
}

func (cpu *gameboyCPU) carry() uint8 {
	if cpu.getFlag("Z") {
		return uint8(1)
	} else {
		return uint8(0)
	}
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

func (cpu *gameboyCPU) INC(opcode uint8) {
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(1)&0x10 == 0x10))

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) DEC(opcode uint8) {
	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(1)) > 0xF)

	cpu.r8Write[opcode](cpu.r8Read[opcode]() - 1)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) ADD(opcode uint8, operand uint8) {
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F)&0x10 == 0x10))
	cpu.setFlag("C", uint32(cpu.r8Read[opcode]())+uint32(operand) > 0xFFFF)

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)

	cpu.currInstruction = "ADD A, r8"
}

func (cpu *gameboyCPU) ADC(opcode uint8, operand uint8) {
	cpu.setFlag("H", (((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F)+cpu.carry())&0x10 == 0x10))
	cpu.setFlag("C", (uint32(cpu.r8Read[opcode]())+uint32(operand)+uint32(cpu.carry())) > 0xFFFF)

	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand + cpu.carry())

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)

	cpu.currInstruction = "ADC A, r8"
}

func (cpu *gameboyCPU) SUB(opcode uint8, operand uint8) {
	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(operand&0x0F)) > 0xF)
	cpu.setFlag("C", (uint32(cpu.r8Read[opcode]())-uint32(operand) > 0xFFFF))

	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)

	cpu.currInstruction = "SUB A, r8"
}

func (cpu *gameboyCPU) SUBC(opcode uint8, operand uint8) {
	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(operand&0x0F)-cpu.carry()) > 0xF)
	cpu.setFlag("C", ((uint32(cpu.r8Read[opcode]()) - uint32(operand) - uint32(cpu.carry())) > 0xFFFF))

	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand - cpu.carry())

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)

	cpu.currInstruction = "SBC A, r8"
}

func (cpu *gameboyCPU) AND(opcode uint8, operand uint8) {

	cpu.r8Write[opcode](cpu.r8Read[opcode]() & cpu.r8Read[operand]())

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("H", true)

	cpu.currInstruction = "AND A, r8"
}

func (cpu *gameboyCPU) XOR(opcode uint8, operand uint8) {

	cpu.r8Write[opcode](cpu.r8Read[opcode]() ^ cpu.r8Read[operand]())

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)

	cpu.currInstruction = "XOR A, r8"
}

func (cpu *gameboyCPU) OR(opcode uint8, operand uint8) {

	cpu.r8Write[opcode](cpu.r8Read[opcode]() | cpu.r8Read[operand]())
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)

	cpu.currInstruction = "OR A, r8"
}

func (cpu *gameboyCPU) CP(opcode uint8, operand uint8) {
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)

	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(operand&0x0F)) > 0xF)
	cpu.setFlag("C", (uint32(cpu.r8Read[opcode]())-uint32(operand) > 0xFFFF))

	cpu.currInstruction = "CP A, r8"
}

//ALU
func (cpu *gameboyCPU) CCF(){
	var carryFlag uint8
	if cpu.getFlag("C"){
		carryFlag = 1
	} else {
		carryFlag = 0
	}

	cpu.setFlag("C",carryFlag != 1) //XOR -> A != B
}

func (cpu *gameboyCPU) SCF(){
	cpu.setFlag("C",true)
}


//EXTENDED--------------------------------
func (cpu *gameboyCPU) BIT(opcode uint8, place uint8){
	//Set zflag is bit not set 
	cpu.setFlag("Z",cpu.r8Read[opcode]() & (1 << place) == 1)
}

func (cpu *gameboyCPU) RES(opcode uint8, place uint8){
	//Reset bit
	//Create a mask that is identical to the opcode except the bit we are resetting
	//Eg: resetting 3rd bit: 0100 -> 1011 is now the mask to and with the opcode
	mask := uint8(^(1 << place))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() & mask)
}

func (cpu *gameboyCPU) SET(opcode uint8,place uint8){
	//Set bit
	cpu.r8Write[opcode](cpu.r8Read[opcode]() | (1 << place))
}
