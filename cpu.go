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
	condition map[uint8]bool        //condition map
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
		0: func() uint8 {return uint8(cpu.BC >> 8)}, 1: func() uint8 {return uint8(cpu.BC & 0xFF)},
		2: func() uint8 {return uint8(cpu.DE >> 8)}, 3: func() uint8 {return uint8(cpu.DE & 0xFF)},
		4: func() uint8 {return uint8(cpu.HL >> 8)}, 5: func() uint8 {return uint8(cpu.HL & 0xFF)},
		6: func() uint8 { return cpu.gb.mmu.readbyte(cpu.HL) },
		7: func() uint8 {return uint8(cpu.AF >> 8)},
	}

	cpu.r8Write = map[uint8]func(uint8){
		0: func (data uint8) {cpu.BC = uint16(data)<<8 | cpu.BC&0xFF},
		1: func (data uint8) {cpu.BC = cpu.BC<<8 | uint16(data)},
		2: func (data uint8) {cpu.DE = uint16(data)<<8 | cpu.DE&0xFF},
		3: func (data uint8) {cpu.DE = cpu.DE<<8 | uint16(data)},
		4: func (data uint8) {cpu.HL = uint16(data)<<8 | cpu.HL&0xFF},
		5: func (data uint8) {cpu.HL = cpu.HL<<8 | uint16(data)},
		6: func (data uint8) {cpu.gb.mmu.writebyte(cpu.HL, data) }, 
		7: func (data uint8) {cpu.AF = uint16(data)<<8 | cpu.AF&0xF0}, //Lower 4 bits always 0
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

	cpu.condition = map[uint8]bool{
		0: !cpu.getFlag("Z"),1: cpu.getFlag("Z"),
		2: !cpu.getFlag("N"),3: cpu.getFlag("N"),
	}

	return cpu
}

func (cpu *gameboyCPU) cycle() {
	if cpu.cycles == 0 {
		fetchedInstruction := cpu.gb.mmu.readbyte(cpu.PC)
		cpu.PC++

		println(cpu.decodeAndExecute(fetchedInstruction))
	}
	cpu.cycles--
}

