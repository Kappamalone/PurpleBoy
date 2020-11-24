package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	//Screen specs
	screenWidth  = 160
	screenHeight = 144
	windowScale  = 5

	tilewindowWidth  = 128
	tilewindowHeight = 192
	tilewindowScale  = 4

	//Gameboy colours
	dark   = 0x2c2137
	ldark  = 0x764462
	lwhite = 0xedb4a1
	white  = 0xa96868
)

var (
	//colours = [4]uint32{dark,white,dark,dark}
	colours = [4]uint32{white, lwhite, ldark, dark}
	//colours = [4]uint32{dark,ldark,lwhite,white}
)

func setRendercolour(renderer *sdl.Renderer, colour uint32, alpha uint8) {
	renderer.SetDrawColor(uint8((colour&0xFF0000)>>16), uint8((colour&0x00FF00)>>8), uint8((colour & 0x0000FF)), alpha)
}

//PPU is the pixel processing unit of the system
//It's a custom GPU that utilises tile based rendering
type PPU struct {
	gb   *gameboy
	VRAM [8 * 1024]uint8

	window      *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	frameBuffer []uint8

	tileWindow      *sdl.Window
	tileRenderer    *sdl.Renderer
	tileTexture     *sdl.Texture
	tileframeBuffer []uint8

	dotClock int //Used to determine what the PPU should be doing
}

func initPPU(gb *gameboy) *PPU {
	ppu := new(PPU)
	ppu.gb = gb
	ppu.window, ppu.renderer = initSDL()
	if isDebugging {
		ppu.tileWindow, ppu.tileRenderer = initSDLDebugging()
		ppu.tileframeBuffer = make([]uint8, tilewindowWidth*tilewindowHeight*4) //RGBA32 uses 4 bytes per pixel
		ppu.tileTexture, _ = ppu.tileRenderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, tilewindowWidth, tilewindowHeight)
		ppu.tileRenderer.SetScale(tilewindowScale, tilewindowScale)
	}

	ppu.renderer.SetScale(windowScale, windowScale)
	ppu.frameBuffer = make([]uint8, screenWidth*screenHeight*4) //RGBA32 uses 4 bytes per pixel
	ppu.texture, _ = ppu.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, screenWidth, screenHeight)
	ppu.clearScreen()

	return ppu
}

func initSDL() (*sdl.Window, *sdl.Renderer) {
	//Does the necessary setup for the SDL library

	//Initialise SDL
	err := sdl.Init(sdl.INIT_VIDEO)
	checkErr(err, "SDL initialisation error")

	//Create window
	window, err := sdl.CreateWindow("Purpleboy!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth*windowScale, screenHeight*windowScale, sdl.WINDOW_SHOWN)
	checkErr(err, "Window creation error")
	window.SetResizable(true)

	//Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "renderer creation error")

	return window, renderer

}

func initSDLDebugging() (*sdl.Window, *sdl.Renderer) {
	//Initialises the required windows for debugging purposes

	tileWindow, err := sdl.CreateWindow("Debug", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, tilewindowWidth*tilewindowScale, tilewindowHeight*tilewindowScale, sdl.WINDOW_SHOWN)
	checkErr(err, "Debug window creation error")
	tileWindow.SetResizable(true)

	tileRenderer, err := sdl.CreateRenderer(tileWindow, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "Debug renderer creation error")

	return tileWindow, tileRenderer

}

func (ppu *PPU) tick() {
	//PPU runs at 2.1MHZ, so this tick function is called every 2 T-cycles
	/*

		LCDC := ppu.gb.mmu.readbyte(0xFF40)
		//LCDSTAT := ppu.gb.mmu.readbyte(0xFF41)

		displayEnable := bitSet(LCDC, 7)
		windowTileMap := 0x9800
		if bitSet(LCDC, 6) {
			windowTileMap = 0x9C00
		}
		windowDisplayEnable := bitSet(LCDC, 5)
		bgwTileData := 0x9000
		if bitSet(LCDC, 4) {
			bgwTileData = 0x8000
		}
		bgTileMap := 0x9800
		if bitSet(LCDC,3) {
			bgTileMap = 0x9C00
		}
		//spriteSize := bitset(LCDC,2)
		//spriteEnable := bitset(LCDC,1)
		//bgwDisplayPriority := bitset(LCDC,0)
	*/

	ppu.dotClock++
	if ppu.dotClock == 456 {
		ppu.dotClock = 0
	}

}

//SDL Helper functions
func (ppu *PPU) drawPixel(x int, y int, colour uint32) {
	ppu.frameBuffer[x*4+(y*4*screenWidth)] = uint8((colour & 0xFF0000) >> 16)
	ppu.frameBuffer[x*4+(y*4*screenWidth)+1] = uint8((colour & 0xFF00) >> 8)
	ppu.frameBuffer[x*4+(y*4*screenWidth)+2] = uint8(colour & 0xFF)
	ppu.frameBuffer[x*4+(y*4*screenWidth)+3] = 0xFF
}

func (ppu *PPU) drawFrame() {
	ppu.renderer.Clear()
	ppu.texture.Update(nil, ppu.frameBuffer, 4*screenWidth)
	ppu.renderer.Copy(ppu.texture, nil, nil)
	ppu.renderer.Present()

}

func (ppu *PPU) clearScreen() {
	ppu.renderer.Clear()
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			ppu.drawPixel(x, y, dark)
		}
	}

	ppu.texture.Update(nil, ppu.frameBuffer, 4*screenWidth)
	ppu.renderer.Copy(ppu.texture, nil, nil)
	ppu.renderer.Present()
}


