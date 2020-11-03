package main

import (
	//"fmt"
	ui "github.com/gizak/termui/v3"
	"os"
	"strconv"
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
	defer ui.Close() //close termui
	gb := initGameboy(true)
	gb.mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")

	cycles, _ := strconv.Atoi(os.Args[1])
	for i := 0; i < cycles; i++ {
		gb.cpu.cycle()
	}

	for {
		handleInput()
	}
}
