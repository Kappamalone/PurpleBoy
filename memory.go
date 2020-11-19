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
	//Implements the memory map from the pandocs

	if addr >= 0x0000 && addr <= 0x3FFF {
		//16KB ROM Bank 00
		mmu.ram[addr] = data

	} else if addr >= 0x4000 && addr <= 0x7FFF {
		//16KB ROM Bank 01~NN
		mmu.ram[addr] = data

	} else if addr >= 0x8000 && addr <= 0x9FFF {
		//8KB VRAM
		mmu.gb.ppu.VRAM[addr - 0x8000] = data

	} else if addr >= 0xA000 && addr <= 0xBFFF {
		//8KB External RAM
		mmu.ram[addr] = data

	} else if addr >= 0xC000 && addr <= 0xCFFF {
		//4KB WRAM Bank 0
		mmu.ram[addr] = data

	} else if addr >= 0xD000 && addr <= 0xDFFF {
		//4KB WRAM Bank 1~N
		mmu.ram[addr] = data

	} else if addr >= 0xE000 && addr <= 0xFDFF {
		//ECHO RAM of C000~DDFF
		mmu.ram[addr - 0x2000] = data

	} else if addr >= 0xFE00 && addr <= 0xFE9F {
		//OAM
		mmu.ram[addr] = data

	} else if addr >= 0xFEA0 && addr <= 0xFEFF {
		//Not usable
		if isDebugging {
			mmu.gb.debug.printConsole("ACCESSING ILLEGAL MEMORY\n","green")
		}

	} else if addr >= 0xFF00 && addr <= 0xFF7F {
		//IO Registers
		mmu.ram[addr] = data

	} else if addr >= 0xFF80 && addr <= 0xFFFE {
		//HRAM
		mmu.ram[addr] = data

	} else if addr >= 0xFFFF && addr <= 0xFFFF {
		//IE Register
		mmu.ram[0xFFFF] = data
	}

}

func (mmu *memory) readbyte(addr uint16) uint8 {
	//Time to do some cool mmu stuff

	var readByte uint8 = 0

	if addr >= 0x0000 && addr <= 0x3FFF {
		//16KB ROM Bank 00
		readByte = mmu.ram[addr]

	} else if addr >= 0x4000 && addr <= 0x7FFF {
		//16KB ROM Bank 01~NN
		readByte = mmu.ram[addr]

	} else if addr >= 0x8000 && addr <= 0x9FFF {
		//8KB VRAM
		readByte = mmu.gb.ppu.VRAM[addr - 0x8000]

	} else if addr >= 0xA000 && addr <= 0xBFFF {
		//8KB External RAM
		readByte = mmu.ram[addr]

	} else if addr >= 0xC000 && addr <= 0xCFFF {
		//4KB WRAM Bank 0
		readByte = mmu.ram[addr]

	} else if addr >= 0xD000 && addr <= 0xDFFF {
		//4KB WRAM Bank 1~N
		readByte = mmu.ram[addr]

	} else if addr >= 0xE000 && addr <= 0xFDFF {
		//ECHO RAM of C000~DDFF
		readByte = mmu.ram[addr - 0x2000]

	} else if addr >= 0xFE00 && addr <= 0xFE9F {
		//OAM
		readByte = mmu.ram[addr]

	} else if addr >= 0xFEA0 && addr <= 0xFEFF {
		//Not usable
		if isDebugging {
			mmu.gb.debug.printConsole("ACCESSING ILLEGAL MEMORY\n","green")
		}

	} else if addr >= 0xFF00 && addr <= 0xFF7F {
		//IO Registers
		readByte = mmu.ram[addr]

	} else if addr >= 0xFF80 && addr <= 0xFFFE {
		//HRAM
		readByte = mmu.ram[addr]

	} else if addr >= 0xFFFF && addr <= 0xFFFF {
		//IE Register
		readByte = mmu.ram[0xFFFF]
	}

	return readByte 


}

func (mmu *memory) readWord(addr uint16) uint16 {
	//Account for low endianness
	low := mmu.readbyte(addr)
	hi := mmu.readbyte(addr + 1)
	return uint16(hi)<<8 | uint16(low)
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