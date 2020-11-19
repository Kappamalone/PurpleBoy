package main

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	//The gameboy is clocked at a speed of 4.194304 MHz,
	//Therefore each frame you'd execute 1/60 of that total amount
	clockspeed     = 4194304
	cyclesPerFrame = clockspeed / 60
)

type gameboy struct {
	cpu   *gameboyCPU
	ppu   *PPU
	mmu   *memory
	debug *debugger
}

func initGameboy(skipBootrom bool, isDebugging bool) *gameboy {
	gb := new(gameboy)
	gb.ppu = initPPU(gb)
	gb.mmu = initMemory(gb)
	gb.cpu = initCPU(gb, skipBootrom)
	if isDebugging {
		gb.debug = initDebugger(gb, isLogging)
	}

	return gb
}

var (
	cfile       string = "10-bit ops"
	skipBootrom bool   = true
	isDebugging bool   = true
	isLogging   bool   = false
)

func main() {
	gb := initGameboy(skipBootrom, isDebugging)
	gb.mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")
	gb.mmu.loadBlaarg(fmt.Sprintf("roms/testroms/cpu_instrs/%s.gb", cfile))

	defer ui.Close()
	defer sdl.Quit()
	defer gb.ppu.window.Destroy()
	defer gb.ppu.renderer.Destroy()

	ticker := time.NewTicker(16 * time.Millisecond)
	for range ticker.C {
		//Use the emitted tick from the ticker to run 1/60th of the required frames every 1/60th of a second
		if isDebugging {
			gb.debug.updateDebugInformation()
			ui.Render(gb.debug.cpuState, gb.debug.consoleOut)
		}

		for i := 0; i < cyclesPerFrame; i++ {
			gb.cpu.cycle()

			if gb.mmu.ram[0xFF02] == 0x81 {
				gb.debug.printConsole(fmt.Sprintf("%c", gb.mmu.ram[0xFF01]),"green")
				gb.mmu.ram[0xFF02] = 0x00
			}
		}

		if handleInput() {
			ticker.Stop() //Stop ticker to exit program
			break
		}
	}
}
