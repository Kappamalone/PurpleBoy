package main

import (
	"io/ioutil"
)

func checkErr(err error, errormsg string) {
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
	mem.ram[addr] = uint8(data & 0x00FF)
	mem.ram[addr+1] = uint8((data & 0xFF00) >> 8)
}

func (mem *memory) readbyte(addr uint16) uint8 {
	return mem.ram[addr]
}

func (mem *memory) readWord(addr uint16) uint16 {
	//Account for low endian
	return uint16(mem.ram[addr+1])<<8 | uint16(mem.ram[addr])
}

func (mem *memory) loadBootrom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find bootrom!")

	for i := 0; i < len(file); i++ {
		mem.ram[i] = file[i]
	}

}

func (mem *memory) loadBlaarg(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	for i := 0; i < len(file); i++ {
		mem.ram[i] = file[i]
	}
}

func (mem *memory) loadRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	for i := 0; i < len(file); i++ {
		mem.ram[0x100+i] = file[i]
	}
}
