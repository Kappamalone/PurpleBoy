package main

import (
	"io/ioutil"
	"fmt"
	"strings"
)

var (
	mbc1BitmaskMap = [7]uint8{0x00, 0x3, 0x7, 0xF, 0x1F, 0x1F, 0x1F} //Get bitmasks for different rom sizes on MBC1
	RAMSizes      = [4]int{0x0000, 0x0500, 0x2000, 0x8000}
)

type cartridge struct {
	memory *memory

	ROM         []uint8 //Contains the whole ROM
	MBC         uint8   //Which MBC the rom uses
	rombankNum  int     //Which rombank is currently in use
	special2Bit int     //Multi purpose 2 bits used as RAM bank num or Upper bits of rom bank num
	bankMode    uint8   //Which rom banking mode is in use

	ROMSize  uint8  //Rom size specified for cartridge
	ERAMSize uint8  //Ram size specified for cartridge
	usingBBRAM bool //Whether or not battery buffered ram is being used
	title    string //Title of game

	ERAM       []uint8 //External RAM
	ERAMEnable bool    //Used to enable/disable eram
}

func initCartridge(memory *memory) *cartridge {
	//TODO: Battery buffered RAM
	//While I could've written this using interfaces,
	//It would've been a lot of extra loc, and since
	//I'm not adding support for every MBC, I'll keep it as is
	cart := new(cartridge)
	cart.memory = memory
	cart.rombankNum = 1

	if useTestRom {
		cart.loadRom(testrom)
	} else {
		cart.loadRom(gamerom)
	}
	cart.initERAM()

	return cart
}

func getMBCNum(hexvalue uint8) uint8 {
	mbcNum := uint8(0xFF)
	switch hexvalue {
	case 0x0:
		mbcNum = 0
	case 0x1, 0x2, 0x3:
		mbcNum = 1
	case 0x5, 0x6:
		mbcNum = 2
	case 0xF, 0x10, 0x11, 0x12, 0x13:
		mbcNum = 3
	}
	if mbcNum == 0xFF {
		println(hexvalue)
		panic("Unsupported MBC Type")
	}
	return mbcNum
}

func getBBRAM(hexvalue uint8) bool {
	usingBBRAM := false 
	switch hexvalue {
	case 0x3,0x6,0x9,0xF,0x10,0x13:
		usingBBRAM = true
	}
	return usingBBRAM
}

func (cart *cartridge) initERAM() {
	//Depending on RAM num initialise properly sized ERAM
	if cart.MBC == 1 || cart.MBC == 3 {
		cart.ERAM = make([]uint8, RAMSizes[cart.ERAMSize])
	} else if cart.MBC == 2 {
		//MBC2 has 0x800 bits of built in ERAM
		cart.ERAM = make([]uint8, 0x800)
	}
	cart.loadBBRAM()
}

func (cart *cartridge) loadBBRAM(){
	if !cart.usingBBRAM {
		return
	}
	path := fmt.Sprintf("roms/gameroms/%s.sav", cart.title)
	file, err := ioutil.ReadFile(path)

	if err == nil {
		//Only load BBRAM if BBRAM exists
		for i := 0; i < len(file); i++ {
			cart.ERAM[i] = file[i]
		}
	}
}

func (cart *cartridge) saveBBRAM(){
	if !cart.usingBBRAM {
		return
	}
	path := fmt.Sprintf("roms/gameroms/%s.sav", cart.title)
	err := ioutil.WriteFile(path,cart.ERAM,0644) //Golang really out here making my life this easy
	checkErr(err,"Could not store battery buffered ram!")
}

func (cart *cartridge) loadRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	cart.ROM = make([]uint8, len(file))
	for i := 0; i < len(file); i++ {
		cart.ROM[i] = file[i]
	}

	cart.MBC = getMBCNum(cart.ROM[0x147])
	cart.usingBBRAM = getBBRAM(cart.ROM[0x147])
	cart.ROMSize = cart.ROM[0x0148]
	cart.ERAMSize = cart.ROM[0x149]

	//Gets title of game from memory
	chars := make([]string, 0)
	for i := 0; i < 16; i++ {
		char := cart.ROM[0x134 + i]
		if char != 0 {
			chars = append(chars, fmt.Sprintf("%c", char))
		}
	}
	cart.title = strings.Join(chars,"")
}

