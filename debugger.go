package main

import (
	"fmt"
	"log"
	"os"
)

type debugger struct {
	gb *gameboy
}

func initDebugger(gb *gameboy, isLogging bool) *debugger {
	debug := new(debugger)
	debug.gb = gb
	if isLogging {
		initLogging()
	}

	return debug
}

func initLogging() {
	//Setup logging
	file, err := os.OpenFile(fmt.Sprintf("logfiles/cpu/%s/%s.txt", cfile, cfile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(file)
}

func (debug *debugger) logTrace() {
	log.Printf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: 00:%04X (%02X %02X %02X %02X)", debug.gb.cpu.getAcc(), debug.gb.cpu.AF&0x00FF, debug.gb.cpu.r8Read[0](), debug.gb.cpu.r8Read[1](), debug.gb.cpu.r8Read[2](), debug.gb.cpu.r8Read[3](), debug.gb.cpu.r8Read[4](), debug.gb.cpu.r8Read[5](), debug.gb.cpu.SP, debug.gb.cpu.PC, debug.gb.mmu.ram[debug.gb.cpu.PC], debug.gb.mmu.ram[debug.gb.cpu.PC+1], debug.gb.mmu.ram[debug.gb.cpu.PC+2], debug.gb.mmu.ram[debug.gb.cpu.PC+3])
}

/*
for i := 0; i < 256; i += 16{
	fmt.Printf("%02X %02X %02X %02X %002X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n",gb.mmu.ram[i],gb.mmu.ram[i+1],gb.mmu.ram[i+2],gb.mmu.ram[i+3],gb.mmu.ram[i+4],gb.mmu.ram[i+5],gb.mmu.ram[i+6],gb.mmu.ram[i+7],gb.mmu.ram[i+8],gb.mmu.ram[i+9],gb.mmu.ram[i+10],gb.mmu.ram[i+11],gb.mmu.ram[i+12],gb.mmu.ram[i+13],gb.mmu.ram[i+14],gb.mmu.ram[i+15])
}
*/

//Write debug windows down here
