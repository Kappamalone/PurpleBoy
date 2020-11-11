package main

import "fmt"

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
	halt            bool

	r8Read    map[uint8]func() uint8       //r8 group 1 reads
	r8Write   map[uint8]func(uint8)        //r8 group 1 write
	r16group1 map[uint8]*uint16            //r16 group 1
	r16group2 map[uint8]*uint16            //r16 group 2
	r16group3 map[uint8]*uint16            //r16 group 3
	condition map[uint8]func() bool        //condition map
	opTable1  map[uint8]func()             //opcode table 1 (accumulator and flag operation)
	opTable2  map[uint8]func(uint8, uint8) //opcode table 2 (ALU operations)
	opTable3  map[uint8]func(uint8)        //opcode table 3(CB shift/rotate operations)
}

//Below are the least amount of cycles taken
//To complete a particular instruction in M-cycles
var regularInstructionTiming = [256]int{
	1, 3, 2, 2, 1, 1, 2, 1, 5, 2, 2, 2, 1, 1, 2, 1,
	0, 3, 2, 2, 1, 1, 2, 1, 3, 2, 2, 2, 1, 1, 2, 1,
	2, 3, 2, 2, 1, 1, 2, 1, 2, 2, 2, 2, 1, 1, 2, 1,
	2, 3, 2, 2, 3, 3, 3, 1, 2, 2, 2, 2, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	2, 2, 2, 2, 2, 2, 0, 2, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	2, 3, 3, 4, 3, 4, 2, 4, 2, 4, 3, 0, 3, 6, 2, 4,
	2, 3, 3, 0, 3, 4, 2, 4, 2, 4, 3, 0, 3, 0, 2, 4,
	3, 3, 2, 0, 0, 4, 2, 4, 4, 1, 4, 0, 0, 0, 2, 4,
	3, 3, 2, 1, 0, 4, 2, 4, 3, 2, 4, 1, 0, 0, 2, 4,
}

var extendedInstructionTiming = [256]int{
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	2, 2, 2, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 3, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
	2, 2, 2, 2, 2, 2, 4, 2, 2, 2, 2, 2, 2, 2, 4, 2,
}

func (cpu *gameboyCPU) skipBootrom() {
	cpu.PC = 0x0100
	cpu.SP = 0xFFFE
	cpu.AF = 0x01B0
	cpu.BC = 0x0013
	cpu.DE = 0x00D8
	cpu.HL = 0x014D
}

