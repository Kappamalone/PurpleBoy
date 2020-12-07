package main

import (
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
	gb.mmu = initMemory(gb,skipBootrom)
	gb.cpu = initCPU(gb, skipBootrom)
	if isDebugging {
		gb.debug = initDebugger(gb, isLogging)
	}

	return gb
}

var (
	cfile       string = "02-interrupts"
	skipBootrom bool   = true
	isDebugging bool   = true
	isLogging   bool   = false

	fullrom string = "roms/gameroms/Dr mario.gb"
)

func main() {
	gb := initGameboy(skipBootrom, isDebugging)
	ticker := time.NewTicker(16 * time.Millisecond)

	//One frame
	for range ticker.C {

		if isDebugging {
			gb.debug.updateDebugInformation()
			ui.Render(gb.debug.cpuState, gb.debug.consoleOut)

			gb.ppu.displayTileset()
			gb.ppu.displayCurrTileMap()
		}

		for i := 0; i < cyclesPerFrame; i++ {
			//System is clocked at 4.2MHZ
			gb.cpu.tick()
			gb.ppu.tick() 

			// if isDebugging {
			// 	if gb.mmu.ram[0xFF02] == 0x81 {
			// 		gb.debug.printConsole(fmt.Sprintf("%c", gb.mmu.ram[0xFF01]), "cyan")
			// 		gb.mmu.ram[0xFF02] = 0x00
			// 	}
			// }
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
		defer gb.ppu.fullWindow.Destroy()
		defer gb.ppu.fullRenderer.Destroy()
	}
	defer sdl.Quit()
	defer gb.ppu.window.Destroy()
	defer gb.ppu.renderer.Destroy()
}
