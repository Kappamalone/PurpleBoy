package main

import (
	"fmt"
)

func (cpu *gameboyCPU) getCRegisters(register string) uint16 {
	//double check the setting of registers with bgb
	var registerData uint16
	switch register {
	case "AF":
		registerData = uint16(cpu.A)<<8 | uint16(cpu.F&0x0F) //Lower 4 bits of F are always 0
	case "BC":
		registerData = uint16(cpu.B)<<8 | uint16(cpu.C)
	case "DE":
		registerData = uint16(cpu.D)<<8 | uint16(cpu.E)
	case "HL":
		registerData = uint16(cpu.H)<<8 | uint16(cpu.L)
	}

	return registerData
}

//Returns pointers to register pairs to be used with LD16,ADD16
func (cpu *gameboyCPU) getCRegisterPointers(register string) [2]*uint8 {
	//double check the setting of registers with bgb
	var high *uint8
	var low *uint8
	switch register {
	case "AF":
		high = &cpu.A
		low = &cpu.F
	case "BC":
		high = &cpu.B
		low = &cpu.C
	case "DE":
		high = &cpu.D
		low = &cpu.E
	case "HL":
		high = &cpu.H
		low = &cpu.L
	}

	return [2]*uint8{high, low}
}

func (cpu *gameboyCPU) setFlag(flag string, value uint8) {
	switch flag {
	case "Z":
		if value == 1 {
			cpu.F |= 128
		} else if value == 0 {
			cpu.F &^= 128
		}
	case "N":
		if value == 1 {
			cpu.F |= 64
		} else if value == 0 {
			cpu.F &^= 64
		}
	case "H":
		if value == 1 {
			cpu.F |= 32
		} else if value == 0 {
			cpu.F &^= 32
		}
	case "C":
		if value == 1 {
			cpu.F |= 16
		} else if value == 0 {
			cpu.F &^= 16
		}
	}
}

func (cpu *gameboyCPU) getFlag(flag string) uint8 {
	var flagbit uint8
	switch flag {
	case "Z":
		flagbit = (cpu.F >> 7) & 1
	case "N":
		flagbit = (cpu.F >> 6) & 1
	case "H":
		flagbit = (cpu.F >> 5) & 1
	case "C":
		flagbit = (cpu.F >> 4) & 1
	}
	return flagbit
}

//Sets zflag if result is equal to zero
func (cpu *gameboyCPU) zflag(result uint8) {
	if result == 0 {
		cpu.setFlag("Z", 1)
	}
}

//No Nflag as no logic required

//Sets hflag if half carry occurs
//I swear I have no idea how this is supposed to be configured on stuff like subtractions
func (cpu *gameboyCPU) hflag(original uint8, val uint8, variant string) {
	//The gameboy manual shows an example where 0xFF + 1 is a half carry, so i hardcoded it in
	switch variant {
	case "ADD", "INC":
		if ((original >> 4) & (original & 1)) == 1 {
			cpu.setFlag("H", 1)
		} else if (original>>4) == 0 && ((original+val)>>4) == 1 {
			cpu.setFlag("H", 1)
		}
	}
}

//Sets cflag if carry occurs with 8 bit num
//ADD: if greater than 0xFF
//SUB,CP: if less than 0x00
//May have to add additional variants for bit rotations
func (cpu *gameboyCPU) cflag(original interface{}, val uint8, variant string) {
	switch original := original.(type) {
	case uint8:
		switch variant {
		case "ADD":
			if ((uint16(original) + uint16(val)) >> 8) == 1 {
				cpu.setFlag("C", 1)
			}
		case "SUB", "CP":
			if (uint16(original) - uint16(val)) > 0xFF {
				cpu.setFlag("C", 1)
			}
		}

	case uint16:
		switch variant {
		case "ADD":
			if ((uint32(original) + uint32(val)) >> 16) == 1 {
				cpu.setFlag("C", 1)
			}
		case "SUB", "CP":
			if (uint32(original) - uint32(val)) > 0xFFFF {
				cpu.setFlag("C", 1)
			}
		}
	default:
		fmt.Printf("This shouldn't happen! %T\n", original)
	}
}

//Helper function to get address from pc and pc+1
func (cpu *gameboyCPU) d16() uint16 {
	//The gameboy is little endian, meaning the least significant byte 
	//is stored first in memory. So reverse the address and the access memory
	hi := cpu.memory.read(cpu.PC+1)
	lo := cpu.memory.read(cpu.PC)
	d16 := uint16(hi) << 8 | uint16(lo)

	return d16
}


//LD location with value
func LD(location *uint8, value uint8) {
	*location = value
}

//LD16 bit variant for combined registers
func LD16(locations [2]*uint8, value uint16) {
	//[high,low]*uint8
	*(locations)[0] = uint8(value >> 8)
	*(locations)[1] = uint8(value & 0xFF)
}

//LD16P variant for PC/SP
func LD16P(location *uint16, value uint16) {
	*location = value
}

//ADD register with value
func ADD(location *uint8, value uint8) {
	*location += value
}

//ADD16 bit variant
func ADD16(locations [2]*uint8, value uint16) {
	cRegister := uint16(*(locations)[0])<<8 | uint16(*(locations)[1])&0xFF
	value += cRegister
	*(locations)[0] = uint8(value >> 8)
	*(locations)[1] = uint8(value & 0xFF)
}
