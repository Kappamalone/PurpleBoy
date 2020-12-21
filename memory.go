package main

import (
	"io/ioutil"
)

type memory struct {
	gb             *gameboy
	cart           *cartridge
	bootromEnabled bool //Used to map bootrom

	bootrom [0x100]uint8 //DMG Bootrom
	//CARTRIDGE: 16KB Rom bank 00 mapped
	//CARTRIDGE: 16KB Rom bank 01~NN mapped
	//CARTRIDGE: 8KB ram bank
	//CARTRIDGE: 4KB WRAM !CGB only bankswitching

	WRAM []uint8     //WRAM
	OAM  [0xA0]uint8 //Object attribute memory aka sprite data
	MMIO [0x80]uint8 //Memory mapped input output
	HRAM [0x7F]uint8 //High ram
}

func initMemory(gb *gameboy, skipBootrom bool) *memory {
	mmu := new(memory)
	mmu.gb = gb
	mmu.bootromEnabled = !skipBootrom

	mmu.initWRAM()
	mmu.cart = initCartridge(mmu)
	if mmu.bootromEnabled {
		mmu.loadBootrom("roms/bootrom/DMG_ROM.gb")
	}
	return mmu
}

func (mmu *memory) initWRAM() {
	//Change this if CGB
	//Initialises 8kb of WRAM
	mmu.WRAM = make([]uint8, 0x2000)
}

func (mmu *memory) executeDMA(data uint8) {
	startAddress := uint16(data) * 0x100
	for i := 0; i < 0xA0; i++ {
		readByte := mmu.readbyte(startAddress + uint16(i))
		mmu.OAM[i] = readByte
	}
}

//MMU LOGIC---------------------------

func (mmu *memory) writebyte(addr uint16, data uint8) {
	//Implements the memory map from the pandocs
	//TODO: make sure to return ppu mode for 0xFF41
	if inRange(addr, 0x0000, 0x7FFF) {
		//16KB ROM Bank 00
		mmu.cart.writeCartridge(addr, data)

	} else if inRange(addr, 0x8000, 0x9FFF) {
		//8KB VRAM
		mmu.gb.ppu.VRAM[addr-0x8000] = data

	} else if inRange(addr, 0xA000, 0xBFFF) {
		//8KB External RAM
		mmu.cart.writeERAM(addr, data)

	} else if inRange(addr, 0xC000, 0xDFFF) {
		//4KB WRAM Bank 0 + 4KB WRAM Bank 1~7
		mmu.WRAM[addr-0xC000] = data

	} else if inRange(addr, 0xE000, 0xFDFF) {
		//ECHO RAM of C000~DDFF
		mmu.WRAM[addr-0xE000] = data

	} else if inRange(addr, 0xFE00, 0xFE9F) {
		//OAM
		mmu.OAM[addr-0xFE00] = data

	} else if inRange(addr, 0xFEA0, 0xFEFF) {
		//Not usable

	} else if inRange(addr, 0xFF00, 0xFF7F) {
		switch addr {
		case 0xFF00:
			mmu.gb.joypad.writeJoypad(data)
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
			mmu.gb.ppu.LCDSTAT = data
		case 0xFF42:
			mmu.gb.ppu.SCY = data
		case 0xFF43:
			mmu.gb.ppu.SCX = data
		case 0xFF44: //LY is read only
		case 0xFF45:
			mmu.gb.ppu.LYC = data
		case 0xFF46:
			mmu.executeDMA(data)
		case 0xFF47:
			mmu.gb.ppu.palette = data
		case 0xFF48:
			mmu.gb.ppu.spritePalette1 = data
		case 0xFF49:
			mmu.gb.ppu.spritePalette2 = data
		case 0xFF4A:
			mmu.gb.ppu.WY = data
		case 0xFF4B:
			mmu.gb.ppu.WX = data
		case 0xFF50:
			mmu.bootromEnabled = (data == 0) //Bootrom writes a non-zero value here to unmap bootrom from memory
		default:
			mmu.MMIO[addr-0xFF00] = data
		}

	} else if inRange(addr, 0xFF80, 0xFFFE) {
		//HRAM
		mmu.HRAM[addr-0xFF80] = data

	} else if inRange(addr, 0xFFFF, 0xFFFF) {
		//IE Register
		mmu.gb.cpu.IE = data
	}

}

func (mmu *memory) readbyte(addr uint16) uint8 {
	readByte := uint8(0)

	if inRange(addr, 0x00, 0xFF) {
		if mmu.bootromEnabled {
			readByte = mmu.bootrom[addr]
		} else {
			readByte = mmu.cart.readCartridge(addr)
		}
	} else if inRange(addr, 0x100, 0x3FFF) {
		//16KB ROM Bank 00
		readByte = mmu.cart.readCartridge(addr)

	} else if inRange(addr, 0x4000, 0x7FFF) {
		//16KB ROM Bank 01~NN
		//Begin the MBC handling
		readByte = mmu.cart.readCartridge(addr)

	} else if inRange(addr, 0x8000, 0x9FFF) {
		//8KB VRAM
		readByte = mmu.gb.ppu.VRAM[addr-0x8000]

	} else if inRange(addr, 0xA000, 0xBFFF) {
		//8KB External RAM
		readByte = mmu.cart.readERAM(addr)

	} else if inRange(addr, 0xC000, 0xDFFF) {
		//4KB WRAM Bank 0 + 4KB WRAM Bank 1~7
		readByte = mmu.WRAM[addr-0xC000]

	} else if inRange(addr, 0xE000, 0xFDFF) {
		//ECHO RAM of C000~DDFF
		readByte = mmu.WRAM[addr-0xE000]

	} else if inRange(addr, 0xFE00, 0xFE9F) {
		//OAM
		readByte = mmu.OAM[addr-0xFE00]

	} else if inRange(addr, 0xFEA0, 0xFEFF) {
		//Not usable

	} else if inRange(addr, 0xFF00, 0xFF7F) {
		switch addr {
		//TIMERS MMIO
		case 0xFF00:
			readByte = mmu.gb.joypad.readJoypad()
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
			readByte = mmu.gb.ppu.palette
		case 0xFF48:
			readByte = mmu.gb.ppu.spritePalette1
		case 0xFF49:
			readByte = mmu.gb.ppu.spritePalette2
		case 0xFF4A:
			readByte = mmu.gb.ppu.WY
		case 0xFF4B:
			readByte = mmu.gb.ppu.WX
		default:
			readByte = mmu.MMIO[addr-0xFF00]
			
		}

	} else if inRange(addr, 0xFF80, 0xFFFE) {
		//HRAM
		readByte = mmu.HRAM[addr-0xFF80]


	} else if inRange(addr, 0xFFFF, 0xFFFF) {
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
