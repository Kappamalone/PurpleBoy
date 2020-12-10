package main

import (
)

type gameboyCPU struct {
	gb *gameboy
	timers *timers

	AF, BC, DE, HL, SP, PC uint16

	cycles int
	IME    bool
	IE     uint8 //Interrupts enabled
	IF     uint8 //Interrupt requests
	HALT   bool

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
	cpu.gb.ppu.LCDC = 0x91
	cpu.gb.mmu.writebyte(0xFF47,0xFC) //Palette
}

func (cpu *gameboyCPU) initMaps() {
	//Implementing the tables found in the opcode decoding chart

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
		0: func() bool { return !cpu.getZ() }, 1: func() bool { return cpu.getZ() },
		2: func() bool { return !cpu.getC() }, 3: func() bool { return cpu.getC() },
	}

	cpu.opTable1 = map[uint8]func(){
		0: cpu.RLCA, 1: cpu.RRCA,
		2: cpu.RLA, 3: cpu.RRA,
		4: cpu.DAA, 5: cpu.CPL,
		6: cpu.SCF, 7: cpu.CCF,
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
}

func initCPU(gb *gameboy, skipBootrom bool) *gameboyCPU {
	cpu := new(gameboyCPU)
	cpu.gb = gb
	cpu.timers = initTimers(cpu)
	cpu.initMaps()
	if skipBootrom {
		cpu.skipBootrom()
	}
	cpu.gb.mmu.ram[0xFF00] = 0xFF //Temp Joypad MMIO stub
	return cpu
}

func (cpu *gameboyCPU) tick(cycle int) {
	//Run one tick of the gameboy's cpu
	
	if !cpu.HALT {
		if cpu.cycles == 0 {
			if isLogging {
				//cpu.gb.debug.logTrace()
				//cpu.gb.debug.logValue(cpu.PC)
			}
			
			cpu.handleInterrupts() //Should handle interrupts on an instruction-by-instruction basis, not every tick!
	
			fetchedInstruction := cpu.gb.mmu.readbyte(cpu.PC)
			cpu.PC++
	
			if fetchedInstruction == 0xCB {
				//add cycles for CB-prefix
				cpu.cycles += extendedInstructionTiming[cpu.gb.mmu.readbyte(cpu.PC)] * 4
			} else {
				//add cycles for regular instruction
				cpu.cycles += regularInstructionTiming[fetchedInstruction] * 4
			}
	
			cpu.decodeAndExecute(fetchedInstruction)
	
		}
		cpu.cycles--
	} else {
		//Interrupts are the only way to disable HALT
		cpu.handleInterrupts()

	}
}