func (cpu *gameboyCPU) decodeAndExecute(opcode uint8) bool {
	extendedInstruction := false

	if opcode == 0x00 {
		//NOP

		cpu.currInstruction = "NOP"

	} else if opcode == 0x08 {
		//LD (u16), SP

		cpu.currInstruction = "LD (u16), SP"

	} else if opcode == 0x10 {
		//STOP

		cpu.currInstruction = "STOP"

	} else if opcode == 0x18 {
		//JR (unconditional)

		cpu.currInstruction = "JR (unconditional)"

	} else if opcodeFormat([8]uint8{0,0,1,2,2,0,0,0},opcode){
		//JR (conditional)

		cpu.currInstruction = "JR (conditional)"

	} else if opcodeFormat([8]uint8{0,0,2,2,0,0,0,1},opcode){
		//LD r16, u16

		cpu.currInstruction = "LD r16, u16"

	} else if opcodeFormat([8]uint8{0,0,2,2,1,0,0,1},opcode){
		//ADD HL, r16

		cpu.currInstruction = "ADD HL, r16"

	} else if opcodeFormat([8]uint8{0,0,2,2,0,0,1,0},opcode){
		//LD (r16 group 2), A

		cpu.currInstruction = "LD (r16), A"

	} else if opcodeFormat([8]uint8{0,0,2,2,1,0,1,0},opcode){
		//LD A, (r16 group 2)

		cpu.currInstruction = "LD A, (r16)"

	} else if opcodeFormat([8]uint8{0,0,2,2,0,0,1,1},opcode){
		//INC r16

		cpu.currInstruction = "INC r16"

	}  else if opcodeFormat([8]uint8{0,0,2,2,1,0,1,1},opcode){
		//DEC r16

		cpu.currInstruction = "DEC r16"

	} else if opcodeFormat([8]uint8{0,0,2,2,2,1,0,0},opcode){
		//INC r8

		cpu.currInstruction = "INC r8"

	} else if opcodeFormat([8]uint8{0,0,2,2,2,1,0,1},opcode){
		//DEC r8

		cpu.currInstruction = "DEC r8"
	
	} else if opcodeFormat([8]uint8{0,0,2,2,2,1,1,0},opcode){
		//LD r8, u8

		cpu.currInstruction = "LD r8, u8"
		
	} else if opcodeFormat([8]uint8{0,0,2,2,2,1,1,1},opcode){
		//OPCODE TABLE 2

		cpu.currInstruction = "OPCODE TABLE 1"
		
	} else if opcode == 0x76 {
		//HALT -> Important to have this occur before LD r8,r8 as it overlaps with LD HL,HL

		cpu.currInstruction = "HALT"

	} else if opcodeFormat([8]uint8{0,1,2,2,2,2,2,2},opcode){
		//LD r8,r8

		cpu.currInstruction = "LD r8, r8"

	} else if opcodeFormat([8]uint8{1,0,2,2,2,2,2,2},opcode){
		//ALU A,r8

		cpu.currInstruction = "ALU A,r8"
		
	} else if opcodeFormat([8]uint8{1,1,0,2,2,0,0,0},opcode){
		//RET condition

		cpu.currInstruction = "RET condition"
		
	} else if opcode == 0xE0 {
		//LD (0xFF00 + u8)

		cpu.currInstruction = "LD (0xFF00 + u8)"

	} else if opcode == 0xE8 {
		//ADD SP,i8

		cpu.currInstruction = "ADD SP, i8"

	} else if opcode == 0xF0 {
		//LD A, (0xFF00 + u8)

		cpu.currInstruction = "LD A, (0xFF00 + u8)"

	} else if opcode == 0xF8 {
		//LD HL, SP + i8

		cpu.currInstruction = "LD HL, SP + i8"

	} else if opcodeFormat([8]uint8{1,1,2,2,0,0,0,1},opcode){
		//POP r16 (group 3)

		cpu.currInstruction = "POP r16"

	} else if opcodeFormat([8]uint8{1,1,2,2,0,0,0,1},opcode){
		//OPCODE TABLE 2

		cpu.currInstruction = "OPCODE TABLE 2"

	}  else if opcodeFormat([8]uint8{1,1,0,2,2,0,1,0},opcode){
		//JP condition

		cpu.currInstruction = "JP condition"

	} else if opcode == 0xE2 {
		//LD (0xFF00 + C), A

		cpu.currInstruction = "LD (0xFF00 + C), A"

	} else if opcode == 0xEA {
		//LD (u16), A

		cpu.currInstruction = "LD (u16), A"

	} else if opcode == 0xF2 {
		//LD A, (0xFF00 + C)

		cpu.currInstruction = "LD A, (0xFF00 + C)"

	} else if opcode == 0xFA {
		//LD A, (u16)

		cpu.currInstruction = "LD A, (u16)"

	} else if opcode == 0xC3 {
		//JP u16

		cpu.currInstruction = "JP u16"

	} else if opcode == 0xCB {
		//CB prefix
		extendedInstruction = true
		cpu.d8()

		if opcodeFormat([8]uint8{0,0,2,2,2,2,2,2},opcode){
			//Shifts/Rotates

			cpu.currInstruction = "0xCB: Shifts/Rotates"

		} else if opcodeFormat([8]uint8{0,1,2,2,2,2,2,2},opcode){
			//BIT bit, r8

			cpu.currInstruction = "0xCB: BIT bit, r8"

		} else if opcodeFormat([8]uint8{1,0,2,2,2,2,2,2},opcode){
			//RES bit, r8

			cpu.currInstruction = "0xCB: RES bit, r8"
		} else if opcodeFormat([8]uint8{1,1,2,2,2,2,2,2},opcode){
			//SET bit, r8

			cpu.currInstruction = "0xCB: SET bit, r8"
		}

	} else if opcode == 0xF3 {
		//Enable interupts 

		cpu.currInstruction = "Enable interrupts"

	} else if opcode == 0xFB {
		//Disable interupts

		cpu.currInstruction = "Disable Interupts"

	} else if opcodeFormat([8]uint8{1,1,0,2,2,1,0,0},opcode){
		//CALL condition

		cpu.currInstruction = "CALL condition"

	} else if opcodeFormat([8]uint8{1,1,2,2,0,1,0,1},opcode){
		//PUSH r16 group 3

		cpu.currInstruction = "PUSH r16"

	} else if opcode == 0xCD {
		//CALL u16

		cpu.currInstruction = "CALL u16"

	} else if opcodeFormat([8]uint8{1,1,2,2,2,1,1,0},opcode){
		//ALU A, u8

		cpu.currInstruction = "ALU A, u8"

	} else if opcodeFormat([8]uint8{1,1,2,2,2,1,1,1},opcode){
		//RST

		cpu.currInstruction = "RST"

	}

	return extendedInstruction
}
