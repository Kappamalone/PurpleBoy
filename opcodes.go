package main

//Addressing modes
func (cpu *gameboyCPU) d8() uint8 {
	immediateData := cpu.gb.mmu.readbyte(cpu.PC)
	cpu.PC++
	return immediateData
}

func (cpu *gameboyCPU) d16() uint16 {
	hi := uint16(cpu.gb.mmu.readbyte(cpu.PC + 1))
	low := uint16(cpu.gb.mmu.readbyte(cpu.PC))
	cpu.PC += 2

	return hi<<8 | low
}

//OPCODES

//8 BIT ARITHMETIC
func (cpu *gameboyCPU) INC(opcode uint8) {
	cpu.setH((((cpu.r8Read[opcode]()&0x0F)+(1))&0x10 == 0x10))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + 1)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
}

func (cpu *gameboyCPU) DEC(opcode uint8) {
	cpu.setH((cpu.r8Read[opcode]()&0x0F - (1)) > 0xF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - 1)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(true)
}

func (cpu *gameboyCPU) ADD(opcode uint8, operand uint8) {
	cpu.setH((((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F))&0x10 == 0x10))
	cpu.setC((uint16(cpu.r8Read[opcode]()) + uint16(operand)) > 0xFF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
}

func (cpu *gameboyCPU) ADC(opcode uint8, operand uint8) {
	carry := boolToInt(cpu.getC())
	cpu.setH((((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F)+carry)&0x10 == 0x10))
	cpu.setC((uint16(cpu.r8Read[opcode]()) + uint16(operand) + uint16(carry)) > 0xFF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand + carry)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
}

func (cpu *gameboyCPU) SUB(opcode uint8, operand uint8) {
	cpu.setH((cpu.r8Read[opcode]()&0x0F - (operand & 0x0F)) > 0xF)
	cpu.setC((uint16(cpu.r8Read[opcode]())-uint16(operand) > 0xFF))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(true)
}

func (cpu *gameboyCPU) SUBC(opcode uint8, operand uint8) {
	carry := boolToInt(cpu.getC())
	cpu.setH(((cpu.r8Read[opcode]()&0x0F - (operand & 0x0F)) - carry) > 0xF)
	cpu.setC(((uint16(cpu.r8Read[opcode]()) - uint16(operand) - uint16(carry)) > 0xFF))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand - carry)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(true)
}

func (cpu *gameboyCPU) AND(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() & operand)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(true)
	cpu.setC(false)
}

func (cpu *gameboyCPU) XOR(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() ^ operand)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)

}

func (cpu *gameboyCPU) OR(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() | operand)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)

}

func (cpu *gameboyCPU) CP(opcode uint8, operand uint8) {
	cpu.setZ((cpu.r8Read[opcode]() - operand) == 0)
	cpu.setN(true)
	cpu.setH(((cpu.r8Read[opcode]() & 0x0F) - (operand & 0x0F)) > 0xF)
	cpu.setC((uint16(cpu.r8Read[opcode]())-uint16(operand) > 0xFF))
}

//ALU OPCODES--------------------------
func (cpu *gameboyCPU) RLCA() {
	//Rotate register A left
	//Rotations shift bits by one place and wrap them
	//Around the byte if necessary
	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(cpu.getAcc()>>7 == 1)
	cpu.setAcc(cpu.getAcc()<<1 | cpu.getAcc()>>7)

}

func (cpu *gameboyCPU) RRCA() {
	//Rotate register A right
	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(cpu.getAcc()&0x01 == 1)
	cpu.setAcc(cpu.getAcc()>>1 | cpu.getAcc()<<7)

}

func (cpu *gameboyCPU) RLA() {
	//Rotate register A left through carry
	carry := boolToInt(cpu.getC())
	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(cpu.getAcc()>>7 == 1)
	cpu.setAcc(cpu.getAcc()<<1 | carry)
}

func (cpu *gameboyCPU) RRA() {
	//Rotate register A right through carry
	carry := boolToInt(cpu.getC())
	cpu.setZ(false)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(cpu.getAcc()&0x01 == 1)
	cpu.setAcc(cpu.getAcc()>>1 | (carry << 7))

}

