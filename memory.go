package main

import (
	"io/ioutil"
)

type memory struct {
	gb  *gameboy
	ram [1024 * 64]uint8
}

func initMemory(gb *gameboy) *memory {
	mmu := new(memory)
	mmu.gb = gb
	return mmu
}

//MMU LOGIC---------------------------

func (mmu *memory) writebyte(addr uint16, data uint8) {
	mmu.ram[addr] = data
}

func (mmu *memory) readbyte(addr uint16) uint8 {
	//Time to do some cool mmu stuff
	return mmu.ram[addr]
}

func (mmu *memory) readWord(addr uint16) uint16 {
	//Account for low endian
	return uint16(mmu.ram[addr+1])<<8 | uint16(mmu.ram[addr])
}

func (mmu *memory) writeword(addr uint16, data uint16) {
	//Account for low endian and store lsb first
	mmu.ram[addr] = uint8(data & 0x00FF)
	mmu.ram[addr+1] = uint8((data & 0xFF00) >> 8)
}

//ROM LOADING--------------------------
func (mmu *memory) loadBootrom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find bootrom!")

	for i := 0; i < len(file); i++ {
		mmu.ram[i] = file[i]
	}

}

func (mmu *memory) loadBlaarg(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	for i := 0; i < len(file); i++ {
		mmu.ram[i] = file[i]
	}
}

func (mmu *memory) loadRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	for i := 0; i < len(file); i++ {
		mmu.ram[0x100+i] = file[i]
	}
}
