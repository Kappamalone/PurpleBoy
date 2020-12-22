package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var (
	//The gameboy is clocked at a speed of 4.194304 MHz,
	//Therefore each frame you'd execute 1/60 of that total amount
	clockspeed     = 4194304
	cyclesPerFrame = int(float32(clockspeed) / 59.73)

	//Global control vars
	skipBootrom bool = true
	isDebugging bool = true
	isLogging   bool = false
	useTestRom  bool = false

	title   string = "FF Adventure"
	testrom string = "roms/testroms/cpu_instrs/instr_timing.gb"
	gamerom string = fmt.Sprintf("roms/gameroms/%s.gb", title)
	inputROM string
)

type gameboy struct {
	cpu    *gameboyCPU
	ppu    *PPU
	mmu    *memory
	joypad *joypad
	debug  *debugger
}

func initGameboy() *gameboy {
	gb := new(gameboy)
	gb.ppu = initPPU(gb)
	gb.mmu = initMemory(gb)
	gb.cpu = initCPU(gb)
	if isDebugging {
		gb.debug = initDebugger(gb)
	}
	gb.joypad = initJoypad(gb)

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
	//Command line flags
	flag.BoolVar(&skipBootrom, "skip", true, "Used for disabling skipping the bootrom")
	flag.StringVar(&inputROM,"path","","Used to input path of ROM")
	flag.BoolVar(&useTestRom,  "t", false, "Used for picking gamerom or testrom")
	flag.Parse()

	if inputROM != "" {
		//This is here because I didn't setup a development branch :P
		gamerom = inputROM + ".gb"
	}

	gb := initGameboy()
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
