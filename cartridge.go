package main

import (
	"io/ioutil"
)

type cartridge struct {
	memory *memory

	ROM  []uint8       //Contains the whole ROM
	ERAM [0x2000]uint8 //External RAM
	mbc  uint8         //Which MBC the rom uses
	rombankNum uint16  //Which rombank is currently in use
}

func initCartridge(memory *memory) *cartridge {
	cart := new(cartridge)
	cart.memory = memory
	//cart.loadRom(fullrom)
	cart.loadRom(testrom)
	cart.rombankNum = 1
	return cart
}

func (cart *cartridge) loadRom(path string){
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	cart.ROM = make([]uint8,len(file))
	for i := 0; i < len(file); i++ {
		cart.ROM[i] = file[i]
	}

	cart.mbc =  cart.ROM[0x0147]
}

func (cart *cartridge) readCartridge(addr uint16) uint8{
	readByte := uint8(0)
	if cart.mbc == 0 {
		//No memory banking
		readByte = cart.ROM[addr]
	} else if cart.mbc == 1  {
		//MBC1
		if inRange(addr,0x0000,0x3FFF){
			readByte = cart.ROM[addr]
		} else {
			//rom bank time!
			readByte = cart.ROM[(cart.rombankNum * 0x4000) + (addr - 0x4000)]
		}
	}

	return readByte
}

func (cart *cartridge) handleRomWrites(addr uint16, data uint8){
	if inRange(addr,0x2000,0x3FFF){
		if data == 0 {
			data = 1
		}
		cart.rombankNum = uint16(data)
	}
}