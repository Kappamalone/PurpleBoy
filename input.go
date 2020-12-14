package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	inputMap = map[int]uint8{ //Maps Scancodes to bit positions
		40: 3,
		81: 3,
		229: 2,
		82: 2,
		29: 1,
		80: 1,
		27: 0,
		79: 0,
	}
)

type joypad struct {
	gb *gameboy
	selectKey bool //False: buttonss True: Directional

	buttons     uint8
	directional uint8
}

func initJoypad(gb *gameboy) *joypad {
	joypad := new(joypad)
	joypad.gb = gb
	joypad.buttons = 0xFF
	joypad.directional = 0xFF

	return joypad
}

func (joypad *joypad) handleInput() bool {
	endProgram := false
	if !isDebugging {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				endProgram = true
				break
			case *sdl.KeyboardEvent:
				joypad.SDLHandleInput(e)
			}
		}
	} else {
		//SDL handles inputs differently if there are mulitple windows open
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					endProgram = true
				}
			case *sdl.KeyboardEvent:
				joypad.SDLHandleInput(e)
			}
		}
	}

	return endProgram
}
func (joypad *joypad) SDLHandleInput(e *sdl.KeyboardEvent){
	if e.Type == sdl.KEYDOWN {
		joypad.gb.cpu.requestJoypad() //TODO: Don't do this
		switch e.Keysym.Scancode {
		//Buttons
		case 40: //Start:  Enter
			clearBit(&joypad.buttons,3)
		case 229: //Select: Shift
			clearBit(&joypad.buttons,2)
		case 29:  //B: Z
			clearBit(&joypad.buttons,1)
		case 27:  //A: X
			clearBit(&joypad.buttons,0)
		//DIRECTIONAL
		case 81: //DOWN
			clearBit(&joypad.directional,3)
		case 82: //UP
			clearBit(&joypad.directional,2)
		case 80://LEFT
			clearBit(&joypad.directional,1)
		case 79://RIGHT
			clearBit(&joypad.directional,0)
		}
	} else if e.Type == sdl.KEYUP {
		switch e.Keysym.Scancode {
		//Buttons
		case 40: //Start: Ctrl
			setBit(&joypad.buttons,3)
		case 229: //Select: Shift
			setBit(&joypad.buttons,2)
		case 29:  //B: Z
			setBit(&joypad.buttons,1)
		case 27:  //A: X
			setBit(&joypad.buttons,0)
		//DIRECTIONAL
		case 81://DOWN
			setBit(&joypad.directional,3)
		case 82: //UP
			setBit(&joypad.directional,2)
		case 80://LEFT
			setBit(&joypad.directional,1)
		case 79://RIGHT
			setBit(&joypad.directional,0)
		}
	}
}

func (joypad *joypad) writeJoypad(data uint8){
	if !bitSet(data,5){
		joypad.selectKey = false
	} else if !bitSet(data,4){
		joypad.selectKey = true
	}
}

func (joypad *joypad) readJoypad() uint8 {
	joypadReadbyte := uint8(0)
	if !joypad.selectKey {
		joypadReadbyte = joypad.buttons
	} else {
		joypadReadbyte = joypad.directional
	}

	return joypadReadbyte
}
