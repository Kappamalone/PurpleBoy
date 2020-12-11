package main

import (
	"io/ioutil"
)

type memory struct {
	gb  *gameboy
	ram [1024 * 64]uint8

	bootromEnabled bool //Used to map bootrom

	bootrom [0x100]uint8 //DMG Bootrom
	OAM     [0x100]uint8    //Object attribute memory aka sprite data
}

func initMemory(gb *gameboy,skipBootrom bool) *memory {
	mmu := new(memory)
	mmu.gb = gb
	mmu.bootromEnabled = !skipBootrom

	if mmu.bootromEnabled {
		mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")
	}
	//mmu.loadFullRom(fullrom)
	mmu.loadFullRom(testrom)
	return mmu
}

//MMU LOGIC---------------------------

func (mmu *memory) writebyte(addr uint16, data uint8) {
	//Implements the memory map from the pandocs
	//TODO: make sure to return ppu mode for 0xFF41

	if inRange(addr,0x0000,0x3FFF) {
		//16KB ROM Bank 00
		//mmu.ram[addr] = data

	} else if inRange(addr,0x4000,0x7FFF) {
		//16KB ROM Bank 01~NN
		//mmu.ram[addr] = data

	} else if inRange(addr,0x8000,0x9FFF) {
		//8KB VRAM
		mmu.gb.ppu.VRAM[addr - 0x8000] = data

	} else if inRange(addr,0xA000,0xBFFF){
		//8KB External RAM
		mmu.ram[addr] = data

	} else if inRange(addr,0xC000,0xCFFF) {
		//4KB WRAM Bank 0
		mmu.ram[addr] = data

	} else if inRange(addr,0xD000,0xDFFF) {
		//4KB WRAM Bank 1~N
		mmu.ram[addr] = data

	} else if inRange(addr,0xE000,0xFDFF) {
		//ECHO RAM of C000~DDFF
		mmu.ram[addr-0x2000] = data

	} else if inRange(addr,0xFE00,0xFE9F) { 
		//OAM
		mmu.OAM[addr-0xFE00] = data

	} else if inRange(addr,0xFEA0,0xFEFF) {
		//Not usable
		/* I'd log this "erronous" behaviour, however it seems that multiple games do this
		if isDebugging {
			mmu.gb.debug.printConsole("ACCESSING ILLEGAL MEMORY\n", "cyan")
		}*/

	} else if inRange(addr,0xFF00,0xFF7F) {
		switch addr {
		case 0xFF00:
			mmu.ram[addr] = (data & 0xF0) | (mmu.ram[addr] & 0x0F)
		//TIMERS MMIO
		case 0xFF04:
			mmu.gb.cpu.timers.DIV = 0 //Writing any value to DIV resets it to 0
			mmu.gb.cpu.timers.divClock = 0
		case 0xFF05:
			mmu.gb.cpu.timers.TIMA = data
		case 0xFF06:
			mmu.gb.cpu.timers.TMA = data
		case 0xFF07:
			mmu.gb.cpu.timers.TAC = data
		//Interrupt MMIO
		case 0xFF0F:
			mmu.gb.cpu.IF = data 
		//PPU MMIO
		case 0xFF40:
			mmu.gb.ppu.LCDC = data
		case 0xFF41:
			mmu.gb.ppu.LCDSTAT = (data & 0xF8) | uint8(mmu.gb.ppu.mode)
		case 0xFF42:
			mmu.gb.ppu.SCY = data
		case 0xFF43:
			mmu.gb.ppu.SCX = data
		case 0xFF44:
			//LY is read only
		case 0xFF45:
			mmu.gb.ppu.LYC = data
		case 0xFF47:
			mmu.ram[addr] = data
		case 0xFF4A:
			mmu.gb.ppu.WY = data
		case 0xFF4B:
			mmu.gb.ppu.WX = data
		case 0xFF50:
			mmu.bootromEnabled = (data == 0) //Bootrom writes a non-zero value here to unmap bootrom from memory
		default:
			//mmu.gb.debug.logWrite(addr)
			mmu.ram[addr] = data
		}

	} else if inRange(addr,0xFF80,0xFFFE){
		//HRAM
		mmu.ram[addr] = data

	} else if inRange(addr,0xFFFF,0xFFFF){
		//IE Register
		mmu.gb.cpu.IE = data
	}

}

