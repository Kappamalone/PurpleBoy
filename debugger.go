package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	logo = []string{
		" ____                   _      _                 _",
		"|  _ \\ _   _ _ __ _ __ | | ___| |__   ___  _   _| |",
		"| |_) | | | | '__| '_ \\| |/ _ \\ '_ \\ / _ \\| | | | |",
		"|  __/| |_| | |  | |_) | |  __/ |_) | (_) | |_| |_|",
		"|_|    \\__,_|_|  | .__/|_|\\___|_.__/ \\___/ \\__, (_)",
		"                 |_|                       |___/",
	}
)

type debugger struct {
	gb *gameboy

	cpuState   *widgets.Paragraph //CPU Internal registers
	consoleOut *widgets.Paragraph //Console for debug info

	console []string //Data to be rendered by console
}

func createWidget(title string, colour ui.Color, dimensions [4]int) *widgets.Paragraph {
	widget := widgets.NewParagraph()
	widget.Title = title
	widget.BorderStyle.Fg = colour
	widget.SetRect(dimensions[0], dimensions[1], dimensions[2], dimensions[3])
	return widget
}

func (debug *debugger) displayLogo() {
	//Print some cool logo stuff
	for _, line := range logo {
		debug.printConsole(line+"\n", "magenta")
	}
	debug.printConsole("\n", "cyan")
	debug.printConsole("Written by Uzman Zawahir", "cyan")
	debug.printConsole("\n", "cyan")
	debug.printConsole("\n", "cyan")
}

func initDebugger(gb *gameboy, isLogging bool) *debugger {
	debug := new(debugger)
	debug.gb = gb

	//Initialise termui
	err := ui.Init()
	checkErr(err, "Failed to intialise termui")
	debug.cpuState = createWidget("[CPU STATE]", ui.ColorCyan, [4]int{0, 0, 26, 33})
	debug.consoleOut = createWidget("[CONSOLE]", ui.ColorCyan, [4]int{0, 33, 163, 52})
	debug.displayLogo()

	//Create a file for logging
	if isLogging {
		initLogging()
	}

	return debug
}

