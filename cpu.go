package main

import (
	"fmt"
	"os"
	"strconv"
)

type gameboyCPU struct {
	A, B, C, D, E, F, H, L uint8 //Registers -> F is the flag register
	SP                     uint16
	PC                     uint16

	cycles int
	currInstruction string

	//include pointers to memory and apu structs
	memory *memory
}

func newCPU() *gameboyCPU {
	cpu := new(gameboyCPU)
	cpu.memory = initMemory(cpu)
	cpu.memory.loadBootrom("roms/bootrom/DMG_ROM.gb")
	return cpu
}

func (cpu *gameboyCPU) cycle(){
	if cpu.cycles == 0{
		fetchedInstruction := cpu.memory.read(cpu.PC)
		cpu.PC++

		cpu.cycles,cpu.currInstruction = cpu.decodeAndExecute(fetchedInstruction)
	}
	cpu.cycles--
}

func main() {
	fmt.Println("PurpleBoy!")
	cpu := newCPU()
	cycleNum,_ := strconv.ParseInt(os.Args[1],10,64)
	fmt.Println("LD SP,d16")

	for i := 0; i < int(cycleNum); i++{
		cpu.cycle()
	}
}
//TODO: WRITE TESTS FOR COMPONENTS
//TODO: GENERATE OPCODES
