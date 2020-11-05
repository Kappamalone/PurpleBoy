package main

import (
	//"fmt"
	// "os"
	// "strconv"
)

type gameboy struct {
	cpu   *gameboyCPU
	mmu   *memory
	debug *debugger
}

func initGameboy(isDebugging bool) *gameboy {
	gb := new(gameboy)
	gb.cpu = initCPU(gb)
	gb.mmu = initMemory(gb)
	if isDebugging {
		gb.debug = initDebugger(gb)
	}

	return gb
}

func main() {
	gb := initGameboy(true)
	gb.mmu.loadBlaarg("roms/testroms/cpu_instrs/cpu_instrs.gb")

	for {
		handleInput()
	}
}
