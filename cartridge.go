package main

import (
	"io/ioutil"
)

var (
	mbc1BitmaskMap = [6]uint8{0x3, 0x7, 0xF, 0x1F, 0x1F, 0x1F} //Get bitmasks for different rom sizes on MBC1
	RAMSizes       = [4]int{0x0000, 0x0500, 0x2000, 0x8000}
)

type cartridge struct {
	memory *memory

	ROM         []uint8 //Contains the whole ROM
	MBC         uint8   //Which MBC the rom uses
	rombankNum  uint16  //Which rombank is currently in use
	special2Bit uint16  //Multi purpose 2 bits used as RAM bank num or Upper bits of rom bank num
	bankMode    uint8   //Which rom banking mode is in use

	ROMSize  uint8 //Rom size specified for cartridge
	ERAMSize uint8 //Ram size specified for cartridge

	ERAM       []uint8 //External RAM
	ERAMEnable bool    //Used to enable/disable eram
}

func initCartridge(memory *memory) *cartridge {
	cart := new(cartridge)
	cart.memory = memory
	cart.rombankNum = 1

	if useTestRom {
		cart.loadRom(testrom)
	} else {
		cart.loadRom(gamerom)
	}
	//cart.loadRom(testrom)
	cart.initERAM()

	return cart
}

func getMBCNum(hexvalue uint8) uint8 {
	mbcNum := uint8(0)
	switch hexvalue {
	case 0x1, 0x2, 0x3:
		mbcNum = 1
	case 0xF, 0x10, 0x11, 0x12, 0x13:
		mbcNum = 3
	}
	return mbcNum
}

func (cart *cartridge) initERAM() {
	//Depending on RAM num initialise properly sized ERAM
	cart.ERAM = make([]uint8, RAMSizes[cart.ERAMSize])
}

func (cart *cartridge) loadRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	cart.ROM = make([]uint8, len(file))
	for i := 0; i < len(file); i++ {
		cart.ROM[i] = file[i]
	}

	cart.MBC = getMBCNum(cart.ROM[0x147])
	cart.ROMSize = cart.ROM[0x0148]
	cart.ERAMSize = cart.ROM[0x149]
}

func (cart *cartridge) readCartridge(addr uint16) uint8 {
	readByte := uint8(0)
	if cart.MBC == 0 {
		//No memory banking
		readByte = cart.ROM[addr]
	} else if cart.MBC == 1 {
		//MBC1
		//Roms bigger than romsize 5 are brokey :(
		if inRange(addr, 0x0000, 0x3FFF) {
			if cart.bankMode == 0x00 || cart.ROMSize <= 0x5 { //Use regular banking if romsize is 512KBytes or lower
				readByte = cart.ROM[addr]
			} else {
				//The 2 special bits can map to 0x00,0x20,0x40,0x60 banks
				readByte = cart.ROM[(uint32(cart.special2Bit)*0x20*0x4000)+uint32(addr)]
			}
		} else {
			//rom bank time!
			if cart.bankMode == 0x00 || cart.ROMSize <= 0x5 { //Use regular banking if rom is <= 512 KBytes or lower
				//Simple Rom banking
				cart.memory.gb.debug.printConsole("Ha found you!\n","red")
				readByte = cart.ROM[(cart.rombankNum*0x4000)+(addr-0x4000)]
			} else {
				//Advanced Rom banking
				readByte = cart.ROM[(cart.special2Bit<<5|cart.rombankNum)*0x4000+(addr-0x4000)]
			}
		}
	} else if cart.MBC == 3 {
		//MBC3
	}

	return readByte
}

func (cart *cartridge) handleRomWrites(addr uint16, data uint8) {
	if inRange(addr, 0x0000, 0x1FFF) {
		//ERAM enable/disable
		if (data & 0xF) == 0xA {
			cart.ERAMEnable = true
		} else {
			cart.ERAMEnable = false
		}
	} else if inRange(addr, 0x2000, 0x3FFF) {
		//ROM Bank select
		if cart.MBC == 1 {
			if (data & 0x1F) == 0 {
				//0x00,0x20,0x40,0x60 Get remapped to one rom bank higher
				//Which means any byte where the lower 5 bits are 0 get mapped to one a rombank one higher
				data++
			}
		}

		cart.rombankNum = uint16(data & mbc1BitmaskMap[cart.ROMSize])

	} else if inRange(addr, 0x4000, 0x5FFF) {
		cart.special2Bit = uint16(data & 0x3)

	} else if inRange(addr, 0x6000, 0x7FFF) {
		//Banking mode select
		cart.bankMode = data & 0x1
	}
}

func (cart *cartridge) readERAM(addr uint16) uint8 {
	readByte := uint8(0xFF)
	if cart.ERAMEnable && cart.ERAMSize != 0 {
		//ERAM sizes 8kb and lower don't have any banking
		if cart.bankMode == 0 || cart.ERAMSize == 0x02 {
			//ERAM Bank 0
			readByte = cart.ERAM[addr-0xA000]
			
		} else {
			//ERAM Banks 0-4
			readByte = cart.ERAM[(cart.special2Bit*0x2000)+(addr-0xA000)]
		}
	}


	return readByte
}

func (cart *cartridge) writeERAM(addr uint16, data uint8) {
	if cart.ERAMEnable && cart.ERAMSize != 0 {
		//ERAM sizes 8kb and lower don't have any banking
		if cart.bankMode == 0 || cart.ERAMSize == 0x02 {
			//ERAM Bank 0
			//TODO: Use modulus to account for 2kb rom banks
			cart.ERAM[addr-0xA000] = data
		} else {
			//ERAM Banks 0-4
			cart.ERAM[(cart.special2Bit*0x2000)+(addr-0xA000)] = data
		}
	}
}
