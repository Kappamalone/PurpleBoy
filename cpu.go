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
	//So most bit specific opcodes are first in their respective groups
	//Which I've marked

	//First Group-----------------------------------
	if opcode == 0x00 {
		//NOP: Do nothing
		cpu.currInstruction = "NOP"

	} else if opcode == 0x08 {
		//LD (u16),SP
		cpu.gb.mmu.writeword(cpu.d16(), cpu.SP)

		cpu.currInstruction = "LD (u16), SP"

	} else if opcode == 0x10 {
		//STOP: Do nothing, Extra spicy flavoured

		cpu.d8() //Since for some reason this instruction is 2 bytes
		cpu.currInstruction = "STOP"

	} else if opcode == 0x18 {
		//JR (unconditional)

		cpu.currInstruction = "JR (unconditional)"

	} else if opcode&0x20 == 0x20 {
		//JR (conditional)

		cpu.currInstruction = "JR (conditional)"

	//Second group---------------------------
	} else if opcode&0x7 == 0x7 {
		//OPCODE GROUP 1

		cpu.currInstruction = "OPCODE"

	} else if opcode&0x6 == 0x6 {
		//LD r8, u8

		cpu.r8Write[opcode>>3](cpu.d8())
		cpu.currInstruction = "LD r8, u8"

	} else if opcode&5 == 0x5 {
		//DEC r8

		cpu.DECR8(opcode>>3)
		cpu.currInstruction = "DEC r8"

	} else if opcode&0x4 == 0x04 {
		//INC r8

		cpu.INCR8(opcode>>3)
		cpu.currInstruction = "INC r8"
	} else if opcode&0xB == 0x0B {
		//DEC r16

		*cpu.r16group1[opcode>>4]--
		cpu.currInstruction = "DEC r16"

	} else if opcode&0x3 == 0x03 {
		//INC r16

		*cpu.r16group1[opcode>>4]++
		cpu.currInstruction = "INC r16"

	} else if opcode&0x0A == 0x0A {
		//LD A, (r16) GROUP 1

		cpu.setA(cpu.gb.mmu.readbyte(*cpu.r16group1[opcode>>4]))
		cpu.currInstruction = "LD A, (r16)"

	} else if opcode&0x02 == 0x02 {
		//LD (r16) GROUP 1, A
		cpu.gb.mmu.writebyte(*cpu.r16group1[opcode>>4], cpu.A())
		cpu.currInstruction = "LD (r16), A"

	} else if opcode&0x09 == 0x09 {
		//ADD HL, r16 GROUP 1

		cpu.HL += *cpu.r16group1[opcode>>4]
		cpu.setFlag("N", false)
		cpu.setFlag("H", ((cpu.HL&0xFFF)+(*cpu.r16group1[opcode>>4]&0xFFF))&0x1000 == 0x1000)
		cpu.setFlag("C", uint32(cpu.HL)+uint32(*cpu.r16group1[opcode>>4]) > 0xFFFF)
		cpu.currInstruction = "ADD HL, r16"

	} else if opcode&0x01 == 0x01 {
		//LD r16 GROUP 1, u16

		*cpu.r16group1[opcode>>4] = cpu.d16()
		cpu.currInstruction = "LD r16,u16"

	//Third Group-------------------
	}
}
