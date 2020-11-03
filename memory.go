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

func (mem *memory) write(addr uint16, data uint8) {
	mem.ram[addr] = data
}

func (mem *memory) read(addr uint16) uint8 {
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
