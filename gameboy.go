package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	//The gameboy is clocked at a speed of 4.194304 MHz,
	//Therefore each frame you'd execute 1/60 of that total amount
	clockspeed     = 4194304
	cyclesPerFrame = clockspeed / 60
)

var (
	skipBootrom bool = true
	isDebugging bool = true
	isLogging   bool = false

	testrom string = "roms/testroms/MBC/bits_ramg.gb"
	//fullrom string = "roms/gameroms/LoZ Link's Awakening.gb"
	fullrom string = "roms/gameroms/Super Mario Land.gb"
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
	gb.mmu = initMemory(gb, skipBootrom)
	gb.cpu = initCPU(gb, skipBootrom)
	if isDebugging {
		gb.debug = initDebugger(gb, isLogging)
	}

	return gb
}

func (gb *gameboy) handleDebug() {
	if isDebugging {
		gb.debug.updateDebugInformation()
		ui.Render(gb.debug.cpuState, gb.debug.consoleOut, gb.debug.firedInterrupts)

		gb.ppu.displayTileset()
	}
}

func (gb *gameboy) handleLogging() {
	if isLogging {
		//gb.debug.logTrace()
		//gb.debug.logValue(gb.mmu.cart.rombankNum)
	}
}

func main() {
	gb := initGameboy(skipBootrom, isDebugging)
	ticker := time.NewTicker(16 * time.Millisecond)

	if isLogging {
		gb.debug.printConsole("Logging Enabled!\n", "green")
	}

	//One frame
	for range ticker.C {
		gb.handleDebug()

		for i := 0; i < cyclesPerFrame; i++ {
			//System is clocked at 4.2MHZ
			gb.cpu.tick(i)
			gb.ppu.tick()
			gb.cpu.timers.tick()

		}
		if handleInput() {
			ticker.Stop() //Stop ticker to exit program
			break
		}
	}

	if isDebugging {
		defer ui.Close()
		defer gb.ppu.tileWindow.Destroy()
		defer gb.ppu.tileRenderer.Destroy()
	}
	defer sdl.Quit()
	defer gb.ppu.window.Destroy()
	defer gb.ppu.renderer.Destroy()
}
