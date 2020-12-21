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

	//PPU internal variables
	ppuEnabled bool

	palette        uint8 //For background
	spritePalette1 uint8 //For sprites
	spritePalette2 uint8 //For sprites
	LCDC           uint8
	LCDSTAT        uint8

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
	checkErr(sdl.Init(sdl.INIT_VIDEO|sdl.INIT_JOYSTICK), "SDL initialisation error")

	//Create window
	if isDebugging {
		//To line up the windows nicely with a fullscreen termui
		mWindowPosX = 13
		mWindowPosY = 80
	}
	window, err := sdl.CreateWindow("Purpleboy!", mWindowPosX, mWindowPosY, screenWidth*windowScale, screenHeight*windowScale, sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	checkErr(err, "Window creation error")

	//Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "renderer creation error")

	return window, renderer

}

func initSDLDebugging() (*sdl.Window, *sdl.Renderer) {
	//Initialises the required windows for debugging purposes

	tileWindow, err := sdl.CreateWindow("Debug", 660, 80, tilewindowWidth*tilewindowScale, tilewindowHeight*tilewindowScale, sdl.WINDOW_SHOWN|sdl.WINDOW_ALWAYS_ON_TOP)
	checkErr(err, "Debug window creation error")

	tileRenderer, err := sdl.CreateRenderer(tileWindow, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err, "Debug renderer creation error")

	return tileWindow, tileRenderer
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
		} else if isZero(ppu.dotClock) {
			if bitSet(ppu.LCDSTAT, 5) {
				ppu.gb.cpu.requestSTAT()
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
				ppu.gb.cpu.requestVblank()
				ppu.mode = Vblank
			} else {
				ppu.mode = OAMSearch
			}
			ppu.compareLYC()
		} else if isZero(ppu.dotClock) {
			if bitSet(ppu.LCDSTAT, 3) {
				ppu.gb.cpu.requestSTAT()
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
			ppu.compareLYC()
		} else if isZero(ppu.dotClock) {
			if bitSet(ppu.LCDSTAT, 4) {
				ppu.gb.cpu.requestSTAT()
			}
		}
	}
	ppu.dotClock++
}

//Compare LY and LYC every tick of the PPU
func (ppu *PPU) compareLYC() {
	//Set bit 2 depending on lyc == ly
	if ppu.LYC == ppu.LY {
		ppu.LCDSTAT |= 0x04

		//If LYC=LY interrupt enable, request
		if bitSet(ppu.LCDSTAT, 6) {
			ppu.gb.cpu.requestSTAT()
		}
	} else {
		ppu.LCDSTAT &^= 0x04
	}
}

func (ppu *PPU) drawScanline() {
	ppu.drawBG()
	ppu.drawSprites()
}

func (ppu *PPU) drawBG() {
	//Draw a scanline here
	if !bitSet(ppu.LCDC, 0) { //Basically background enable master flag
		return
	}
	usingWindow := false
	windowX := uint8(0)
	if ppu.WX == 7 { //Pesky underflows!
		windowX = 0 //TODO: Fix one width pixel offset in Aladdin
	} else if ppu.WX <= 6 {
		windowX = 6 - ppu.WX
	} else {
		windowX = ppu.WX - 7
	}

	if bitSet(ppu.LCDC, 5) && ppu.WY <= ppu.LY {
		//For future reference: We check the line to see if it has entered the window region, since ppu.WY is usually fixed (i think)
		usingWindow = true
	}

	tileMap := 0x1800               //Both BG and window use 0x1800 as the default map
	tileDataStart := uint16(0x1000) //Signed!
	if bitSet(ppu.LCDC, 4) {
		tileDataStart = 0x0000
	}

	for x := 0; x < 160; x++ {
		xcoordOffset := int(0)
		ycoordOffset := int(0)
		if usingWindow && windowX <= uint8(x) {
			//Window
			ycoordOffset = int(ppu.LY - ppu.WY)
			xcoordOffset = int(uint8(x) - windowX)

			if bitSet(ppu.LCDC, 6) { //Window tilemap select
				tileMap = 0x1C00
			}
		} else {
			//BG
			xcoordOffset = int(uint8(x) + ppu.SCX)
			ycoordOffset = int(ppu.LY + ppu.SCY)

			if bitSet(ppu.LCDC, 3) { //BG tilemap select
				tileMap = 0x1C00
			}
		}

		row := ycoordOffset % 8                  //Which row of the tile is used for the line
		tileMapOffset := (ycoordOffset / 8) * 32 //Offset for the tilemap

		tile := xcoordOffset / 8 //Which tile we're using from tile map
		col := xcoordOffset % 8  //Which bit from the 2 bytes are drawing

		tileNum := ppu.VRAM[tileMapOffset+tileMap+tile]
		byte1 := uint8(0)
		byte2 := uint8(0)
		if tileDataStart == 0x1000 {
			//Signed tile access
			target := (tileDataStart + (uint16(int16(int8(tileNum))) * 16) + uint16(row*2))
			byte1 = ppu.VRAM[target]
			byte2 = ppu.VRAM[target+1]
		} else {
			//Regular tile access
			target := (tileDataStart + (uint16(tileNum) * 16) + uint16(row*2))
			byte1 = ppu.VRAM[target]
			byte2 = ppu.VRAM[target+1]
		}

		//Get colour from palette
		colourIndex := ((byte2 >> (7 - col) & 1) << 1) | (byte1 >> (7 - col) & 1)
		colour := colours[(ppu.palette>>(colourIndex*2))&0x3]

		ppu.drawPixel(ppu.frameBuffer, screenWidth, x, int(ppu.LY), colour)
	}
}

