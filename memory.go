package main

import (
	"io/ioutil"
)

func checkErr(errormsg string, err error) {
	if err != nil {
		panic(errormsg)
	}
}

type memory struct {
	gameboy *gameboy
	ram     [1024 * 64]uint8
}

func initMemory(gb *gameboy) *memory {
	mem := new(memory)
	mem.gameboy = gb
	return mem
}

func (mem *memory) writebyte(addr uint16, data uint8) {
	mem.ram[addr] = data
}

func (mem *memory) writeword(addr uint16, data uint16) {
	//Account for low endian and store lsb first
	mem.ram[addr] = uint8(data) & 0xFF
	mem.ram[addr+1] = uint8(uint16(data) >> 8)
}

func (mem *memory) readbyte(addr uint16) uint8 {
	return mem.ram[addr]
}

func (mem *memory) loadBootrom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr("Could not find bootrom!", err)

	for i := 0; i < len(file); i++ {
		mem.ram[i] = file[i]
	}

}

func (mem *memory) loadRom(path string) {

}
