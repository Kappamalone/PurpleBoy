package main

type memory struct {
	ram [16 * 1024]uint8
}

func initMemory() *memory {
	mem := new(memory)
	return mem
}

func (mem *memory) write(addr uint16, data uint8) {
	mem.ram[addr] = data
}

func (mem *memory) read(addr uint16) uint8 {
	return mem.ram[addr]
}
