package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	//Enum PPU modes
	Hblank = iota
	Vblank
	OAMSearch
	LCDTransfer 

	//Main tile window sizes
	screenWidth  = 160
	screenHeight = 144
	windowScale  = 4

	//Full window sizes
	fullwindowWidth  = 256
	fullwindowHeight = 256
	fullwindowScale  = 1

	//Debug tile window sizes
	tilewindowWidth  = 128
	tilewindowHeight = 192
	tilewindowScale  = 3
)

var (
	//Gameboy colours
	dark   = uint32(0x46425e)
	ldark  = uint32(0x5b768d)
	lwhite = uint32(0xd17c7c)
	white  = uint32(0xf6c6a8)

	colours = [4]uint32{white, lwhite, ldark, dark}
)

//PPU is the pixel processing unit of the system
//It's a custom GPU that utilises tile based rendering
type PPU struct {
	gb   *gameboy
	VRAM [8 * 1024]uint8

	//Main window
	window      *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	frameBuffer []uint8

	//Tileset (debugging)
	tileWindow      *sdl.Window
	tileRenderer    *sdl.Renderer
	tileTexture     *sdl.Texture
	tileFramebuffer []uint8

	//Full window (debugging)
	fullWindow      *sdl.Window
	fullRenderer    *sdl.Renderer
	fulltexture     *sdl.Texture
	fullFramebuffer []uint8

	//PPU internal variables
	ppuEnabled bool

	LCDC    uint8
	LCDSTAT uint8

	SCX uint8
	SCY uint8

	LY  uint8 //Used to hold line number that the scanline renderer is on
	LYC uint8

	WX uint8 //Remember this is window position - 7
	WY uint8

	//Mode of the PPU
	mode int

	dotClock int //Used to determine what the PPU should be doing
}

func initPPU(gb *gameboy) *PPU {
	ppu := new(PPU)
	ppu.gb = gb
	ppu.window, ppu.renderer = initSDL()
	if isDebugging {
		ppu.tileWindow, ppu.tileRenderer = initSDLDebugging()

		ppu.tileFramebuffer = make([]uint8, tilewindowWidth*tilewindowHeight*4) //RGBA32
		ppu.tileTexture, _ = ppu.tileRenderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, tilewindowWidth, tilewindowHeight)
		ppu.tileRenderer.SetScale(tilewindowScale, tilewindowScale)

		ppu.fullFramebuffer = make([]uint8, fullwindowWidth*fullwindowHeight*4) //RGBA32
		ppu.fulltexture, _ = ppu.fullRenderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, fullwindowWidth, fullwindowHeight)
		ppu.fullRenderer.SetScale(fullwindowScale, fullwindowScale)
	}

	ppu.renderer.SetScale(windowScale, windowScale)
	ppu.frameBuffer = make([]uint8, screenWidth*screenHeight*4) //RGBA32
	ppu.texture, _ = ppu.renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, screenWidth, screenHeight)

	ppu.mode = 2
	ppu.ppuEnabled = true

	return ppu
}

func initSDL() (*sdl.Window, *sdl.Renderer) {
	//Does the necessary setup for the SDL library
	mWindowPosX := int32(sdl.WINDOWPOS_UNDEFINED)
	mWindowPosY := int32(sdl.WINDOWPOS_UNDEFINED)

	//Initialise SDL
	err := sdl.Init(sdl.INIT_VIDEO)
	checkErr(err, "SDL initialisation error")

	//Create window
	if isDebugging {
		//To line up the windows nicely with a fullscreen termui
		mWindowPosX = 13
		mWindowPosY = 80
	}
	window, err := sdl.CreateWindow("Purpleboy!", mWindowPosX, mWindowPosY, screenWidth*windowScale, screenHeight*windowScale, sdl.WINDOW_SHOWN)
	checkErr(err, "Window creation error")
	window.SetResizable(true)

	//Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "renderer creation error")

	return window, renderer

}

