package main

type gameboyCPU struct {
	gb *gameboy

	AF uint16
	BC uint16
	DE uint16
	HL uint16

	SP uint16
	PC uint16

	cycles          int
	currInstruction string

	r8Read    map[uint8]func() uint8 //r8 group 1 reads
	r8Write   map[uint8]func(uint8)  //r8 group 1 write
	r16group1 map[uint8]*uint16      //r16 group 1
	r16group2 map[uint8]*uint16      //r16 group 2
	r16group3 map[uint8]*uint16      //r16 group 3
}

func test() uint8 {
	return uint8(0)
}

type fn func() uint8

func initCPU(gb *gameboy) *gameboyCPU {
	cpu := new(gameboyCPU)
	cpu.currInstruction = "LD SP, d16"
	cpu.gb = gb

	//Some premptively messy code to make other parts easier to write :)
	cpu.r8Read = map[uint8]func() uint8{
		0: cpu.B, 1: cpu.C,
		2: cpu.D, 3: cpu.E,
		4: cpu.H, 5: cpu.L,
		6: func() uint8 { return cpu.gb.mmu.readbyte(cpu.HL) }, 7: cpu.A,
	}

	cpu.r8Write = map[uint8]func(uint8){
		0: cpu.setB, 1: cpu.setC,
		2: cpu.setD, 3: cpu.setD,
		4: cpu.setH, 5: cpu.setL,
		6: func(data uint8) { cpu.gb.mmu.writebyte(cpu.HL, data) }, 7: cpu.setA,
	}

	cpu.r16group1 = map[uint8]*uint16{
		0: &cpu.BC, 1: &cpu.DE,
		2: &cpu.HL, 3: &cpu.SP,
	}

	cpu.r16group2 = map[uint8]*uint16{
		0: &cpu.BC, 1: &cpu.DE,
		2: &cpu.HL, 3: &cpu.HL, //--> Remember to increment HL on 2 and decrement HL on 3
	}

	cpu.r16group3 = map[uint8]*uint16{
		0: &cpu.BC, 1: &cpu.DE,
		2: &cpu.HL, 3: &cpu.AF,
	}

	return cpu
}

func (cpu *gameboyCPU) cycle() {
	if cpu.cycles == 0 {
		fetchedInstruction := cpu.gb.mmu.readbyte(cpu.PC)
		cpu.PC++

		cpu.decodeAndExecute(fetchedInstruction)
	}
	cpu.cycles--
}

func (cpu *gameboyCPU) decodeAndExecute(opcode uint8) {
	//I've structured the elifs to correctly decode the correct instruction
	//It looks a bit wonky since I'm making sure everything that isn't an operand
	//Exactly matches the decodings

	
	
}