func initCPU(gb *gameboy) *gameboyCPU {
	cpu := new(gameboyCPU)
	cpu.gb = gb
	cpu.skipBootrom()

	cpu.gb.mmu.writebyte(0xFF44, 0x90) //TEMP!!!

	//Some premptively messy code to make other parts easier to write :)
	//Actually nevermind it's still hard
	cpu.r8Read = map[uint8]func() uint8{
		0: func() uint8 { return uint8((cpu.BC & 0xFF00) >> 8) }, 1: func() uint8 { return uint8(cpu.BC & 0xFF) },
		2: func() uint8 { return uint8((cpu.DE & 0xFF00) >> 8) }, 3: func() uint8 { return uint8(cpu.DE & 0xFF) },
		4: func() uint8 { return uint8((cpu.HL & 0xFF00) >> 8) }, 5: func() uint8 { return uint8(cpu.HL & 0xFF) },
		6: func() uint8 { return cpu.gb.mmu.readbyte(cpu.HL) },
		7: func() uint8 { return uint8((cpu.AF & 0xFF00) >> 8) },
	}

	cpu.r8Write = map[uint8]func(uint8){
		0: func(data uint8) { cpu.BC = uint16(data)<<8 | cpu.BC&0xFF },
		1: func(data uint8) { cpu.BC = (cpu.BC & 0xFF00) | uint16(data) },
		2: func(data uint8) { cpu.DE = uint16(data)<<8 | cpu.DE&0xFF },
		3: func(data uint8) { cpu.DE = (cpu.DE & 0xFF00) | uint16(data) },
		4: func(data uint8) { cpu.HL = uint16(data)<<8 | cpu.HL&0xFF },
		5: func(data uint8) { cpu.HL = (cpu.HL & 0xFF00) | uint16(data) },
		6: func(data uint8) { cpu.gb.mmu.writebyte(cpu.HL, data) },
		7: func(data uint8) { cpu.AF = uint16(data)<<8 | cpu.AF&0xF0 }, //Lower 4 bits always 0
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

	cpu.condition = map[uint8]func() bool{
		0: func() bool { return !cpu.getFlag("Z") }, 1: func() bool { return cpu.getFlag("Z") },
		2: func() bool { return !cpu.getFlag("C") }, 3: func() bool { return cpu.getFlag("C") },
	}

	cpu.opTable1 = map[uint8]func(){
		0: cpu.RLCA, 1: cpu.RRCA,
		2: cpu.RLA, 3: cpu.RRA,
		4: cpu.DAA, 5: cpu.CPL,
		6: cpu.CCF, 7: cpu.SCF,
	}

	cpu.opTable2 = map[uint8]func(uint8, uint8){
		0: cpu.ADD, 1: cpu.ADC, 2: cpu.SUB, 3: cpu.SUBC,
		4: cpu.AND, 5: cpu.XOR, 6: cpu.OR, 7: cpu.CP,
	}

	cpu.opTable3 = map[uint8]func(uint8){
		0: cpu.RLC, 1: cpu.RRC,
		2: cpu.RL, 3: cpu.RR,
		4: cpu.SLA, 5: cpu.SRA,
		6: cpu.SWAP, 7: cpu.SRL,
	}

	return cpu
}

func (cpu *gameboyCPU) cycle() {
	/*
		if cpu.cycles == 0 {
			fetchedInstruction := cpu.gb.mmu.readbyte(cpu.PC)
			cpu.PC++

			if fetchedInstruction == 0xCB {
				//add cycles for cb
				cpu.cycles += extendedInstructionTiming[cpu.gb.mmu.readbyte(cpu.PC)] * 4
			} else {
				//add cycles for regular instruction
				cpu.cycles += regularInstructionTiming[fetchedInstruction] * 4
			}

			cpu.decodeAndExecute(fetchedInstruction)

		}
		cpu.cycles-- */
	fetchedInstruction := cpu.gb.mmu.readbyte(cpu.PC)
	cpu.PC++
	cpu.decodeAndExecute(fetchedInstruction)
}

func (cpu *gameboyCPU) decodeAndExecute(opcode uint8) {
	//TODO: make sure instruction timings for branch instructions are correct by
	//Referencing the opcode table

	if opcode == 0x00 {
		//NOP
		cpu.currInstruction = "NOP"

	} else if opcode == 0x08 {
		//LD (u16), SP
		cpu.gb.mmu.writeword(cpu.d16(), cpu.SP)
		cpu.currInstruction = "LD (u16), SP"

	} else if opcode == 0x10 {
		//STOP <- Not really sure
		cpu.currInstruction = "STOP"

	} else if opcode == 0x18 {
		//JR (unconditional)
		cpu.JR(cpu.d8())
		cpu.currInstruction = "JR (unconditional)"

	} else if opcodeFormat([8]uint8{0, 0, 1, 2, 2, 0, 0, 0}, opcode) {
		//JR (conditional)
		offset := cpu.d8() //Important to call cpu.d8() to increment PC
		if cpu.condition[opcode>>3&0x03]() {
			cpu.JR(offset)
			cpu.cycles += 4
		}
		cpu.currInstruction = "JR (conditional)"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 0, 1}, opcode) {
		//LD r16, u16

		*cpu.r16group1[opcode>>4&0x3] = cpu.d16()
		cpu.currInstruction = "LD r16, u16"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 0, 1}, opcode) {
		//ADD HL, r16
		cpu.setFlag("H", ((cpu.HL&0x0F00)+(*cpu.r16group1[opcode>>4&0x03])&0x0F00)&0x1000 == 0x1000)
		cpu.setFlag("C", (uint32(cpu.HL)+uint32(*cpu.r16group1[opcode>>4&0x03]) > 0xFFFF))
		cpu.HL += *cpu.r16group1[opcode>>4&0x03]
		cpu.setFlag("N", false)
		cpu.currInstruction = "ADD HL, r16"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 1, 0}, opcode) {
		//LD (r16 group 2), A

		cpu.gb.mmu.writebyte(*cpu.r16group2[opcode>>4&0x3], cpu.r8Read[7]())

		if opcode>>4&0x3 == 2 {
			cpu.HL++
		} else if opcode>>4&0x3 == 3 {
			cpu.HL--
		}
		cpu.currInstruction = "LD (r16), A"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 1, 0}, opcode) {
		//LD A, (r16 group 2)

		cpu.r8Write[7](cpu.gb.mmu.readbyte(*cpu.r16group2[opcode>>4&0x3]))

		if opcode>>4&0x3 == 2 {
			cpu.HL++
		} else if opcode>>4&0x3 == 3 {
			cpu.HL--
		}
		cpu.currInstruction = "LD A, (r16)"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 1, 1}, opcode) {
		//INC r16

		*cpu.r16group1[opcode>>4&0x3]++
		cpu.currInstruction = "INC r16"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 1, 1}, opcode) {
		//DEC r16

		*cpu.r16group1[opcode>>4&0x3]--
		cpu.currInstruction = "DEC r16"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 0, 0}, opcode) {
		//INC r8

		cpu.INC(opcode >> 3 & 0x07)
		cpu.currInstruction = "INC r8"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 0, 1}, opcode) {
		//DEC r8

		cpu.DEC(opcode >> 3 & 0x07)
		cpu.currInstruction = "DEC r8"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 1, 0}, opcode) {
		//LD r8, u8

		cpu.r8Write[opcode>>3&0x07](cpu.d8())
		cpu.currInstruction = "LD r8, u8"

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 1, 1}, opcode) {
		//Accumulator/Flag operations : opcode table 1
		cpu.opTable1[opcode>>3&0x07]()

	} else if opcode == 0x76 {
		//HALT -> Important to have this occur before LD r8,r8 as it overlaps with LD HL,HL
		cpu.halt = true
		cpu.currInstruction = "HALT"

	} else if opcodeFormat([8]uint8{0, 1, 2, 2, 2, 2, 2, 2}, opcode) {
		//LD r8,r8
		cpu.r8Write[opcode>>3&0x07](cpu.r8Read[opcode&0x07]())
		cpu.currInstruction = "LD r8, r8"

	} else if opcodeFormat([8]uint8{1, 0, 2, 2, 2, 2, 2, 2}, opcode) {
		//ALU A,r8
		cpu.opTable2[opcode>>3&0x07](7, cpu.r8Read[opcode&0x07]())

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 0, 0, 0}, opcode) {
		//RET condition
		if cpu.condition[opcode>>3&0x03]() {
			cpu.RET()
			cpu.cycles += 12
		}
		cpu.currInstruction = "RET condition"

	} else if opcode == 0xE0 {
		//LD (0xFF00 + u8), A
		cpu.gb.mmu.writebyte(0xFF00+uint16(cpu.d8()), cpu.r8Read[7]())
		cpu.currInstruction = "LD (0xFF00 + u8)"

	} else if opcode == 0xE8 {
		//ADD SP,i8
		cpu.SP = addSigned(cpu.SP, cpu.d8())
		cpu.currInstruction = "ADD SP, i8"

	} else if opcode == 0xF0 {
		//LD A, (0xFF00 + u8)
		cpu.r8Write[7](cpu.gb.mmu.readbyte(0xFF00 + uint16(cpu.d8())))
		cpu.currInstruction = "LD A, (0xFF00 + u8)"

	} else if opcode == 0xF8 {
		//LD HL, SP + i8
		cpu.HL = addSigned(cpu.SP, cpu.d8())
		cpu.currInstruction = "LD HL, SP + i8"

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 0, 0, 0, 1}, opcode) {
		//POP r16 (group 3)
		*cpu.r16group3[opcode>>4&0x03] = cpu.gb.mmu.readWord(cpu.SP)
		if opcode>>4&0x03 == 3 {
			cpu.AF &= 0xFFF0 //Always set lower 4 bits of AF to 0
		}
		cpu.SP += 2 //Stack regresses upwards
		cpu.currInstruction = "POP r16"

	} else if opcodeFormat([8]uint8{1, 1, 0, 0, 1, 0, 0, 1}, opcode) {
		//RET
		cpu.RET()
		cpu.currInstruction = "RET"

	} else if opcodeFormat([8]uint8{1, 1, 0, 1, 1, 0, 0, 1}, opcode) {
		//RETI
		cpu.RET()
		cpu.gb.mmu.writebyte(0xFFFF, 1)
		cpu.currInstruction = "RETI"

	} else if opcodeFormat([8]uint8{1, 1, 1, 0, 1, 0, 0, 1}, opcode) {
		//JP HL
		cpu.PC = cpu.HL
		cpu.currInstruction = "JP HL"

	} else if opcodeFormat([8]uint8{1, 1, 1, 1, 1, 0, 0, 1}, opcode) {
		//LD SP,HL
		cpu.SP = cpu.HL
		cpu.currInstruction = "LD SP, HL"

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 0, 1, 0}, opcode) {
		//JP condition
		jmp := cpu.d16()
		if cpu.condition[opcode>>3&0x03]() {
			cpu.JP(jmp)
			cpu.cycles += 4
		}
		cpu.currInstruction = "JP condition"

	} else if opcode == 0xE2 {
		//LD (0xFF00 + C), A
		cpu.gb.mmu.writebyte(0xFF00+uint16(cpu.r8Read[1]()), cpu.r8Read[7]())
		cpu.currInstruction = "LD (0xFF00 + C), A"

	} else if opcode == 0xEA {
		//LD (u16), A
		cpu.gb.mmu.writebyte(cpu.d16(), cpu.r8Read[7]())
		cpu.currInstruction = "LD (u16), A"

	} else if opcode == 0xF2 {
		//LD A, (0xFF00 + C)
		cpu.r8Write[7](cpu.gb.mmu.readbyte(0xFF00 + uint16(cpu.r8Read[1]())))
		cpu.currInstruction = "LD A, (0xFF00 + C)"

	} else if opcode == 0xFA {
		//LD A, (u16)
		cpu.r8Write[7](cpu.gb.mmu.readbyte(cpu.d16()))
		cpu.currInstruction = "LD A, (u16)"

	} else if opcode == 0xC3 {
		//JP u16
		cpu.JP(cpu.d16())
		cpu.currInstruction = fmt.Sprintf("JP Conditional")

	} else if opcode == 0xCB {
		//CB prefix
		opcode = cpu.gb.mmu.readbyte(cpu.PC) //update opcode

		if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 2, 2, 2}, opcode) {
			//Shifts/Rotates
			cpu.opTable3[opcode>>3&0x07](opcode & 0x07)

		} else if opcodeFormat([8]uint8{0, 1, 2, 2, 2, 2, 2, 2}, opcode) {
			//BIT bit, r8
			cpu.BIT(opcode&0x7, opcode>>3&0x7)
			cpu.currInstruction = "0xCB: BIT bit, r8"

		} else if opcodeFormat([8]uint8{1, 0, 2, 2, 2, 2, 2, 2}, opcode) {
			//RES bit, r8
			cpu.RES(opcode&0x7, opcode>>3&0x7)
			cpu.currInstruction = "0xCB: RES bit, r8"
		} else if opcodeFormat([8]uint8{1, 1, 2, 2, 2, 2, 2, 2}, opcode) {
			//SET bit, r8
			cpu.SET(opcode&0x7, opcode>>3&0x7)
			cpu.currInstruction = "0xCB: SET bit, r8"
		}

		cpu.PC++ //Adjust PC after dealing with extended opcodes

	} else if opcode == 0xF3 {
		//Enable interupts
		cpu.gb.mmu.writebyte(0xFFFF, 1)
		cpu.currInstruction = "Enable interrupts"

	} else if opcode == 0xFB {
		//Disable interupts
		cpu.gb.mmu.writebyte(0xFFFF, 0)
		cpu.currInstruction = "Disable Interupts"

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 1, 0, 0}, opcode) {
		//CALL condition
		addr := cpu.d16()
		if cpu.condition[opcode>>3&0x03]() {
			cpu.CALL(addr)
			cpu.cycles += 12
		}
		cpu.currInstruction = "CALL condition"

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 0, 1, 0, 1}, opcode) {
		//PUSH r16 group 3 POSSIBLE PROBLEM
		cpu.SP -= 2 //Stack grows downwards
		cpu.gb.mmu.writeword(cpu.SP, *cpu.r16group3[opcode>>4&0x03])
		cpu.currInstruction = "PUSH r16"

	} else if opcode == 0xCD {
		//CALL u16
		cpu.CALL(cpu.d16())
		cpu.currInstruction = "CALL u16"

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 2, 1, 1, 0}, opcode) {
		//ALU A, u8
		cpu.opTable2[opcode>>3&0x07](7, cpu.d8())
		cpu.currInstruction = "ALU A, u8"

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 2, 1, 1, 1}, opcode) {
		//RST: Call to a given vector
		//Not really sure how it works :0
		cpu.currInstruction = "RST"

	} else {
		println("ILLEGAL OPCODE ", opcode)
	}
}
