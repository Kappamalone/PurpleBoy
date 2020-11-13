package main

//import "fmt"
//"strings"

func opcodeFormat(patternArray [8]uint8, opcode uint8) bool {
	//Takes an input in the form of a string such as
	//"11220011" and return true if the opcode matches
	//the pattern (2 are ignored bits)

	match := true
	for i := 0; i < 8; i++ {
		if patternArray[i] != 2 {
			if patternArray[i] == 1 {
				if (opcode & (1 << (7 - i))) == 0 { //Checks if (7-ith) bit is not set
					match = false
				}
			} else if patternArray[i] == 0 {
				if (opcode & (1 << (7 - i))) > 0 { //Checks if (7-ith) bit is not set
					match = false
				}
			}
		}
	}

	return match
}

func (cpu *gameboyCPU) setFlag(flag string, operand bool) {
	switch flag {
	case "Z":
		if operand {
			cpu.AF |= 128
		} else {
			cpu.AF &^= 128
		}
	case "N":
		if operand {
			cpu.AF |= 64
		} else {
			cpu.AF &^= 64
		}
	case "H":
		if operand {
			cpu.AF |= 32
		} else {
			cpu.AF &^= 32
		}
	case "C":
		if operand {
			cpu.AF |= 16
		} else {
			cpu.AF &^= 16
		}
	}
}

func (cpu *gameboyCPU) getFlag(flag string) bool {
	flagSet := false
	switch flag {
	case "Z":
		if (cpu.AF>>7)&1 == 1 {
			flagSet = true
		}
	case "N":
		if (cpu.AF>>6)&1 == 1 {
			flagSet = true
		}
	case "H":
		if (cpu.AF>>5)&1 == 1 {
			flagSet = true
		}
	case "C":
		if (cpu.AF>>4)&1 == 1 {
			flagSet = true
		}
	}
	return flagSet
}

func (cpu *gameboyCPU) getAcc() uint8 {
	return uint8((cpu.AF & 0xFF00) >> 8)
}

func (cpu *gameboyCPU) setAcc(value uint8) {
	cpu.AF = uint16(value)<<8 | (cpu.AF & 0x00F0)
}

func (cpu *gameboyCPU) carry() uint8 {
	if cpu.getFlag("C") {
		return uint8(1)
	} else {
		return uint8(0)
	}
}

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

func addSigned(opcode uint16, signedValue uint8) uint16 {
	//Th 2s Complement representation is a method of storing
	//Negative numbers in a byte. The MSB indicates if the bit is
	//negative, with the 0x80 being -128 and 0x7F being 127
	//The reason I'm not directly computing the twos complement
	//Is because these additions are adding uints of different sizes
	if signedValue>>7 == 1 {
		subtract := (1 << 7) - (signedValue & 0x7F)
		return opcode - uint16(subtract)
	} else {
		add := signedValue & 0x7F
		return opcode + uint16(add)
	}
}

//OPCODES
//Wonderful explanation for half carry flags at https://robdor.com/2016/08/10/gameboy-emulator-half-carry-flag/

func (cpu *gameboyCPU) INC(opcode uint8) {
	cpu.setFlag("H", (((cpu.r8Read[opcode]()&0x0F)+(1))&0x10 == 0x10))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + 1)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) DEC(opcode uint8) {
	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(1)) > 0xF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - 1)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)
}

func (cpu *gameboyCPU) ADD(opcode uint8, operand uint8) {
	cpu.setFlag("H", (((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F))&0x10 == 0x10))
	cpu.setFlag("C", (uint16(cpu.r8Read[opcode]())+uint16(operand)) > 0xFF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) ADC(opcode uint8, operand uint8) {
	carry := cpu.carry()
	cpu.setFlag("H", (((cpu.r8Read[opcode]()&0x0F)+(operand&0x0F)+carry)&0x10 == 0x10))
	cpu.setFlag("C", (uint16(cpu.r8Read[opcode]())+uint16(operand)+uint16(carry)) > 0xFF)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() + operand + carry)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
}

