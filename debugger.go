package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type debugger struct {
	gb *gameboy

	registerDebug *widgets.Paragraph
	flagDebug     *widgets.Paragraph
}

func initDebugger(gb *gameboy) *debugger {
	debug := new(debugger)
	debug.gb = gb

	ui.Init()
	debug.initTermui()
	debug.renderDebug()

	return debug
}

func (debug *debugger) initTermui() {

	debug.registerDebug = widgets.NewParagraph()
	debug.flagDebug = widgets.NewParagraph()

	debug.registerDebug.Title = "Registers"
	debug.flagDebug.Title = "Flags"

	debug.registerDebug.BorderStyle.Fg = ui.ColorWhite
	debug.flagDebug.BorderStyle.Fg = ui.ColorWhite

	debug.registerDebug.SetRect(0, 0, 59, 40)
	debug.flagDebug.SetRect(60, 0, 120, 119)
}

func (debug *debugger) renderDebug() {
	ui.Render(debug.registerDebug, debug.flagDebug)
}