func (mmu *memory) readbyte(addr uint16) uint8 {
	readByte := uint8(0)

	if inRange(addr,0x00,0xFF)  {
		if mmu.bootromEnabled {
			readByte = mmu.bootrom[addr]
		} else {
			readByte = mmu.ram[addr]
		}
	} else if inRange(addr,0x100,0x3FFF){
		//16KB ROM Bank 00
		readByte = mmu.ram[addr]

	} else if inRange(addr,0x4000,0x7FFF){
		//16KB ROM Bank 01~NN
		readByte = mmu.ram[addr]

	} else if inRange(addr,0x8000,0x9FFF){
		//8KB VRAM
		//readByte = mmu.gb.ppu.readVRAM(addr - 0x8000)
		readByte = mmu.gb.ppu.VRAM[addr - 0x8000]

	} else if inRange(addr,0xA000,0xBFFF){
		//8KB External RAM
		readByte = mmu.ram[addr]

	} else if inRange(addr,0xC000,0xCFFF){
		//4KB WRAM Bank 0
		readByte = mmu.ram[addr]

	} else if inRange(addr,0xD000,0xDFFF){
		//4KB WRAM Bank 1~N
		readByte = mmu.ram[addr]

	} else if inRange(addr,0xE000,0xFDFF){
		//ECHO RAM of C000~DDFF
		readByte = mmu.ram[addr-0x2000]

	} else if inRange(addr,0xFE00,0xFE9F){
		//OAM
		readByte = mmu.OAM[addr-0xFE00]

	} else if inRange(addr,0xFEA0,0xFEFF){
		//Not usable
		if isDebugging {
			mmu.gb.debug.printConsole("ACCESSING ILLEGAL MEMORY\n", "cyan")
		}

	} else if inRange(addr,0xFF00,0xFF7F){
		switch addr {
		//TIMERS MMIO
		case 0xFF04:
			readByte = mmu.gb.cpu.timers.DIV
		case 0xFF05:
			readByte = mmu.gb.cpu.timers.TIMA
		case 0xFF06:
			readByte = mmu.gb.cpu.timers.TMA
		case 0xFF07:
			readByte = mmu.gb.cpu.timers.TAC
		//Interrupt MMIO
		case 0xFF0F:
			readByte = mmu.gb.cpu.IF
		//PPU MMIO
		case 0xFF40:
			readByte = mmu.gb.ppu.LCDC
		case 0xFF41:
			readByte = (mmu.gb.ppu.LCDSTAT & 0xFC) | uint8(mmu.gb.ppu.mode) 
		case 0xFF42:
			readByte = mmu.gb.ppu.SCY
		case 0xFF43:
			readByte = mmu.gb.ppu.SCX
		case 0xFF44:
			readByte = mmu.gb.ppu.LY
		case 0xFF45:
			readByte = mmu.gb.ppu.LYC
		case 0xFF47:
			readByte = mmu.ram[addr]
		case 0xFF4A:
			readByte = mmu.gb.ppu.WY
		case 0xFF4B:
			readByte = mmu.gb.ppu.WX
		default:
			//mmu.gb.debug.logRead(addr)
			readByte = mmu.ram[addr]
		}

	} else if inRange(addr,0xFF80,0xFFFE){
		//HRAM
		readByte = mmu.ram[addr]

	} else if inRange(addr,0xFFFF,0xFFFF){
		//IE Register
		readByte = mmu.gb.cpu.IE
	}


	return readByte

}

func (mmu *memory) readWord(addr uint16) uint16 {
	//Account for low endianness and read lsb first
	low := mmu.readbyte(addr)
	hi := mmu.readbyte(addr + 1)
	return uint16(hi)<<8 | uint16(low)
}

func (mmu *memory) writeword(addr uint16, data uint16) {
	//Account for low endian and store lsb first
	low := uint8(data & 0x00FF)
	hi := uint8((data & 0xFF00) >> 8)
	mmu.writebyte(addr, low)
	mmu.writebyte(addr+1, hi)
}

//ROM LOADING--------------------------
func (mmu *memory) loadBootrom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find bootrom!")

	for i := 0; i < len(file); i++ {
		mmu.bootrom[i] = file[i]
	}

}

func (mmu *memory) loadFullRom(path string) {
	file, err := ioutil.ReadFile(path)
	checkErr(err, "Could not find rom specified!")

	for i := 0; i < len(file); i++ {
		mmu.ram[i] = file[i]
	}
}