func (cpu *gameboyCPU) DAA() {
	//Decimal Adjust Accumulator: To get correct BCD
	//Representation after an arithmetic instruction
	var acc uint8 = cpu.getAcc() //Temp var ensure there is proper overflow when dealing with BCD
	if !cpu.getN() {
		//Previous instruction was an addition
		//Deal with high nibble
		if cpu.getC() || (acc > 0x99) { //0x99 instead of 0x90
			acc += 0x60
			cpu.setC(true)
		}
		//Deal with low nibble
		if cpu.getH() || ((acc & 0x0F) > 0x09) {
			acc += 0x06
		}
	} else {
		//Previous instruction was a subtraction
		if cpu.getC() {
			acc -= 0x60
			cpu.setC(true)
		}
		if cpu.getH() {
			acc -= 0x06
		}
	}
	cpu.setAcc(acc)
	cpu.setZ(cpu.getAcc() == 0)
	cpu.setH(false)
}

func (cpu *gameboyCPU) CPL() {
	//Complement accumulator
	cpu.setAcc(^cpu.getAcc())
	cpu.setN(true)
	cpu.setH(true)
}

func (cpu *gameboyCPU) CCF() {
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(!cpu.getC())
}

func (cpu *gameboyCPU) SCF() {
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(true)
}

//PROGRAM FLOW OPCODES--------------------------
func (cpu *gameboyCPU) PUSH(data uint16) {
	//Stack progresses downwards
	cpu.SP -= 2
	cpu.gb.mmu.writeword(cpu.SP, data)
}

func (cpu *gameboyCPU) POP() uint16 {
	//Stack regresses upwards
	data := cpu.gb.mmu.readWord(cpu.SP)
	cpu.SP += 2
	return data
}

func (cpu *gameboyCPU) CALL(addr uint16) {
	//A call is essentially jumping to an address and placing the PC on
	//The stack to return to once the subroutine is finished executing
	cpu.PUSH(cpu.PC)
	cpu.PC = addr
}

func (cpu *gameboyCPU) RET() {
	cpu.PC = cpu.POP()
}

func (cpu *gameboyCPU) RST(addr uint8) {
	cpu.PUSH(cpu.PC)
	cpu.PC = uint16(addr)
}

func (cpu *gameboyCPU) JP(addr uint16) {
	cpu.PC = addr
}

func (cpu *gameboyCPU) JR(relativeJumpValue uint8) {
	//A relative jump using a signed int
	cpu.PC = addSigned(cpu.PC, relativeJumpValue)
}

//EXTENDED CB OPCODES--------------------------------
func (cpu *gameboyCPU) RLC(opcode uint8) {
	cpu.setC(cpu.r8Read[opcode]()>>7 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()<<1 | cpu.r8Read[opcode]()>>7)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) RRC(opcode uint8) {
	cpu.setC(cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | cpu.r8Read[opcode]()<<7)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) RL(opcode uint8) {
	carry := boolToInt(cpu.getC())
	cpu.setC(cpu.r8Read[opcode]()>>7 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()<<1 | carry)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) RR(opcode uint8) {
	carry := boolToInt(cpu.getC())
	cpu.setC(cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | (carry << 7))
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) SLA(opcode uint8) {
	//So it seems arithmetic and logical shifts are
	//actually different
	cpu.setC((cpu.r8Read[opcode]() & 0x80) == 0x80)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() << 1 & 0xFE)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) SRA(opcode uint8) {
	cpu.setC((cpu.r8Read[opcode]() & 0x01) == 0x01)
	signedBit := cpu.r8Read[opcode]() & 0x80
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | signedBit)

	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) SWAP(opcode uint8) {
	hi := (cpu.r8Read[opcode]() & 0xF0) >> 4
	low := cpu.r8Read[opcode]() & 0x0F
	cpu.r8Write[opcode](low<<4 | hi)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
	cpu.setC(false)
}

func (cpu *gameboyCPU) SRL(opcode uint8) {
	cpu.setC(cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() >> 1)
	cpu.setZ(cpu.r8Read[opcode]() == 0)
	cpu.setN(false)
	cpu.setH(false)
}

func (cpu *gameboyCPU) BIT(opcode uint8, place uint8) {
	//Set zflag is bit not set starting from rhs
	cpu.setZ((!bitSet(cpu.r8Read[opcode](), place)))
	cpu.setN(false)
	cpu.setH(true)
}

func (cpu *gameboyCPU) RES(opcode uint8, place uint8) {
	//Reset bit
	//Create a mask that is identical to the opcode except the bit we are resetting
	//Eg: resetting 3rd bit: 0100 -> 1011 is now the mask to and with the opcode
	mask := uint8(^(1 << place))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() & mask)
}

func (cpu *gameboyCPU) SET(opcode uint8, place uint8) {
	//Set bit
	cpu.r8Write[opcode](cpu.r8Read[opcode]() | (1 << place))
}
