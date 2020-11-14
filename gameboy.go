package main

import (
	"fmt"
	// "strconv"
)

type gameboy struct {
	cpu   *gameboyCPU
	mmu   *memory
	debug *debugger
}

func initGameboy(isDebugging bool) *gameboy {
	gb := new(gameboy)
	gb.mmu = initMemory(gb)
	gb.cpu = initCPU(gb)
	if isDebugging {
		gb.debug = initDebugger(gb, true)
	}

	return gb
}

var cfile string = "02-interrupts"

func main() {
	gb := initGameboy(true)
	//gb.mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")
	gb.mmu.loadBlaarg(fmt.Sprintf("roms/testroms/cpu_instrs/%s.gb", cfile))

	for {
		gb.debug.logTrace()
		gb.cpu.cycle()
		if gb.mmu.ram[0xFF02] == 0x81 {
			fmt.Printf("%c", gb.mmu.ram[0xFF01])
			gb.mmu.ram[0xFF02] = 0x00
		}
	}
}

//Passed tests
//01
//03
//04
//05
//06
//07
//08
//09
//10
//11
