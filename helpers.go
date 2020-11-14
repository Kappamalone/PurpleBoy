package main

import (

)

func checkErr(err error, errormsg string) {
	if err != nil {
		panic(errormsg)
	}
}

func bitSet(data uint8, place uint8) bool {
	//Checks if bit is set starting from the rhs
	//I really wish I wrote this function earlier...
	if (data>>place)&0x01 == 0x01 {
		return true
	} else {
		return false
	}
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