func (cpu *gameboyCPU) SUB(opcode uint8, operand uint8) {
	cpu.setFlag("H", (cpu.r8Read[opcode]()&0x0F-(operand&0x0F)) > 0xF)
	cpu.setFlag("C", (uint16(cpu.r8Read[opcode]())-uint16(operand) > 0xFF))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)
}

func (cpu *gameboyCPU) SUBC(opcode uint8, operand uint8) {
	carry := cpu.carry()
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F-(operand&0x0F))-carry) > 0xF)
	cpu.setFlag("C", ((uint16(cpu.r8Read[opcode]()) - uint16(operand) - uint16(carry)) > 0xFF))
	cpu.r8Write[opcode](cpu.r8Read[opcode]() - operand - carry)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", true)
}

func (cpu *gameboyCPU) AND(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() & operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", true)
	cpu.setFlag("C", false)
}

func (cpu *gameboyCPU) XOR(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() ^ operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", false)

}

func (cpu *gameboyCPU) OR(opcode uint8, operand uint8) {
	cpu.r8Write[opcode](cpu.r8Read[opcode]() | operand)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", false)

}

func (cpu *gameboyCPU) CP(opcode uint8, operand uint8) {
	cpu.setFlag("Z", (cpu.r8Read[opcode]()-operand) == 0)
	cpu.setFlag("N", true)
	cpu.setFlag("H", ((cpu.r8Read[opcode]()&0x0F)-(operand&0x0F)) > 0xF)
	cpu.setFlag("C", (uint16(cpu.r8Read[opcode]())-uint16(operand) > 0xFF))
}

//ALU OPCODES--------------------------
func (cpu *gameboyCPU) RLCA() {
	//Rotate register A left
	//Rotations shift bits by one place and wrap them
	//Around the byte if necessary
	cpu.setFlag("Z", false)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", cpu.getAcc()>>7 == 1)
	cpu.setAcc(cpu.getAcc()<<1 | cpu.getAcc()>>7)

}

func (cpu *gameboyCPU) RRCA() {
	//Rotate register A right
	cpu.setFlag("Z", false)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", cpu.getAcc()&0x01 == 1)
	cpu.setAcc(cpu.getAcc()>>1 | cpu.getAcc()<<7)

}

func (cpu *gameboyCPU) RLA() {
	//Rotate register A left through carry
	var cflag uint8 //Get int version from bool
	if cpu.getFlag("C") {
		cflag = 1
	}
	cpu.setFlag("Z", false)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", cpu.getAcc()>>7 == 1)
	cpu.setAcc(cpu.getAcc()<<1 | cflag)
}

func (cpu *gameboyCPU) RRA() {
	//Rotate register A right through carry
	var cflag uint8 //Get int version from
	if cpu.getFlag("C") {
		cflag = 1
	}
	cpu.setFlag("Z", false)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", cpu.getAcc()&0x01 == 1)
	cpu.setAcc(cpu.getAcc()>>1 | (cflag << 7))

}

func (cpu *gameboyCPU) DAA() {
	//Decimal Adjust Accumulator: To get correct BCD
	//Representation after an arithmetic instruction
	//Basically the scariest thing about the CPU
	//I'll be completely honest I gave up and copied another implementation

	var acc uint8 = cpu.getAcc() //Temp var ensure there is proper overflow when dealing with BCD

	if !cpu.getFlag("N") {
		//Previous instruction was an addition
		//Deal with high nibble
		if cpu.getFlag("C") || (acc > 0x99) { //0x99 instead of 0x90
			acc += 0x60
			cpu.setFlag("C", true)
		}
		//Deal with low nibble
		if cpu.getFlag("H") || ((acc & 0x0F) > 0x09) {
			acc += 0x06
		}
	} else {
		//Previous instruction was a subtraction
		if cpu.getFlag("C") {
			acc -= 0x60
			cpu.setFlag("C", true)
		}
		if cpu.getFlag("H") {
			acc -= 0x06
		}

	}

	cpu.setAcc(acc)
	cpu.setFlag("Z", cpu.getAcc() == 0)
	cpu.setFlag("H", false)
}

