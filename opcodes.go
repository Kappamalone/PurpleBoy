package main

import (
//"fmt"
)

//LD location with value
func LD(register *uint8, value uint8) {
	*register = value
}

//LD16 bit variant
func LD16(registers [2]*uint8, value uint16) {
	//[high,low]*uint8
	*(registers)[0] = uint8(value >> 8)
	*(registers)[1] = uint8(value & 0xFF)
}

//ADD register with value
func ADD(register *uint8, value uint8) {
	*register += value
}

//ADD16 bit variant
func ADD16(registers [2]*uint8, value uint16) {
	cRegister := uint16(*(registers)[0])<<8 | uint16(*(registers)[1])&0xFF
	value += cRegister
	*(registers)[0] = uint8(value >> 8)
	*(registers)[1] = uint8(value & 0xFF)
}