func (cpu *gameboyCPU) decodeAndExecute(opcode uint8) {
	//TODO: make sure instruction timings for branch instructions are correct by
	//Referencing the opcode table
	if opcode == 0x00 {
		//NOP

	} else if opcode == 0x08 {
		//LD (u16), SP
		cpu.gb.mmu.writeword(cpu.d16(), cpu.SP)

	} else if opcode == 0x10 {
		//STOP <- Not really sure

	} else if opcode == 0x18 {
		//JR (unconditional)
		cpu.JR(cpu.d8())

	} else if opcodeFormat([8]uint8{0, 0, 1, 2, 2, 0, 0, 0}, opcode) {
		//JR (conditional)
		offset := cpu.d8() //Important to call cpu.d8() to increment PC
		if cpu.condition[opcode>>3&0x03]() {
			cpu.JR(offset)
			cpu.cycles += 4
		}

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 0, 1}, opcode) {
		//LD r16, u16
		*cpu.r16group1[opcode>>4&0x3] = cpu.d16()

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 0, 1}, opcode) {
		//ADD HL, r16
		cpu.setH(((cpu.HL&0x0FFF)+(*cpu.r16group1[opcode>>4&0x03]&0x0FFF))&0x1000 == 0x1000)
		cpu.setC((uint32(cpu.HL)+uint32(*cpu.r16group1[opcode>>4&0x03]) > 0xFFFF))
		cpu.HL += *cpu.r16group1[opcode>>4&0x03]
		cpu.setN(false)

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 1, 0}, opcode) {
		//LD (r16 group 2), A
		register := (opcode >> 4 & 0x3)
		cpu.gb.mmu.writebyte(*cpu.r16group2[register], cpu.getAcc())
		if register == 2 {
			cpu.HL++
		} else if register == 3 {
			cpu.HL--
		}

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 1, 0}, opcode) {
		//LD A, (r16 group 2)
		register := (opcode >> 4 & 0x3)
		cpu.setAcc(cpu.gb.mmu.readbyte(*cpu.r16group2[opcode>>4&0x3]))
		if register == 2 {
			cpu.HL++
		} else if register == 3 {
			cpu.HL--
		}

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 0, 0, 1, 1}, opcode) {
		//INC r16
		*cpu.r16group1[opcode>>4&0x3]++

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 1, 0, 1, 1}, opcode) {
		//DEC r16
		*cpu.r16group1[opcode>>4&0x3]--

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 0, 0}, opcode) {
		//INC r8
		cpu.INC(opcode >> 3 & 0x07)

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 0, 1}, opcode) {
		//DEC r8
		cpu.DEC(opcode >> 3 & 0x07)

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 1, 0}, opcode) {
		//LD r8, u8
		cpu.r8Write[opcode>>3&0x07](cpu.d8())

	} else if opcodeFormat([8]uint8{0, 0, 2, 2, 2, 1, 1, 1}, opcode) {
		//Accumulator/Flag operations : opcode table 1
		cpu.opTable1[opcode>>3&0x07]()

	} else if opcode == 0x76 {
		//HALT -> Important to have this occur before LD r8,r8 as it overlaps with LD HL,HL
		cpu.HALT = true

	} else if opcodeFormat([8]uint8{0, 1, 2, 2, 2, 2, 2, 2}, opcode) {
		//LD r8,r8
		cpu.r8Write[opcode>>3&0x07](cpu.r8Read[opcode&0x07]())

	} else if opcodeFormat([8]uint8{1, 0, 2, 2, 2, 2, 2, 2}, opcode) {
		//ALU A,r8
		cpu.opTable2[opcode>>3&0x07](7, cpu.r8Read[opcode&0x07]())

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 0, 0, 0}, opcode) {
		//RET condition
		if cpu.condition[opcode>>3&0x03]() {
			cpu.RET()
			cpu.cycles += 12
		}

	} else if opcode == 0xE0 {
		//LD (0xFF00 + u8), A
		cpu.gb.mmu.writebyte(0xFF00+uint16(cpu.d8()), cpu.getAcc())

	} else if opcode == 0xE8 {
		//ADD SP,i8
		signedValue := cpu.d8()
		cpu.setH(((cpu.SP&0x0F)+(uint16(signedValue)&0x0F))&0x10 == 0x10)
		cpu.setC((uint16(cpu.SP&0xFF) + uint16(signedValue)) > 0xFF)
		cpu.SP = addSigned(cpu.SP, signedValue)
		cpu.setZ(false)
		cpu.setN(false)

	} else if opcode == 0xF0 {
		//LD A, (0xFF00 + u8)
		cpu.setAcc(cpu.gb.mmu.readbyte(0xFF00 + uint16(cpu.d8())))

	} else if opcode == 0xF8 {
		//LD HL, SP + i8
		signedValue := cpu.d8()
		cpu.setH(((cpu.SP&0x0F)+(uint16(signedValue)&0x0F))&0x10 == 0x10)
		cpu.setC((uint16(cpu.SP&0xFF) + uint16(signedValue)) > 0xFF)
		cpu.HL = addSigned(cpu.SP, signedValue)
		cpu.setZ(false)
		cpu.setN(false)

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 0, 0, 0, 1}, opcode) {
		//POP r16 (group 3)
		register := opcode >> 4 & 0x03
		*cpu.r16group3[register] = cpu.POP()
		if register == 3 {
			cpu.AF &= 0xFFF0 //Always set lower 4 bits of AF to 0
		}

	} else if opcodeFormat([8]uint8{1, 1, 0, 0, 1, 0, 0, 1}, opcode) {
		//RET
		cpu.RET()

	} else if opcodeFormat([8]uint8{1, 1, 0, 1, 1, 0, 0, 1}, opcode) {
		//RETI POSSIBLE PROBLEM
		cpu.RET()
		cpu.IME = true

	} else if opcodeFormat([8]uint8{1, 1, 1, 0, 1, 0, 0, 1}, opcode) {
		//JP HL
		cpu.PC = cpu.HL

	} else if opcodeFormat([8]uint8{1, 1, 1, 1, 1, 0, 0, 1}, opcode) {
		//LD SP,HL
		cpu.SP = cpu.HL

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 0, 1, 0}, opcode) {
		//JP condition
		jmp := cpu.d16()
		if cpu.condition[opcode>>3&0x03]() {
			cpu.JP(jmp)
			cpu.cycles += 4
		}

	} else if opcode == 0xE2 {
		//LD (0xFF00 + C), A
		cpu.gb.mmu.writebyte(0xFF00+uint16(cpu.r8Read[1]()), cpu.getAcc())

	} else if opcode == 0xEA {
		//LD (u16), A
		addr := cpu.d16()
		cpu.gb.mmu.writebyte(addr, cpu.getAcc())

	} else if opcode == 0xF2 {
		//LD A, (0xFF00 + C)
		cpu.setAcc(cpu.gb.mmu.readbyte(0xFF00 + uint16(cpu.r8Read[1]())))

	} else if opcode == 0xFA {
		//LD A, (u16)
		cpu.setAcc(cpu.gb.mmu.readbyte(cpu.d16()))

	} else if opcode == 0xC3 {
		//JP u16
		cpu.JP(cpu.d16())

	} else if opcode == 0xCB {
		//CB prefix
		opcode = cpu.gb.mmu.readbyte(cpu.PC) //update opcode
		opcodeType := opcode >> 6
		if opcodeType == 0 {
			//Shifts/Rotates
			cpu.opTable3[(opcode>>3)&0x07](opcode & 0x7)

		} else if opcodeType == 1 {
			//BIT bit, r8
			cpu.BIT(opcode&0x7, (opcode>>3)&0x7)

		} else if opcodeType == 2 {
			//RES bit, r8
			cpu.RES(opcode&0x7, (opcode>>3)&0x7)

		} else if opcodeType == 3 {
			//SET bit, r8
			cpu.SET(opcode&0x7, (opcode>>3)&0x7)
		}

		cpu.PC++ //Adjust PC after dealing with extended opcodes

	} else if opcode == 0xF3 {
		//Disable interupts
		cpu.IME = false

	} else if opcode == 0xFB {
		//Enable interupts
		cpu.IME = true

	} else if opcodeFormat([8]uint8{1, 1, 0, 2, 2, 1, 0, 0}, opcode) {
		//CALL condition
		addr := cpu.d16()
		if cpu.condition[opcode>>3&0x03]() {
			cpu.CALL(addr)
			cpu.cycles += 12
		}

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 0, 1, 0, 1}, opcode) {
		//PUSH r16 group 3
		cpu.PUSH(*cpu.r16group3[opcode>>4&0x03])

	} else if opcode == 0xCD {
		//CALL u16
		cpu.CALL(cpu.d16())

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 2, 1, 1, 0}, opcode) {
		//ALU A, u8
		cpu.opTable2[opcode>>3&0x07](7, cpu.d8())

	} else if opcodeFormat([8]uint8{1, 1, 2, 2, 2, 1, 1, 1}, opcode) {
		//RST: Call to a given vector
		cpu.RST(opcode & 0x38)

	} else {
		cpu.gb.debug.printConsole("ILLEGAL OPCODE!\n", "cyan")
	}

}
