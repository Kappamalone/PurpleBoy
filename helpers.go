package main

import ()

func checkErr(err error, errormsg string) {
	if err != nil {
		panic(errormsg)
	}
}

//F register getters and setters
func (cpu *gameboyCPU) setZ(flag bool) {
	if flag {
		cpu.AF |= 128
	} else {
		cpu.AF &^= 128
	}
}

func (cpu *gameboyCPU) setN(flag bool) {
	if flag {
		cpu.AF |= 64
	} else {
		cpu.AF &^= 64
	}
}

func (cpu *gameboyCPU) setH(flag bool) {
	if flag {
		cpu.AF |= 32
	} else {
		cpu.AF &^= 32
	}
}

func (cpu *gameboyCPU) setC(flag bool) {
	if flag {
		cpu.AF |= 16
	} else {
		cpu.AF &^= 16
	}
}

func (cpu *gameboyCPU) getZ() bool {
	return (cpu.AF>>7)&1 == 1
}

func (cpu *gameboyCPU) getN() bool {
	return (cpu.AF>>6)&1 == 1
}

func (cpu *gameboyCPU) getH() bool {
	return (cpu.AF>>5)&1 == 1
}

func (cpu *gameboyCPU) getC() bool {
	return (cpu.AF>>4)&1 == 1
}

func (cpu *gameboyCPU) getAcc() uint8 {
	return uint8((cpu.AF & 0xFF00) >> 8)
}

func (cpu *gameboyCPU) setAcc(value uint8) {
	cpu.AF = uint16(value)<<8 | (cpu.AF & 0x00F0)
}

func boolToInt(flag bool) uint8 {
	if flag {
		return uint8(1)
	} else {
		return uint8(0)
	}
}

func bitSet(data interface{}, place uint8) bool {
	//Checks if bit is set starting from the rhs
	//I really wish I wrote this function earlier...
	var isSet bool
	switch d := data.(type) {
	case uint8:
		if (d>>place)&0x01 == 0x01 {
			isSet = true
		} else {
			isSet = false
		}
	case uint16:
		if (d>>place)&0x01 == 0x01 {
			isSet = true
		} else {
			isSet = false
		}
	}
	return isSet
}

func setBit(data *uint8, place uint8){
	*data |= (1 << place)
}

func clearBit(data *uint8, place uint8){
	*data &^= (1 << place)
}


func inRange(value uint16, lowerBound uint16, upperBound uint16) bool {
	//Used by MMU to have cleaner looking code
	return value >= lowerBound && value <= upperBound
}

func isZero(value int) bool {
	return value == 0
}

func addSigned(opcode uint16, signedValue uint8) uint16 {
	//This new method I found courtesy of the emudev discord server is a lot cleaner
	//Basically we convert the uint8 => int8 so that the compiler knows its signed
	//Then we convert this to an int16, which retains this sign
	//Finally we convert it to an uint16 to add together with the requested value
	return opcode + uint16(int16(int8(signedValue)))
}

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
					break
				}
			} else if patternArray[i] == 0 {
				if (opcode & (1 << (7 - i))) > 0 { //Checks if (7-ith) bit is not set
					match = false
					break
				}
			}
		}
	}

	return match
}
