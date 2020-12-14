package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"time"
	"flag"
)

const (
	//The gameboy is clocked at a speed of 4.194304 MHz,
	//Therefore each frame you'd execute 1/60 of that total amount
	clockspeed     = 4194304
	cyclesPerFrame = clockspeed / 60
)

var (
	//rom_8mb fails for some unholy reason

	skipBootrom bool = true
	isDebugging bool = false
	isLogging   bool = false

	title   string = "Tetris"
	testrom string = "roms/testroms/MBC/rom_8mb.gb"
	gamerom string = fmt.Sprintf("roms/gameroms/%s.gb",title)
	useTestRom bool = false
)

type gameboy struct {
	cpu    *gameboyCPU
	ppu    *PPU
	mmu    *memory
	joypad *joypad
	debug  *debugger
}

func initGameboy(skipBootrom bool, isDebugging bool) *gameboy {
	gb := new(gameboy)
	gb.ppu = initPPU(gb)
	gb.mmu = initMemory(gb, skipBootrom)
	gb.cpu = initCPU(gb, skipBootrom)
	gb.joypad = initJoypad(gb)
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
	flag.BoolVar(&useTestRom,"t",false,"Used for picking gamerom or testrom")
	flag.Parse()

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
		if gb.joypad.handleInput() {
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
