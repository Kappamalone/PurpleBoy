package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func handleInput() bool {
	endProgram := false
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			endProgram = true
			println("Quit")
			break
		}
	}
	return endProgram
}
