package main

import ()

type debugger struct {
	gb *gameboy
}

func initDebugger(gb *gameboy) *debugger {
	debug := new(debugger)
	debug.gb = gb

	return debug
}
