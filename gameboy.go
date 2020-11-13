package main

import (
	"fmt"
	"log"
	"os"
	// "strconv"
)

type gameboy struct {
	cpu   *gameboyCPU
	mmu   *memory
	debug *debugger
}

func initGameboy(isDebugging bool) *gameboy {
	gb := new(gameboy)
	gb.mmu = initMemory(gb)
	gb.cpu = initCPU(gb)
	if isDebugging {
		gb.debug = initDebugger(gb)
	}

	return gb
}

var cfile string = "11-op a,(hl)"

func initLogging() *os.File {
	//Setup logging
	file, err := os.OpenFile(fmt.Sprintf("logfiles/cpu/%s/%s.txt", cfile, cfile), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(file)

	return file
}

func main() {
	gb := initGameboy(true)
	//gb.mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")
	gb.mmu.loadBlaarg(fmt.Sprintf("roms/testroms/cpu_instrs/%s.gb", cfile))
	//logging := initLogging()
	//defer logging.Close()

	for i := 0; i < 10000000; i++{
		/*
		log.Printf("A: %02X F: %02X B: %02X C: %02X D: %02X E: %02X H: %02X L: %02X SP: %04X PC: 00:%04X (%02X %02X %02X %02X)",
		gb.cpu.getAcc(), gb.cpu.AF&0x00FF, gb.cpu.r8Read[0](), gb.cpu.r8Read[1](), gb.cpu.r8Read[2](), gb.cpu.r8Read[3](), gb.cpu.r8Read[4](), gb.cpu.r8Read[5](), gb.cpu.SP, gb.cpu.PC, gb.mmu.ram[gb.cpu.PC], gb.mmu.ram[gb.cpu.PC+1], gb.mmu.ram[gb.cpu.PC+2], gb.mmu.ram[gb.cpu.PC+3]) */
		//log.Printf("%02X\n",gb.mmu.readbyte(0xDFFA))
		gb.cpu.cycle()
		if gb.mmu.ram[0xFF02] == 0x81 {
			fmt.Printf("%c", gb.mmu.ram[0xFF01])
			gb.mmu.ram[0xFF02] = 0x00
		}
	}
}

//Passed tests
//03
//04
//05
//06
//07
//08
//09
//10
//11
