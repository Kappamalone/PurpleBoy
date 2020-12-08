package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func handleInput() bool {
	endProgram := false
	if !isDebugging {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				endProgram = true
				break
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
			}
		}
	}

	return endProgram
}