func initSDLDebugging() (*sdl.Window, *sdl.Renderer, ) {
	//Initialises the required windows for debugging purposes

	tileWindow, err := sdl.CreateWindow("Debug", 660, 80, tilewindowWidth*tilewindowScale, tilewindowHeight*tilewindowScale, sdl.WINDOW_SHOWN)
	checkErr(err, "Debug window creation error")
	tileWindow.SetResizable(true)

	tileRenderer, err := sdl.CreateRenderer(tileWindow, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "Debug renderer creation error")

	//fullWindow, err := sdl.CreateWindow("Full window", 1051, 78, fullwindowWidth*fullwindowScale, fullwindowHeight*fullwindowScale, sdl.WINDOW_SHOWN)
	checkErr(err, "Debug window creation error")

	//fullRenderer, err := sdl.CreateRenderer(fullWindow, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "Debug renderer creation error")

	return tileWindow, tileRenderer //, fullWindow, fullRenderer

}

func (ppu *PPU) tick() {
	//TODO: proper cpu privileges when accessing data
	ppu.ppuEnabled = bitSet(ppu.LCDC, 7)
	if !ppu.ppuEnabled {
		ppu.mode = Hblank
		ppu.LY = 0
		return
	}

	switch ppu.mode {
	//Request LCDSTAT interrupts if corresponding bit set in LCDSTAT
	case OAMSearch: 
		if ppu.dotClock == 80 {
			ppu.dotClock = -1
			ppu.mode = LCDTransfer
		} else if isZero(ppu.dotClock){
			if bitSet(ppu.LCDSTAT,5){
				ppu.gb.cpu.IF |= 0x2
			}
		}
	case LCDTransfer:
		if ppu.dotClock == 172 {
			ppu.dotClock = -1
			ppu.mode = Hblank
			ppu.drawScanline()
		} 
	case Hblank:
		if ppu.dotClock == 204 {
			ppu.dotClock = -1
			ppu.LY++
			if ppu.LY == 144 {
				ppu.gb.cpu.IF |= 0x01 //Request vblank 
				ppu.mode = Vblank
			} else {
				ppu.mode = OAMSearch
			}
		} else if isZero(ppu.dotClock){
			if bitSet(ppu.LCDSTAT,3){
				ppu.gb.cpu.IF |= 0x2
			}
		}
	case Vblank:
		if ppu.dotClock == 456 {
			ppu.dotClock = -1
			ppu.LY++
			if ppu.LY == 154 {
				ppu.drawBuffer(ppu.renderer, ppu.texture, ppu.frameBuffer, screenWidth)
				ppu.LY = 0
				ppu.mode = OAMSearch
			}
		} else if isZero(ppu.dotClock){
			if bitSet(ppu.LCDSTAT,4){
				ppu.gb.cpu.IF |= 0x2
			}
		}

	}

	//Compare LYC and LY here every tick
	ppu.compareLYC()
	ppu.dotClock++
}

//Compare LY and LYC every tick of the PPU
func (ppu *PPU) compareLYC() {

}

func (ppu *PPU) drawScanline() {
	//Draw a scanline here
	//Since I'm getting absolutely punked by this part of the project,
	//I shall approach it slowly and cautiously, and implement things one
	//bit at a time

	tileMap := 0x1800
	if bitSet(ppu.LCDC, 3) {
		tileMap = 0x1C00
	}

	tileDataStart := uint16(0x1000) //Signed!
	if bitSet(ppu.LCDC, 4) {
		tileDataStart = 0x0000
	}

	palette := ppu.gb.mmu.readbyte(0xFF47)

	ycoordOffset := int(ppu.LY + ppu.SCY)
	row := ycoordOffset % 8                     //Which row of the tile is used for the line
	tileMapOffset := (ycoordOffset / 8) * 32   //Offset for the tilemap

	for x := 0; x < 160; x++ {
		xcoordOffset := x + int(ppu.SCX)
		tile := (xcoordOffset) / 8 //Which tile we're using for 8 bits
		col := x % 8  //Which bit from the 2 bytes are drawing //POSSIBLE PROBLEM

		tileNum := ppu.VRAM[tileMapOffset+tileMap+(tile)]
		byte1 := uint8(0)
		byte2 := uint8(0)
		if tileDataStart == 0x1000 {
			//Signed tile access
			byte1 = ppu.VRAM[tileDataStart+(uint16(int16(int8(tileNum))) * 16)+ uint16(row*2)]
			byte2 = ppu.VRAM[tileDataStart+(uint16(int16(int8(tileNum))) * 16)+ uint16(row*2) + 1]
		} else {
			//Regular tile access
			byte1 = ppu.VRAM[tileDataStart+(uint16(tileNum) * 16) + uint16(row*2)]
			byte2 = ppu.VRAM[tileDataStart+(uint16(tileNum) * 16) + uint16(row*2) + 1]
		}

		//Get colour from palette
		colourIndex := ((byte2 >> (7 - col) & 1) << 1) | (byte1 >> (7 - col) & 1)
		colour := colours[(palette>>(colourIndex*2))&0x3]

		ppu.drawPixel(ppu.frameBuffer, screenWidth, x, int(ppu.LY), colour)

	}
}