func initLogging() {
	//Setup logging
	file, err := os.OpenFile("logfiles/log.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(file)
}

func (debug *debugger) logTrace() {
	log.Printf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: 00:%04X (%02X %02X %02X %02X)", debug.gb.cpu.getAcc(), debug.gb.cpu.AF&0x00FF, debug.gb.cpu.r8Read[0](), debug.gb.cpu.r8Read[1](), debug.gb.cpu.r8Read[2](), debug.gb.cpu.r8Read[3](), debug.gb.cpu.r8Read[4](), debug.gb.cpu.r8Read[5](), debug.gb.cpu.SP, debug.gb.cpu.PC, debug.gb.mmu.ram[debug.gb.cpu.PC], debug.gb.mmu.ram[debug.gb.cpu.PC+1], debug.gb.mmu.ram[debug.gb.cpu.PC+2], debug.gb.mmu.ram[debug.gb.cpu.PC+3])
}

func (debug *debugger) logVRAM() {
	for i := 0; i < len(debug.gb.ppu.VRAM); i += 16 {
		log.Printf("%02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X",debug.gb.ppu.VRAM[i],debug.gb.ppu.VRAM[i+1],debug.gb.ppu.VRAM[i+2],debug.gb.ppu.VRAM[i+3],
																	 								 debug.gb.ppu.VRAM[i+4],debug.gb.ppu.VRAM[i+5],debug.gb.ppu.VRAM[i+6],debug.gb.ppu.VRAM[i+7],
																									 debug.gb.ppu.VRAM[i+8],debug.gb.ppu.VRAM[i+9],debug.gb.ppu.VRAM[i+10],debug.gb.ppu.VRAM[i+11],
																									 debug.gb.ppu.VRAM[i+12],debug.gb.ppu.VRAM[i+13],debug.gb.ppu.VRAM[i+14],debug.gb.ppu.VRAM[i+15])
		}
	}

//Write debug windows down here

func (debug *debugger) printConsole(data string, colour string) {
	//Works through a primitive line by line basis
	if len(debug.console) > 18 {
		debug.console = debug.console[1:]
	}
	debug.console = append(debug.console, fmt.Sprintf("[%s](fg:%s)", data, colour))
}

//UPDATE DEBUG INFO
func (debug *debugger) updateDebugInformation() {
	//Update text for all components inside TUI
	//Basically line by line updates

	debugCPU := make([]string, 0)
	debugCPU = append(debugCPU, fmt.Sprintf("\n [PC](fg:cyan) = [$%04X](fg:yellow)", debug.gb.cpu.PC))
	debugCPU = append(debugCPU, fmt.Sprintf(" [SP](fg:cyan) = [$%04X](fg:yellow)\n", debug.gb.cpu.SP))
	debugCPU = append(debugCPU, "[------------------------\n](fg:cyan)")
	debugCPU = append(debugCPU, fmt.Sprintf(" [A](fg:cyan) = [$%02X](fg:yellow)     [F](fg:cyan) = [$%02X](fg:yellow)\n", debug.gb.cpu.AF>>8, debug.gb.cpu.AF&0xFF))
	debugCPU = append(debugCPU, fmt.Sprintf(" [B](fg:cyan) = [$%02X](fg:yellow)     [C](fg:cyan) = [$%02X](fg:yellow)\n", debug.gb.cpu.BC>>8, debug.gb.cpu.BC&0xFF))
	debugCPU = append(debugCPU, fmt.Sprintf(" [D](fg:cyan) = [$%02X](fg:yellow)     [E](fg:cyan) = [$%02X](fg:yellow)\n", debug.gb.cpu.DE>>8, debug.gb.cpu.DE&0xFF))
	debugCPU = append(debugCPU, fmt.Sprintf(" [H](fg:cyan) = [$%02X](fg:yellow)     [L](fg:cyan) = [$%02X](fg:yellow)\n", debug.gb.cpu.HL>>8, debug.gb.cpu.HL&0xFF))
	debugCPU = append(debugCPU, "[\n------------------------\n](fg:cyan)")
	debugCPU = append(debugCPU, fmt.Sprintf(" [Z](fg:cyan) = [%d](fg:yellow)       [N](fg:cyan) = [%d](fg:yellow)", boolToInt(debug.gb.cpu.getZ()), boolToInt(debug.gb.cpu.getN())))
	debugCPU = append(debugCPU, fmt.Sprintf(" [H](fg:cyan) = [%d](fg:yellow)       [C](fg:cyan) = [%d](fg:yellow)", boolToInt(debug.gb.cpu.getH()), boolToInt(debug.gb.cpu.getC())))
	debugCPU = append(debugCPU, "[\n\n------------------------\n](fg:cyan)")
	debugCPU = append(debugCPU, fmt.Sprintf(" [IME](fg:cyan) = [%d](fg:yellow)     [HALT](fg:cyan) = [%d](fg:yellow)", boolToInt(debug.gb.cpu.IME), boolToInt(debug.gb.cpu.HALT)))

	debug.cpuState.Text = strings.Join(debugCPU, "\n")
	debug.consoleOut.Text = strings.Join(debug.console, "")
}

/*
for i := 0; i < 256; i += 16{
	fmt.Printf("%02X %02X %02X %02X %002X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n",gb.mmu.ram[i],gb.mmu.ram[i+1],gb.mmu.ram[i+2],gb.mmu.ram[i+3],gb.mmu.ram[i+4],gb.mmu.ram[i+5],gb.mmu.ram[i+6],gb.mmu.ram[i+7],gb.mmu.ram[i+8],gb.mmu.ram[i+9],gb.mmu.ram[i+10],gb.mmu.ram[i+11],gb.mmu.ram[i+12],gb.mmu.ram[i+13],gb.mmu.ram[i+14],gb.mmu.ram[i+15])
}
*/