//I spent god knows how long getting this stuff to work...

func (ppu *PPU) drawTilePixel(x int, y int, colour uint32) {
	ppu.tileframeBuffer[x*4+(y*4*tilewindowWidth)] = uint8((colour & 0xFF0000) >> 16)
	ppu.tileframeBuffer[x*4+(y*4*tilewindowWidth)+1] = uint8((colour & 0xFF00) >> 8)
	ppu.tileframeBuffer[x*4+(y*4*tilewindowWidth)+2] = uint8(colour & 0xFF)
	ppu.tileframeBuffer[x*4+(y*4*tilewindowWidth)+3] = 0xFF
}

func (ppu *PPU) drawTile(bitmap []uint8, tileCoord int) {
	baseRow := (tileCoord / 16) * 8
	baseCol := (tileCoord % 16) * 8
	for row := 0; row < 8; row++ {
		byte1 := bitmap[row * 2]
		byte2 := bitmap[row * 2 + 1]
		for col := 0; col < 8; col++ {
			colour := ((byte2 >> (7-col) & 1) << 1) | (byte1 >> (7-col) & 1)
			ppu.drawTilePixel(baseCol + col, baseRow + row, colours[colour])
		}
	}
}

func (ppu *PPU) displayTileset() {
	//I ended up struggling here quite a bit because of a misunderstanding of how vram stores data
	//It stores data in TILES, not as a row-by-row bitmap
	for i := 0; i < 0x1800; i +=16 {
		ppu.drawTile(ppu.VRAM[i:i+16],i/16)
	}


	ppu.tileRenderer.Clear()
	ppu.tileTexture.Update(nil, ppu.tileframeBuffer, 4*tilewindowWidth)
	ppu.tileRenderer.Copy(ppu.tileTexture, nil, nil)
	ppu.tileRenderer.Present()
}

/*
Some self documentation just to wrap my head around these ppu concepts:

Tile Data: Every 2 bytes represents eight pixels in a row.
		   It can be accessed through two different methods,

		   1)Access tiles through ($8000 + (Tile_num * 16))
		   2)Access tiles through ($9000 + (Signed_tile_num * 16))

		   The addressing method is controlled by the LCDC register

Tile Maps: This map is used to determine which tiles are displayed to the screen,
		   in the case of the background and window grids

		   It's size is 32x32, where each byte corresponds to a given
		   tile data. There are two tile maps found at $9C00-$9FFF and $9800-$9BFF

OAM Memory: This memory location from $FE00-$FE9F holds data corresponding to
			sprites. Each sprite takes up 4 bytes, which means a total of
			40 sprites can be displayed at any given time.

			Byte 0: Y position - 16
			Byte 1: X position - 8
			Byte 2: Tile number used for the graphics of the sprite (always uses $8000 method)
			Byte 3: Sprite Flags -> For cool and snazzy effects
					Bit 7: OBJ-vs-BG priority
					Bit 6: Y flip
					Bit 5: X flip
					Bit 4: Palette number (OBP0 vs OBP1 register)
					bit 3-0: Only used by the CGB

Registers: LCDC : $FF40
		   Bit 7 => LCD Enable
		   Bit 6 => Window Tile Map select: 0 is $9800-9BFF and 1 is $9C00-9FFF
		   Bit 5 => Window enable
		   Bit 4 => Tile Data select: 0 is $8000 unsigned method, and 1 is $9000 signed method
		   Bit 3 => BG Tile Map select: 0 is $9800-9BFF and 1 is $9C00-9FFF
		   Bit 2 => Sprite size: 0 is 1x1 tiles and 1 is 1x2 tiles
		   Bit 1 => Sprite Enable
		   Bit 0 => BG and Window Enable

		   LCD Status: $FF41
		   Bit 7 => Always 1
		   Bit 6 => LYC = LY Stat interrupt enable, which allows the afformentioned condition to trigger a STAT interrupt (See Bit 2)
		   Bit 5 => Toggles Mode 2 condition to trigger a STAT interrupt
		   Bit 4 => Toggles Mode 1 condition to trigger a STAT interrupt
		   Bit 3 => Toggles Mode 0 condition to trigger a STAT interrupt
		   Bit 2 => Set if LYC = LY
		   Bit 1-0 => PPU mode:
					  0 : H-blank
					  1 : V-blank
					  2 : OAM Scan
					  3 : Drawing
*/