//Draws a pixel on a buffer
func (ppu *PPU) drawPixel(buffer []uint8, lineWidth int, x int, y int, colour uint32) {
	buffer[x*4+(y*4*lineWidth)] = uint8((colour & 0xFF0000) >> 16)
	buffer[x*4+(y*4*lineWidth)+1] = uint8((colour & 0xFF00) >> 8)
	buffer[x*4+(y*4*lineWidth)+2] = uint8(colour & 0xFF)
	buffer[x*4+(y*4*lineWidth)+3] = 0xFF
}

func (ppu *PPU) drawBuffer(renderer *sdl.Renderer, texture *sdl.Texture, buffer []uint8, lineWidth int) {
	renderer.Clear()
	texture.Update(nil, buffer, 4*lineWidth)
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

//Render background
func (ppu *PPU) drawTile(framebuffer []uint8, lineWidth int, bitmap []uint8, tileCoord int) {
	//Get proper offsets for drawing the current tile
	baseRow := (tileCoord / (lineWidth / 8)) * 8
	baseCol := (tileCoord % (lineWidth / 8)) * 8

	//Begin drawing tile from supplied bitmap
	for row := 0; row < 8; row++ {
		byte1 := bitmap[row*2]
		byte2 := bitmap[row*2+1]
		for col := 0; col < 8; col++ {
			//Get colour from palette
			palette := ppu.gb.mmu.readbyte(0xFF47)
			colourIndex := ((byte2 >> (7 - col) & 1) << 1) | (byte1 >> (7 - col) & 1)
			colour := colours[(palette>>(colourIndex*2))&0x3]
			ppu.drawPixel(framebuffer, lineWidth, baseCol+col, baseRow+row, colour)
		}
	}
}

func (ppu *PPU) displayTileset() {
	//I ended up struggling here quite a bit because of a misunderstanding of how vram stores data
	//It stores data in TILES, not as a row-by-row bitmap
	for i := 0; i < 0x1800; i += 16 {
		ppu.drawTile(ppu.tileFramebuffer, tilewindowWidth, ppu.VRAM[i:i+16], i/16)
	}

	ppu.drawBuffer(ppu.tileRenderer, ppu.tileTexture, ppu.tileFramebuffer, tilewindowWidth)
}

func (ppu *PPU) displayCurrTileMap(){
	for i := 0; i < 32*32; i++ {
		//Loop through each of the 32 x 32 tiles in one of the tilemaps

		tileMap := 0x1800 //0x9800
		if bitSet(ppu.LCDC,3) {
			tileMap = 0x1C00 //0x9C00
		}

		tileDataStart := uint16(0x1000) //0x9000
		if bitSet(ppu.LCDC, 4) { 
			tileDataStart = 0x0000 
		}

		tileDataStart = uint16(0x1000) //Temp fix?

		
		if tileDataStart == 0x0000 {
			//Unsigned tile access
			tileNum := ppu.VRAM[tileMap + i] //Get tile num from tilemap
			tileDataIndex := tileNum * 16 //Get index to start of tile in vram
			ppu.drawTile(ppu.fullFramebuffer,fullwindowWidth,ppu.VRAM[tileDataIndex:tileDataIndex+16],i)
		} else if tileDataStart == 0x1000 {
			//Signed tile access
			tileNum := ppu.VRAM[tileMap + i] //Get tile num from tilemap
			tileDataIndex := tileDataStart + uint16(int16(int8(uint8(tileNum)))) * 16
			ppu.drawTile(ppu.fullFramebuffer,fullwindowWidth,ppu.VRAM[tileDataIndex:tileDataIndex+16],i)  
		}
	}
	ppu.drawBuffer(ppu.fullRenderer, ppu.fulltexture, ppu.fullFramebuffer, fullwindowWidth)
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
		   Bit 4 => Tile Data select: 1 is $8000 unsigned method, and 0 is $9000 signed method
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
