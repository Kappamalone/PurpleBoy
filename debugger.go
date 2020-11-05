package main

import ()

type debugger struct {
	gb *gameboy
}

/*
for i := 0; i < 256; i += 16{
	fmt.Printf("%02X %02X %02X %02X %002X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n",gb.mmu.ram[i],gb.mmu.ram[i+1],gb.mmu.ram[i+2],gb.mmu.ram[i+3],gb.mmu.ram[i+4],gb.mmu.ram[i+5],gb.mmu.ram[i+6],gb.mmu.ram[i+7],gb.mmu.ram[i+8],gb.mmu.ram[i+9],gb.mmu.ram[i+10],gb.mmu.ram[i+11],gb.mmu.ram[i+12],gb.mmu.ram[i+13],gb.mmu.ram[i+14],gb.mmu.ram[i+15])
}
*/

func initDebugger(gb *gameboy) *debugger {
	debug := new(debugger)
	debug.gb = gb

	return debug
}