func (ppu *PPU) drawSprites() {
	if !bitSet(ppu.LCDC, 0) { //Sprite Enable
		return
	}

	sprites := [][]uint8{}
	for i := 0; i < 0xA0; i += 4 {
		//TODO: Fix x <= 8 and y <= 16
		if len(sprites) <= 10 {
			spriteSize := uint8(8)
			if bitSet(ppu.LCDC, 2) {
				spriteSize = 16
			}
			y := ppu.gb.mmu.OAM[i] - 16
			x := ppu.gb.mmu.OAM[i+1] - 8 
			if y <= ppu.LY && ppu.LY < y+spriteSize {
				sprites = append(sprites, []uint8{y, x, ppu.gb.mmu.OAM[i+2], ppu.gb.mmu.OAM[i+3]})
			}
		} else {
			break
		}
	}

	for i := 0; i < len(sprites); i++ {
		x := int(sprites[i][1])
		y := sprites[i][0]
		tileNum := sprites[i][2]
		attrs := sprites[i][3]

		spritePalSelect := bitSet(attrs, 4)
		xflip := bitSet(attrs, 5)
		bgPriority := bitSet(attrs, 7)

		row := uint16(ppu.LY - y) //TOOD: y flip
		byte1 := ppu.VRAM[(uint16(tileNum)*16)+(row*2)]
		byte2 := ppu.VRAM[(uint16(tileNum)*16)+(row*2)+1]

		if x >= 0 && x <= 160 {
			for col := 0; col < 8; col++ {
				colour := uint32(0)
				colourIndex := uint8(0)
				if xflip {
					colourIndex = ((byte2 >> col & 1) << 1) | (byte1 >> col & 1)
				} else {
					colourIndex = ((byte2 >> (7 - col) & 1) << 1) | (byte1 >> (7 - col) & 1)
				}

				if !spritePalSelect { //Sprite Palette select
					colour = colours[(ppu.spritePalette1>>(colourIndex*2))&0x3]
				} else {
					colour = colours[(ppu.spritePalette2>>(colourIndex*2))&0x3]
				}

				if !bgPriority { //BG-OBJ priority
					if x+col <= 160 && colourIndex != 0 {
						ppu.drawPixel(ppu.frameBuffer, screenWidth, x+col, int(ppu.LY), colour)
					}
				} else {
					if x+col <= 160 && colourIndex != 0 && ppu.getPixelColour(x+col, int(ppu.LY)) == white {
						ppu.drawPixel(ppu.frameBuffer, screenWidth, x+col, int(ppu.LY), colour)
					}
				}
			}
		}
	}
}

//Draws a pixel on a buffer
func (ppu *PPU) drawPixel(buffer []uint8, lineWidth int, x int, y int, colour uint32) {
	buffer[x*4+(y*4*lineWidth)] = uint8((colour & 0xFF0000) >> 16)
	buffer[x*4+(y*4*lineWidth)+1] = uint8((colour & 0xFF00) >> 8)
	buffer[x*4+(y*4*lineWidth)+2] = uint8(colour & 0xFF)
	buffer[x*4+(y*4*lineWidth)+3] = 0xFF
}

func (ppu *PPU) getPixelColour(x int, y int) uint32 {
	//Gets the colour from a given coordinate
	byte1 := ppu.frameBuffer[x*4+(y*4*screenWidth)]
	byte2 := ppu.frameBuffer[x*4+(y*4*screenWidth)+1]
	byte3 := ppu.frameBuffer[x*4+(y*4*screenWidth)+2]

	return uint32(byte1)<<16 | uint32(byte2)<<8 | uint32(byte3)
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
			colourIndex := ((byte2 >> (7 - col) & 1) << 1) | (byte1 >> (7 - col) & 1)
			colour := colours[(ppu.palette>>(colourIndex*2))&0x3]
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
