package main

import ()

type gameboyCPU struct {
	gb *gameboy

	AF uint16
	BC uint16
	DE uint16
	HL uint16

	SP uint16
	PC uint16

	cycles          int
	currInstruction string
}

func initCPU(gb *gameboy) *gameboyCPU {
	cpu := new(gameboyCPU)
	cpu.currInstruction = "LD SP, d16"
	cpu.gb = gb

	return cpu
}

func (cpu *gameboyCPU) cycle() {
	if cpu.cycles == 0 {
		fetchedInstruction := cpu.gb.mmu.read(cpu.PC)
		cpu.PC++

		cpu.cycles, cpu.currInstruction = cpu.decodeAndExecute(fetchedInstruction)
	}
	cpu.cycles--
}

func (cpu *gameboyCPU) A() uint8{
	return uint8(cpu.AF >> 8)
}

func (cpu *gameboyCPU) F() uint8{
	return uint8(cpu.AF & 0xFF)
}

func (cpu *gameboyCPU) B() uint8{
	return uint8(cpu.BC >> 8)
}

func (cpu *gameboyCPU) C() uint8{
	return uint8(cpu.BC & 0xFF)
}

func (cpu *gameboyCPU) D() uint8{
	return uint8(cpu.DE >> 8)
}

func (cpu *gameboyCPU) E() uint8{
	return uint8(cpu.DE & 0xFF)
}

func (cpu *gameboyCPU) H() uint8{
	return uint8(cpu.HL >> 8)
}

func (cpu *gameboyCPU) L() uint8{
	return uint8(cpu.HL & 0xFF)
}



func (cpu *gameboyCPU) decodeAndExecute(opcode uint8) (int, string) {
	return 1, "sad"
}