func (cpu *gameboyCPU) CPL() {
	//Complement accumulator
	cpu.setAcc(^cpu.getAcc())

	cpu.setFlag("N", true)
	cpu.setFlag("H", true)
}

func (cpu *gameboyCPU) CCF() {
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", !cpu.getFlag("C"))
}

func (cpu *gameboyCPU) SCF() {
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", true)
}

//JUMP OPCODES--------------------------
func (cpu *gameboyCPU) CALL(addr uint16) {
	//A call is essentially jumping to an address and placing the PC on
	//The stack to return to once the subroutine is finished executing
	cpu.SP -= 2
	cpu.gb.mmu.writeword(cpu.SP, cpu.PC)
	cpu.PC = addr
}

func (cpu *gameboyCPU) RET() {
	cpu.PC = cpu.gb.mmu.readWord(cpu.SP)
	cpu.SP += 2
}

func (cpu *gameboyCPU) JP(addr uint16) {
	cpu.PC = addr
}

func (cpu *gameboyCPU) JR(relativeJumpValue uint8) {
	//A relative jump using a signed int
	cpu.PC = addSigned(cpu.PC, relativeJumpValue)
}

func (cpu *gameboyCPU) RST(addr uint8) {
	cpu.SP -= 2
	cpu.gb.mmu.writeword(cpu.SP, cpu.PC)
	cpu.PC = uint16(addr)
}

//EXTENDED--------------------------------
func (cpu *gameboyCPU) RLC(opcode uint8) {
	cpu.setFlag("C", cpu.r8Read[opcode]()>>7 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()<<1 | cpu.r8Read[opcode]()>>7)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)

	cpu.currInstruction = "RLC"
}

func (cpu *gameboyCPU) RRC(opcode uint8) {
	cpu.setFlag("C", cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | cpu.r8Read[opcode]()<<7)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)

	cpu.currInstruction = "RRC"
}

func (cpu *gameboyCPU) RL(opcode uint8) {
	var cflag uint8
	if cpu.getFlag("C") {
		cflag = 1
	}
	cpu.setFlag("C", cpu.r8Read[opcode]()>>7 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()<<1 | cflag)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)

	cpu.currInstruction = "RL"
}

func (cpu *gameboyCPU) RR(opcode uint8) {
	var cflag uint8
	if cpu.getFlag("C") {
		cflag = 1
	}
	cpu.setFlag("C", cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | (cflag << 7))
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)

	cpu.currInstruction = "RR"
}

func (cpu *gameboyCPU) SLA(opcode uint8) {
	//So it seems arithmetic and logical shifts are
	//actually different
	cpu.setFlag("C", (cpu.r8Read[opcode]()&0x80) == 0x80)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() << 1 & 0xFE)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.currInstruction = "SLA"
}

func (cpu *gameboyCPU) SRA(opcode uint8) {
	cpu.setFlag("C", (cpu.r8Read[opcode]()&0x01) == 0x01)
	signedBit := cpu.r8Read[opcode]() & 0x80
	cpu.r8Write[opcode](cpu.r8Read[opcode]()>>1 | signedBit)

	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.currInstruction = "SRA"

}

func (cpu *gameboyCPU) SWAP(opcode uint8) {
	hi := (cpu.r8Read[opcode]() & 0xF0) >> 4
	low := cpu.r8Read[opcode]() & 0x0F
	cpu.r8Write[opcode](low<<4 | hi)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)
	cpu.setFlag("C", false)
	cpu.currInstruction = "SWAP"
}

func (cpu *gameboyCPU) SRL(opcode uint8) {
	cpu.setFlag("C", cpu.r8Read[opcode]()&0x01 == 1)
	cpu.r8Write[opcode](cpu.r8Read[opcode]() >> 1)
	cpu.setFlag("Z", cpu.r8Read[opcode]() == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", false)

	cpu.currInstruction = "SRL"
}

func (cpu *gameboyCPU) BIT(opcode uint8, place uint8) {
	//Set zflag is bit not set
	cpu.setFlag("Z", (cpu.r8Read[opcode]()&(1<<place)) == 0)
	cpu.setFlag("N", false)
	cpu.setFlag("H", true)
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
