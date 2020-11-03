package main

import (
//"fmt"
)

func (cpu *gameboyCPU) setFlag(flag string, value uint8) {
	switch flag {
	case "Z":
		if value == 1 {
			cpu.AF |= 128
		} else if value == 0 {
			cpu.AF &^= 128
		}
	case "N":
		if value == 1 {
			cpu.AF |= 64
		} else if value == 0 {
			cpu.AF &^= 64
		}
	case "H":
		if value == 1 {
			cpu.AF |= 32
		} else if value == 0 {
			cpu.AF &^= 32
		}
	case "C":
		if value == 1 {
			cpu.AF |= 16
		} else if value == 0 {
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


//Addressing modes
func (cpu *gameboyCPU) d8() uint8 {
	immediateData := cpu.gb.mmu.read(cpu.PC)
	cpu.PC++
	return immediateData
}

func (cpu *gameboyCPU) d16() uint16{
	hi := uint16(cpu.gb.mmu.read(cpu.PC+1))
	low := uint16(cpu.gb.mmu.read(cpu.PC))
	cpu.PC += 2
	
	return hi << 8 | low
}