func (cart *cartridge) readCartridge(addr uint16) uint8 {
	readByte := uint8(0)
	//println(cart.ROM[0x3F * 0x4000 + 0x3FFF])
	if cart.MBC == 0 {
		//No memory banking
		readByte = cart.ROM[addr]
	} else if cart.MBC == 1 {
		//MBC1
		if inRange(addr, 0x0000, 0x3FFF) {
			if cart.bankMode == 0x00 || cart.ROMSize <= 0x4 { //Use regular banking if romsize is 512KBytes or lower
				readByte = cart.ROM[addr]
			} else {
				//The 2 special bits can map to 0x00,0x20,0x40,0x60 banks
				if cart.ROMSize == 0x5 {
					//Wrap the additional bit
					readByte = cart.ROM[((cart.special2Bit&1)*0x20*0x4000)+int(addr)]
				} else {
					readByte = cart.ROM[((cart.special2Bit)*0x20*0x4000)+int(addr)]
				}
			}
		} else {
			if cart.ROMSize <= 0x4 { //Use regular banking if rom is <= 512 KBytes or lower
				//Simple Rom banking
				readByte = cart.ROM[(cart.rombankNum*0x4000)+(int(addr)-0x4000)]
			} else {
				//Advanced Rom banking
				if cart.ROMSize == 0x5 {
					//Wrap the additional bit
					readByte = cart.ROM[((cart.special2Bit&1)<<5|cart.rombankNum)*0x4000+(int(addr)-0x4000)]
				} else {
					readByte = cart.ROM[(cart.special2Bit<<5|cart.rombankNum)*0x4000+(int(addr)-0x4000)]
				}
			}
		}
	} else if cart.MBC == 2 {
		//MBC2
		if inRange(addr, 0x0000, 0x3FFF) {
			readByte = cart.ROM[addr]
		} else {
			readByte = cart.ROM[(cart.rombankNum*0x4000)+(int(addr)-0x4000)]
		}
	} else if cart.MBC == 3 {
		//MBC3
		if inRange(addr, 0x0000, 0x3FFF) {
			readByte = cart.ROM[addr]
		} else {
			//Simple Rom banking
			readByte = cart.ROM[(cart.rombankNum*0x4000)+(int(addr)-0x4000)]
		}
	}

	return readByte
}

func (cart *cartridge) writeCartridge(addr uint16, data uint8) {
	if inRange(addr, 0x0000, 0x1FFF) {
		//ERAM enable/disable
		if cart.MBC == 1 || cart.MBC == 3 {
			cart.ERAMEnable = (data & 0xF) == 0xA
		} else if cart.MBC == 2 {
			//Handle writes
			if !bitSet(addr, 8) {
				//RAM Write
				cart.ERAMEnable = (data & 0xF) == 0xA
			} else {
				//ROM Write
				if data == 0 {
					data = 1
				}
				data &= mbc1BitmaskMap[cart.ROMSize]
				cart.rombankNum = int(data)
			}
		}
	} else if inRange(addr, 0x2000, 0x3FFF) {
		//ROM Bank select
		if cart.MBC == 1 {
				if data&0x1F == 0 {
					//0x00,0x20,0x40,0x60 Get remapped to one rom bank higher
					//Which means any byte where the lower 5 bits are 0 get mapped to one a rombank one higher
					data++
				}
			data &= mbc1BitmaskMap[cart.ROMSize]
			cart.rombankNum = int(data)
		} else if cart.MBC == 2 {
			//Handle writes
			if !bitSet(addr, 8) {
				//RAM Write
				cart.ERAMEnable = (data & 0xF) == 0xA
			} else {
				//ROM Write
				if data == 0 {
					data = 1
				}
				data &= mbc1BitmaskMap[cart.ROMSize]
				cart.rombankNum = int(data)
			}
		} else if cart.MBC == 3 {
			//0x20,0x40 and 0x60 aren't affected in MBC3
			//Meaning we don't have to apply the mask
			if data == 0 {
				data = 1
			}
			cart.rombankNum = int(data)
		}

	} else if inRange(addr, 0x4000, 0x5FFF) {
		//Used to control RAM/ROM 5th/6th bits
		cart.special2Bit = int(data) & 0x3

	} else if inRange(addr, 0x6000, 0x7FFF) {
		//Banking mode select
		cart.bankMode = data & 0x1
	}
}

func (cart *cartridge) readERAM(addr uint16) uint8 {
	readByte := uint8(0xFF)
	if cart.MBC == 1 || cart.MBC == 3 {
		if cart.ERAMEnable && cart.ERAMSize != 0 {
			//ERAM sizes 8kb and lower don't have any banking
			if cart.bankMode == 0 || cart.ERAMSize <= 0x02 {
				//ERAM Bank 0
				readByte = cart.ERAM[addr-0xA000]

			} else {
				//ERAM Banks 0-4
				readByte = cart.ERAM[(cart.special2Bit*0x2000)+(int(addr)-0xA000)]
			}
		}
	} else if cart.MBC == 2 {
		//MBC2 has 0x800 bits of built in ERAM
		readByte = cart.ERAM[(addr-0xA000)%800]
	}

	return readByte
}

func (cart *cartridge) writeERAM(addr uint16, data uint8) {
	if cart.MBC == 1 || cart.MBC == 3 {
		if cart.ERAMEnable && cart.ERAMSize != 0 {
			//ERAM sizes 8kb and lower don't have any banking
			if cart.bankMode == 0 || cart.ERAMSize <= 0x02 {
				//ERAM Bank 0
				//TODO: Use modulus to account for 2kb rom banks
				cart.ERAM[addr-0xA000] = data
			} else {
				//ERAM Banks 0-4
				cart.ERAM[(cart.special2Bit*0x2000)+(int(addr)-0xA000)] = data
			}
		}
	} else if cart.MBC == 2 {
		cart.ERAM[(addr-0xA000)%0x800] = data
	}
}
