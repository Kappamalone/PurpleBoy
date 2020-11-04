package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
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
	defer ui.Close() //close termui
	gb := initGameboy(true)
	gb.mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")

	a := 0x0000
	b := 0x0
	fmt.Println(uint32(a)-uint32(b) > 0xFFFF)

	for {
		handleInput()
	}
}
