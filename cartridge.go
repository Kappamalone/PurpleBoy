package main

import (
	"io/ioutil"
)

type cartridge struct {
	memory *memory

	ROM         []uint8 //Contains the whole ROM
	MBC         uint8   //Which MBC the rom uses
	rombankNum  uint16  //Which rombank is currently in use
	special2Bit uint16   //Multi purpose 2 bits used as RAM bank num or Upper bits of rom bank num
	bankMode    uint8   //Which rom banking mode is in use

	ERAM       [0x2000]uint8 //External RAM
	ERAMEnable bool          //Used to enable/disable eram
}

func initCartridge(memory *memory) *cartridge {
	cart := new(cartridge)
	cart.memory = memory
	//cart.loadRom(fullrom)
	cart.loadRom(testrom)
	cart.rombankNum = 1

	return cart
}

func (cart *cartridge) initERAM(){
	//Depending on MBC initialise properly sized ERAM
}

func (cart *cartridge) loadRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	cart.ROM = make([]uint8, len(file))
	for i := 0; i < len(file); i++ {
		cart.ROM[i] = file[i]
	}

	cart.MBC = cart.ROM[0x0147]
}

func (cart *cartridge) readCartridge(addr uint16) uint8 {
	readByte := uint8(0)
	if cart.MBC == 0 {
		//No memory banking
		readByte = cart.ROM[addr]
	} else if cart.MBC == 1 {
		//MBC1
		if inRange(addr, 0x0000, 0x3FFF) {
			if cart.bankMode == 0x00 {
				readByte = cart.ROM[addr]
			} else {
				//The 2 special bits can map to 0x00,0x20,0x40,0x60 banks
				readByte = cart.ROM[(uint32(cart.special2Bit) * 0x20 * 0x4000) + uint32(addr)]
			}
		} else {
			//rom bank time!
			if cart.bankMode == 0x00 {
				//Simple Rom banking
				readByte = cart.ROM[(cart.rombankNum*0x4000)+(addr-0x4000)]
			} else {
				//Advanced Rom banking
				readByte = cart.ROM[(cart.special2Bit<<5|cart.rombankNum)*0x4000+(addr-0x4000)]
			}
		}
	}

	return readByte
}

func (cart *cartridge) handleRomWrites(addr uint16, data uint8) {
	if inRange(addr, 0x0000, 0x1FFF) {
		//ERAM enable/disable
		if addr == 0x0A {
			cart.ERAMEnable = true
		} else {
			cart.ERAMEnable = false
		}
	} else if inRange(addr, 0x2000, 0x3FFF) {
		//ROM Bank select
		if data & 0x1F == 0 {
			//0x00,0x20,0x40,0x60 Get remapped to one rom bank higher
			data++
		}
		cart.rombankNum = uint16(data)

	} else if inRange(addr, 0x4000, 0x5FFF) {
		//Why is the writes here using a value greater than 0x3?
		if data > 0x3 {
			//cart.memory.gb.debug.printConsole("ROM write for special 2 bits too :beeg\n","red")
		}
		cart.special2Bit = uint16(data & 0x3)

	} else if inRange(addr, 0x6000, 0x7FFF) {
		//Banking mode select
		if data > 0x1 {
			cart.memory.gb.debug.printConsole("ROM write for banking mode select incorrect!\n", "red")
		}

		cart.bankMode = data
	}
}


func (cart *cartridge) readERAM(addr uint16) uint8 {
	readByte := uint8(0)
	if cart.ERAMEnable {
		if cart.bankMode == 0 || len(cart.ERAM) == 0x2000 {
			//ERAM Bank 0
			readByte = cart.ERAM[addr - 0xA000]
		} else {
			//ERAM Banks 0-4
			readByte = cart.ERAM[(cart.special2Bit * 0x2000) + (addr - 0xA000)]
		}
	}

	return readByte
}

func (cart *cartridge) writeERAM(addr uint16,data uint8){
	if cart.ERAMEnable {
		if cart.bankMode == 0 {
			//ERAM Bank 0
			cart.ERAM[addr - 0xA000] = data
		} else {
			//ERAM Banks 0-4
			cart.ERAM[(cart.special2Bit * 0x2000) + (addr - 0xA000)] = data
		}
	}
}