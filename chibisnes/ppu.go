package chibisnes

type BGLayer struct {
	hScroll       uint16
	vScroll       uint16
	tilemapWider  bool
	tilemapHigher bool
	tilemapAddr   uint16
	tileAddr      uint16
	bigTiles      bool
	mosaicEnabled bool
}

type Layer struct {
	mainScreenEnabled  bool
	subScreenEnabled   bool
	mainScreenWindowed bool
	subScreenWindowed  bool
}

type WindowLayer struct {
	window1Enabled  bool
	window2Enabled  bool
	window1Inversed bool
	window2Inversed bool
	maskLogic       byte
}

type PPU struct {
	console *Console

	// VRAM access
	vram                [0x8000]uint16
	vramPointer         uint16
	vramIncrementOnHigh bool
	vramIncrement       uint16
	vramRemapMode       byte
	vramReadBuffer      uint16

	// CGRAM access
	cgram            [0x100]uint16
	cgramPointer     byte
	cgramSecondWrite bool
	cgramBuffer      byte

	// OAM access
	oam              [0x100]uint16
	highOAM          [0x20]byte
	oamAddr          byte
	oamAddrWritten   byte
	oamInHigh        bool
	oamInHighWritten bool
	oamSecondWrite   bool
	oamBuffer        byte

	// object, sprites
	objPriority       bool
	objTileAddr1      uint16
	objTileAddr2      uint16
	objSize           byte
	objPixelBuffer    [256]byte // line buffers
	objPriorityBuffer [256]byte
	timeOver          bool
	rangeOver         bool
	objInterlace      bool

	// background layers
	bgLayer         [4]BGLayer
	scrollPrev      byte
	scrollPrev2     byte
	mosaicSize      byte
	mosaicStartLine byte

	// layers
	layer [5]Layer

	// mode 7
	mode7Matrix     [8]int16 // a, b, c, d, x, y, h, v
	mode7Prev       byte
	mode7LargeField bool
	mode7CharFill   bool
	mode7XFlip      bool
	mode7YFlip      bool
	mode7ExtBG      bool

	// mode 7 internal
	mode7StartX int32
	mode7StartY int32

	// windows
	windowLayer  [6]WindowLayer
	window1Left  byte
	window1Right byte
	window2Left  byte
	window2Right byte

	// color math
	clipMode        byte
	preventMathMode byte
	addSubscreen    bool
	subtractColor   bool
	halfColor       bool
	mathEnabled     [6]bool
	fixedColorR     byte
	fixedColorG     byte
	fixedColorB     byte

	// settings
	forcedBlank    bool
	brightness     byte
	mode           byte
	bg3priority    bool
	evenFrame      bool
	pseudoHires    bool
	overscan       bool
	frameOverscan  bool
	interlace      bool
	frameInterlace bool
	directColor    bool

	// latching
	hCount          uint16
	vCount          uint16
	hCountSecond    bool
	vCountSecond    bool
	countersLatched bool
	ppu1OpenBus     byte
	ppu2OpenBus     byte

	// pixel buffer (xbgr)
	// times 2 for event and odd frame
	pixelBuffer [512 * 4 * 239 * 2]byte
}

// array for layer definitions per mode:
//   0-7: mode 0-7; 8: mode 1 + l3prio; 9: mode 7 + extbg

//   0-3; layers 1-4; 4: sprites; 5: nonexistent
var layersPerMode [10][12]int = [10][12]int{
	{4, 0, 1, 4, 0, 1, 4, 2, 3, 4, 2, 3},
	{4, 0, 1, 4, 0, 1, 4, 2, 4, 2, 5, 5},
	{4, 0, 4, 1, 4, 0, 4, 1, 5, 5, 5, 5},
	{4, 0, 4, 1, 4, 0, 4, 1, 5, 5, 5, 5},
	{4, 0, 4, 1, 4, 0, 4, 1, 5, 5, 5, 5},
	{4, 0, 4, 1, 4, 0, 4, 1, 5, 5, 5, 5},
	{4, 0, 4, 4, 0, 4, 5, 5, 5, 5, 5, 5},
	{4, 4, 4, 0, 4, 5, 5, 5, 5, 5, 5, 5},
	{2, 4, 0, 1, 4, 0, 1, 4, 4, 2, 5, 5},
	{4, 4, 1, 4, 0, 4, 1, 5, 5, 5, 5, 5},
}

var prioritysPerMode [10][12]int = [10][12]int{
	{3, 1, 1, 2, 0, 0, 1, 1, 1, 0, 0, 0},
	{3, 1, 1, 2, 0, 0, 1, 1, 0, 0, 5, 5},
	{3, 1, 2, 1, 1, 0, 0, 0, 5, 5, 5, 5},
	{3, 1, 2, 1, 1, 0, 0, 0, 5, 5, 5, 5},
	{3, 1, 2, 1, 1, 0, 0, 0, 5, 5, 5, 5},
	{3, 1, 2, 1, 1, 0, 0, 0, 5, 5, 5, 5},
	{3, 1, 2, 1, 0, 0, 5, 5, 5, 5, 5, 5},
	{3, 2, 1, 0, 0, 5, 5, 5, 5, 5, 5, 5},
	{1, 3, 1, 1, 2, 0, 0, 1, 0, 0, 5, 5},
	{3, 2, 1, 1, 0, 0, 0, 5, 5, 5, 5, 5},
}

var layerCountPerMode [10]int = [10]int{
	12, 10, 8, 8, 8, 8, 6, 5, 10, 7,
}

var bitDepthsPerMode [10][4]int = [10][4]int{
	{2, 2, 2, 2},
	{4, 4, 2, 5},
	{4, 4, 5, 5},
	{8, 4, 5, 5},
	{8, 2, 5, 5},
	{4, 2, 5, 5},
	{4, 5, 5, 5},
	{8, 5, 5, 5},
	{4, 4, 2, 5},
	{8, 7, 5, 5},
}

var spriteSizes [8][2]int = [8][2]int{
	{8, 16}, {8, 32}, {8, 64}, {16, 32},
	{16, 64}, {32, 64}, {16, 32}, {16, 32},
}

func NewPPU(console *Console) *PPU {
	return &PPU{
		console: console,
	}
}

func (ppu *PPU) Reset() {
	ppu.vramIncrement = 1
	ppu.mosaicSize = 1
	ppu.mosaicStartLine = 1
	ppu.forcedBlank = true
}

func (ppu *PPU) checkOverscan() bool {
	ppu.frameOverscan = ppu.overscan // set if we have a overscan-frame
	return ppu.frameOverscan
}

func (ppu *PPU) handleVBlank() {
	// called either right after checkOverscan at (0, 255), or at (0, 240)
	if !ppu.forcedBlank {
		ppu.oamAddr = ppu.oamAddrWritten
		ppu.oamInHigh = ppu.oamInHighWritten
		ppu.oamSecondWrite = false
	}
	ppu.frameInterlace = ppu.interlace // set if we have a interlaced frame
}

func (ppu *PPU) runLine(line int) {
	if line == 0 {
		// pre-render line
		// TODO: this now happens halfway into the first line
		ppu.mosaicStartLine = 1
		ppu.rangeOver = false
		ppu.timeOver = false
		ppu.evenFrame = !ppu.evenFrame
		if !ppu.forcedBlank {
			ppu.evaluateSprites(0)
		}
	} else {
		for i := 0; i < len(ppu.objPixelBuffer); i++ {
			ppu.objPixelBuffer[i] = 0
		}
		// evaluate sprites
		if !ppu.forcedBlank {
			ppu.evaluateSprites(line - 1)
		}
		// actual line
		if ppu.mode == 7 {
			ppu.calculateMode7Starts(line)
		}
		for x := 0; x < 256; x++ {
			ppu.handlePixel(x, line)
		}
	}
}

func (ppu *PPU) handlePixel(x int, y int) {
	var r, r2 int = 0, 0
	var g, g2 int = 0, 0
	var b, b2 int = 0, 0
	if !ppu.forcedBlank {
		var mainLayer int = ppu.getPixel(x, y, false, &r, &g, &b)
		var colorWindowState bool = ppu.getWindowState(5, x)
		if ppu.clipMode == 3 ||
			(ppu.clipMode == 2 && colorWindowState) ||
			(ppu.clipMode == 1 && !colorWindowState) {
			r = 0
			g = 0
			b = 0
		}
		var secondLayer int = 5 // backdrop
		var mathEnabled bool = mainLayer < 6 && ppu.mathEnabled[mainLayer] && !(ppu.preventMathMode == 3 ||
			(ppu.preventMathMode == 2 && colorWindowState) ||
			(ppu.preventMathMode == 1 && !colorWindowState))

		if (mathEnabled && ppu.addSubscreen) || ppu.pseudoHires || ppu.mode == 5 || ppu.mode == 6 {
			secondLayer = ppu.getPixel(x, y, true, &r2, &g2, &b2)
		}
		// TODO: subscreen pixels can be clipped to black as well
		// TODO: math for subscreen pixels (add/sub sub to main)
		if mathEnabled {
			if ppu.subtractColor {
				if ppu.addSubscreen && secondLayer != 5 {
					r -= r2
					g -= g2
					b -= b2
				} else {
					r -= int(ppu.fixedColorR)
					g -= int(ppu.fixedColorG)
					b -= int(ppu.fixedColorB)
				}
			} else {
				if ppu.addSubscreen && secondLayer != 5 {
					r += r2
					g += g2
					b += b2
				} else {
					r += int(ppu.fixedColorR)
					g += int(ppu.fixedColorG)
					b += int(ppu.fixedColorB)
				}
			}
			if ppu.halfColor && (secondLayer != 5 || !ppu.addSubscreen) {
				r >>= 1
				g >>= 1
				b >>= 1
			}
			if r > 31 {
				r = 31
			}
			if g > 31 {
				g = 31
			}
			if b > 31 {
				b = 31
			}
			if r < 0 {
				r = 0
			}
			if g < 0 {
				g = 0
			}
			if b < 0 {
				b = 0
			}
		}
		if !(ppu.pseudoHires || ppu.mode == 5 || ppu.mode == 6) {
			r2 = r
			g2 = g
			b2 = b
		}
	}

	var row int
	if ppu.evenFrame {
		row = (y - 1) + 0
	} else {
		row = (y - 1) + 239
	}

	ppu.pixelBuffer[row*2048+x*8+0] = byte(((r2 << 3) | (r2 >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+1] = byte(((g2 << 3) | (g2 >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+2] = byte(((b2 << 3) | (b2 >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+3] = 0xFF
	ppu.pixelBuffer[row*2048+x*8+4] = byte(((r << 3) | (r >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+5] = byte(((g << 3) | (g >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+6] = byte(((b << 3) | (b >> 2)) * int(ppu.brightness) / 15)
	ppu.pixelBuffer[row*2048+x*8+7] = 0xFF
}

func (ppu *PPU) getPixel(x int, y int, sub bool, r *int, g *int, b *int) int {
	// figure out which color is on this location on main- or subscreen, sets it in r, g, b
	// returns which layer it is: 0-3 for bg layer, 4 or 6 for sprites (depending on palette), 5 for backdrop
	var actMode int
	if ppu.mode == 1 && ppu.bg3priority {
		actMode = 8
	} else {
		actMode = int(ppu.mode)
	}
	if ppu.mode == 7 && ppu.mode7ExtBG {
		actMode = 9
	}

	var layer int = 5
	var pixel int = 0
	for i := 0; i < layerCountPerMode[actMode]; i++ {
		var curLayer int = layersPerMode[actMode][i]
		var curPriority int = prioritysPerMode[actMode][i]
		var layerActive bool = false
		if !sub {
			layerActive = ppu.layer[curLayer].mainScreenEnabled && (!ppu.layer[curLayer].mainScreenWindowed || !ppu.getWindowState(curLayer, x))
		} else {
			layerActive = ppu.layer[curLayer].subScreenEnabled && (!ppu.layer[curLayer].subScreenWindowed || !ppu.getWindowState(curLayer, x))
		}
		if layerActive {
			if curLayer < 4 {
				// bg layer
				var lx int = x
				var ly int = y
				if ppu.bgLayer[curLayer].mosaicEnabled && ppu.mosaicSize > 1 {
					lx -= lx % int(ppu.mosaicSize)
					ly -= (ly - int(ppu.mosaicStartLine)) % int(ppu.mosaicSize)
				}
				if ppu.mode == 7 {
					pixel = ppu.getPixelForMode7(lx, curLayer, curPriority > 0)
				} else {
					lx += int(ppu.bgLayer[curLayer].hScroll)
					if ppu.mode == 5 || ppu.mode == 6 {
						lx *= 2
						if sub || ppu.bgLayer[curLayer].mosaicEnabled {
							lx += 0
						} else {
							lx += 1
						}
						if ppu.interlace {
							ly *= 2
							if ppu.evenFrame || ppu.bgLayer[curLayer].mosaicEnabled {
								ly += 0
							} else {
								ly += 1
							}
						}
					}
					ly += int(ppu.bgLayer[curLayer].vScroll)
					if ppu.mode == 2 || ppu.mode == 4 || ppu.mode == 6 {
						ppu.handleOPT(curLayer, &lx, &ly)
					}
					pixel = ppu.getPixelForBGLayer(
						lx&0x3ff, ly&0x3ff,
						curLayer, curPriority > 0)
				}
			} else {
				// get a pixel from the sprite buffer
				pixel = 0
				if int(ppu.objPriorityBuffer[x]) == curPriority {
					pixel = int(ppu.objPixelBuffer[x])
				}
			}
		}
		if pixel > 0 {
			layer = curLayer
			break
		}
	}
	if ppu.directColor && layer < 4 && bitDepthsPerMode[actMode][layer] == 8 {
		*r = ((pixel & 0x7) << 2) | ((pixel & 0x100) >> 7)
		*g = ((pixel & 0x38) >> 1) | ((pixel & 0x200) >> 8)
		*b = ((pixel & 0xc0) >> 3) | ((pixel & 0x400) >> 8)
	} else {
		var color uint16 = ppu.cgram[pixel&0xff]
		*r = int(color & 0x1f)
		*g = int((color >> 5) & 0x1f)
		*b = int((color >> 10) & 0x1f)
	}
	if layer == 4 && pixel < 0xc0 {
		layer = 6 // sprites with palette color < 0xc0
	}
	return layer
}

func (ppu *PPU) handleOPT(layer int, lx *int, ly *int) {
	var x int = *lx
	var y int = *ly
	var column int = 0

	if ppu.mode == 6 {
		column = ((x - (x & 0xf)) - ((int(ppu.bgLayer[layer].hScroll) * 2) & 0xfff0)) >> 4
	} else {
		column = ((x - (x & 0x7)) - (int(ppu.bgLayer[layer].hScroll) & 0xfff8)) >> 3
	}

	if column > 0 {
		// fetch offset values from layer 3 tilemap
		var valid int
		if layer == 0 {
			valid = 0x2000
		} else {
			valid = 0x4000
		}
		var hOffset uint16 = ppu.getOffsetValue(column-1, 0)
		var vOffset uint16 = 0
		if ppu.mode == 4 {
			if (hOffset & 0x8000) > 0 {
				vOffset = hOffset
				hOffset = 0
			}
		} else {
			vOffset = ppu.getOffsetValue(column-1, 1)
		}
		if ppu.mode == 6 {
			// TODO: not sure if correct
			if (int(hOffset) & valid) > 0 {
				*lx = (((int(hOffset) & 0x3f8) + (column * 8)) * 2) | (x & 0xf)
			}
		} else {
			if (int(hOffset) & valid) > 0 {
				*lx = ((int(hOffset) & 0x3f8) + (column * 8)) | (x & 0x7)
			}
		}

		// TODO: not sure if correct for interlace
		if (int(vOffset) & valid) > 0 {
			*ly = (int(vOffset) & 0x3ff) + (y - int(ppu.bgLayer[layer].vScroll))
		}
	}
}

func (ppu *PPU) getOffsetValue(col int, row int) uint16 {
	var x int = col*8 + int(ppu.bgLayer[2].hScroll)
	var y int = row*8 + int(ppu.bgLayer[2].vScroll)
	var tileBits int
	var tileHighBit int
	if ppu.bgLayer[2].bigTiles {
		tileBits = 4
		tileHighBit = 0x200
	} else {
		tileBits = 3
		tileHighBit = 0x100
	}
	var tilemapAddr uint16 = ppu.bgLayer[2].tilemapAddr + uint16(((y>>tileBits)&0x1f)<<5|((x>>tileBits)&0x1f))
	if (x&int(tileHighBit)) > 0 && ppu.bgLayer[2].tilemapWider {
		tilemapAddr += 0x400
	}
	if (y&int(tileHighBit)) > 0 && ppu.bgLayer[2].tilemapHigher {
		if ppu.bgLayer[2].tilemapWider {
			tilemapAddr += 0x800
		} else {
			tilemapAddr += 0x400
		}
	}
	return ppu.vram[tilemapAddr&0x7fff]
}

func (ppu *PPU) getPixelForBGLayer(x int, y int, layer int, priority bool) int {
	// figure out address of tilemap word and read it
	var wideTiles bool = ppu.bgLayer[layer].bigTiles || ppu.mode == 5 || ppu.mode == 6
	var tileBitsX int
	var tileHighBitX int
	if wideTiles {
		tileBitsX = 4
		tileHighBitX = 0x200
	} else {
		tileBitsX = 3
		tileHighBitX = 0x100
	}
	var tileBitsY int
	var tileHighBitY int
	if ppu.bgLayer[layer].bigTiles {
		tileBitsY = 4
		tileHighBitY = 0x200
	} else {
		tileBitsY = 3
		tileHighBitY = 0x100
	}

	var tilemapAddr uint16 = uint16(int(ppu.bgLayer[layer].tilemapAddr) + (((y>>tileBitsY)&0x1f)<<5 | ((x >> tileBitsX) & 0x1f)))

	if (x&tileHighBitX) > 0 && ppu.bgLayer[layer].tilemapWider {
		tilemapAddr += 0x400
	}

	if (y&tileHighBitY) > 0 && ppu.bgLayer[layer].tilemapHigher {
		if ppu.bgLayer[layer].tilemapWider {
			tilemapAddr += 0x800
		} else {
			tilemapAddr += 0x400
		}
	}

	var tile uint16 = ppu.vram[tilemapAddr&0x7fff]
	// check priority, get palette
	if ((tile & 0x2000) > 0) != priority {
		return 0 // wrong priority
	}
	var paletteNum int = (int(tile) & 0x1c00) >> 10
	// figure out position within tile
	var row int
	if (tile & 0x8000) > 0 {
		row = 7 - (y & 0x7)
	} else {
		row = (y & 0x7)
	}
	var col int
	if (tile & 0x4000) > 0 {
		col = (x & 0x7)
	} else {
		col = 7 - (x & 0x7)
	}
	var tileNum int = int(tile) & 0x3ff
	if wideTiles {
		// if unflipped right half of tile, or flipped left half of tile
		// (y & 8) xor (tile & 0x4000)
		if ((x & 8) > 0) != ((tile & 0x4000) > 0) {
			tileNum += 1
		}
	}
	if ppu.bgLayer[layer].bigTiles {
		// if unflipped bottom half of tile, or flipped upper half of tile
		// (y & 8) xor (tile & 0x8000)
		if ((y & 8) > 0) != ((tile & 0x8000) > 0) {
			tileNum += 0x10
		}
	}
	// read tiledata, ajust palette for mode 0
	var bitDepth int = bitDepthsPerMode[ppu.mode][layer]
	if ppu.mode == 0 {
		paletteNum += 8 * layer
	}
	// plane 1 (always)
	var paletteSize int = 4
	var plane1 uint16 = ppu.vram[(int(ppu.bgLayer[layer].tileAddr)+((tileNum&0x3ff)*4*bitDepth)+row)&0x7fff]
	var pixel int = (int(plane1) >> col) & 1
	pixel |= ((int(plane1) >> (8 + col)) & 1) << 1

	// plane 2 (for 4bpp, 8bpp)
	if bitDepth > 2 {
		paletteSize = 16
		var plane2 uint16 = ppu.vram[(int(ppu.bgLayer[layer].tileAddr)+((tileNum&0x3ff)*4*bitDepth)+8+row)&0x7fff]
		pixel |= ((int(plane2) >> col) & 1) << 2
		pixel |= ((int(plane2) >> (8 + col)) & 1) << 3
	}

	// plane 3 & 4 (for 8bpp)
	if bitDepth > 4 {
		paletteSize = 256
		var plane3 uint16 = ppu.vram[(int(ppu.bgLayer[layer].tileAddr)+((tileNum&0x3ff)*4*bitDepth)+16+row)&0x7fff]
		pixel |= ((int(plane3) >> col) & 1) << 4
		pixel |= ((int(plane3) >> (8 + col)) & 1) << 5
		var plane4 uint16 = ppu.vram[(int(ppu.bgLayer[layer].tileAddr)+((tileNum&0x3ff)*4*bitDepth)+24+row)&0x7fff]
		pixel |= ((int(plane4) >> col) & 1) << 6
		pixel |= ((int(plane4) >> (8 + col)) & 1) << 7
	}

	// return cgram index, or 0 if transparent, palette number in bits 10-8 for 8-color layers
	if pixel == 0 {
		return 0
	}
	return paletteSize*paletteNum + pixel
}

func (ppu *PPU) calculateMode7Starts(y int) {
	// expand 13-bit values to signed values
	// cast to int16
	var hScroll int = int((int16(ppu.mode7Matrix[6] << 3)) >> 3)
	var vScroll int = int((int16(ppu.mode7Matrix[7] << 3)) >> 3)
	var xCenter int = int((int16(ppu.mode7Matrix[4] << 3)) >> 3)
	var yCenter int = int((int16(ppu.mode7Matrix[5] << 3)) >> 3)
	// do calculation
	var clippedH int = hScroll - xCenter
	var clippedV int = vScroll - yCenter
	if (clippedH & 0x2000) > 0 {
		clippedH = (clippedH | ^1023)
	} else {
		clippedH = (clippedH & 1023)
	}
	if (clippedV & 0x2000) > 0 {
		clippedV = (clippedV | ^1023)
	} else {
		clippedV = (clippedV & 1023)
	}
	if ppu.bgLayer[0].mosaicEnabled && ppu.mosaicSize > 1 {
		y -= (y - int(ppu.mosaicStartLine)) % int(ppu.mosaicSize)
	}
	var ry byte
	if ppu.mode7YFlip {
		ry = byte(255 - y)
	} else {
		ry = byte(y)
	}
	ppu.mode7StartX = int32(((int(ppu.mode7Matrix[0]) * clippedH) & ^63) +
		((int(ppu.mode7Matrix[1]) * int(ry)) & ^63) +
		((int(ppu.mode7Matrix[1]) * clippedV) & ^63) +
		(xCenter << 8))
	ppu.mode7StartY = int32(((int(ppu.mode7Matrix[2]) * clippedH) & ^63) +
		((int(ppu.mode7Matrix[3]) * int(ry)) & ^63) +
		((int(ppu.mode7Matrix[3]) * clippedV) & ^63) +
		(yCenter << 8))
}

func (ppu *PPU) getPixelForMode7(x int, layer int, priority bool) int {
	var rx byte
	if ppu.mode7XFlip {
		rx = byte(255 - x)
	} else {
		rx = byte(x)
	}
	var xPos int = int((ppu.mode7StartX + (int32(ppu.mode7Matrix[0]) * int32(rx))) >> 8)
	var yPos int = int((ppu.mode7StartY + (int32(ppu.mode7Matrix[2]) * int32(rx))) >> 8)

	var outsideMap bool = xPos < 0 || xPos >= 1024 || yPos < 0 || yPos >= 1024
	xPos &= 0x3ff
	yPos &= 0x3ff
	if !ppu.mode7LargeField {
		outsideMap = false
	}
	var tile byte
	if outsideMap {
		tile = 0
	} else {
		tile = byte(ppu.vram[(yPos>>3)*128+(xPos>>3)] & 0xff)
	}
	var pixel byte
	if outsideMap && !ppu.mode7CharFill {
		pixel = 0
	} else {
		pixel = byte(ppu.vram[int(tile)*64+(yPos&7)*8+(xPos&7)] >> 8)
	}
	if layer == 1 {
		if ((pixel & 0x80) > 0) != priority {
			return 0
		}
		return int(pixel & 0x7f)
	}
	return int(pixel)
}

func (ppu *PPU) getWindowState(layer int, x int) bool {
	if !ppu.windowLayer[layer].window1Enabled && !ppu.windowLayer[layer].window2Enabled {
		return false
	}
	if ppu.windowLayer[layer].window1Enabled && !ppu.windowLayer[layer].window2Enabled {
		var test bool = x >= int(ppu.window1Left) && x <= int(ppu.window1Right)
		if ppu.windowLayer[layer].window1Inversed {
			return !test
		} else {
			return test
		}
	}
	if !ppu.windowLayer[layer].window1Enabled && ppu.windowLayer[layer].window2Enabled {
		var test bool = x >= int(ppu.window2Left) && x <= int(ppu.window2Right)
		if ppu.windowLayer[layer].window2Inversed {
			return !test
		} else {
			return test
		}
	}
	var test1 bool = x >= int(ppu.window1Left) && x <= int(ppu.window1Right)
	var test2 bool = x >= int(ppu.window2Left) && x <= int(ppu.window2Right)
	if ppu.windowLayer[layer].window1Inversed {
		test1 = !test1
	}
	if ppu.windowLayer[layer].window2Inversed {
		test2 = !test2
	}
	switch ppu.windowLayer[layer].maskLogic {
	case 0:
		return test1 || test2
	case 1:
		return test1 && test2
	case 2:
		return test1 != test2
	case 3:
		return test1 == test2
	}
	return false
}

func (ppu *PPU) evaluateSprites(line int) {
	// TODO: iterate over oam normally to determine in-range sprites,
	//   then iterate those in-range sprites in reverse for tile-fetching
	// TODO: rectangular sprites, wierdness with sprites at -256
	var index byte
	if ppu.objPriority {
		index = (ppu.oamAddr & 0xfe)
	} else {
		index = 0
	}

	var spritesFound int = 0
	var tilesFound int = 0
	for i := 0; i < 128; i++ {
		var y byte = byte(ppu.oam[index] >> 8)
		// check if the sprite is on this line and get the sprite size
		// WARNING: convert to signed int after calculated by unsinged int.
		var row byte = byte(line) - y
		var spriteSize int = spriteSizes[ppu.objSize][(ppu.highOAM[index>>3]>>((index&7)+1))&1]
		var spriteHeight int
		if ppu.objInterlace {
			spriteHeight = spriteSize / 2
		} else {
			spriteHeight = spriteSize
		}
		if int(row) < spriteHeight {
			// in y-range, get the x location, using the high bit as well
			var x int = int(ppu.oam[index] & 0xff)
			x |= ((int(ppu.highOAM[index>>3]) >> (int(index) & 7)) & 1) << 8
			if x > 255 {
				x -= 512
			}
			// if in x-range
			if x > -spriteSize {
				// break if we found 32 sprites already
				spritesFound++
				if spritesFound > 32 {
					ppu.rangeOver = true
					// break
				}
				// update row according to obj-interlace
				if ppu.objInterlace {
					if ppu.evenFrame {
						row = row*2 + 1
					} else {
						row = row*2 + 0
					}
				}
				// get some data for the sprite and y-flip row if needed
				var tile int = int(ppu.oam[index+1] & 0xff)
				var palette int = (int(ppu.oam[index+1]) & 0xe00) >> 9
				var hFlipped bool = (ppu.oam[index+1] & 0x4000) > 0
				if (ppu.oam[index+1] & 0x8000) > 0 {
					row = byte(spriteSize) - 1 - row
				}
				// fetch all tiles in x-range
				for col := 0; col < spriteSize; col += 8 {
					if col+x > -8 && col+x < 256 {
						// break if we found 34 8*1 slivers already
						tilesFound++
						if tilesFound > 34 {
							// XXX: must be break??
							ppu.timeOver = true
							// break
						}
						// figure out which tile this uses, looping within 16x16 pages, and get it's data
						var usedCol int
						if hFlipped {
							usedCol = spriteSize - 1 - col
						} else {
							usedCol = col
						}
						var usedTile byte = byte((((tile >> 4) + (int(row) / 8)) << 4) | (((tile & 0xf) + (usedCol / 8)) & 0xF))
						var objAdr uint16
						if (ppu.oam[index+1] & 0x100) > 0 {
							objAdr = ppu.objTileAddr2
						} else {
							objAdr = ppu.objTileAddr1
						}
						var plane1 uint16 = ppu.vram[(objAdr+uint16(usedTile)*16+(uint16(row)&0x7))&0x7fff]
						var plane2 uint16 = ppu.vram[(objAdr+uint16(usedTile)*16+8+(uint16(row)&0x7))&0x7fff]
						// go over each pixel
						for px := 0; px < 8; px++ {
							var shift int
							if hFlipped {
								shift = px
							} else {
								shift = 7 - px
							}
							var pixel int = (int(plane1) >> shift) & 1
							pixel |= ((int(plane1) >> (8 + shift)) & 1) << 1
							pixel |= ((int(plane2) >> shift) & 1) << 2
							pixel |= ((int(plane2) >> (8 + shift)) & 1) << 3
							// draw it in the buffer if there is a pixel here, and the buffer there is still empty
							var screenCol int = col + x + px
							if pixel > 0 && screenCol >= 0 && screenCol < 256 && ppu.objPixelBuffer[screenCol] == 0 {
								ppu.objPixelBuffer[screenCol] = byte(0x80 + 16*palette + pixel)
								ppu.objPriorityBuffer[screenCol] = byte((ppu.oam[index+1] & 0x3000) >> 12)
							}
						}
					}
				}
				// XXX: must be break??
				// if tilesFound > 34 {
				// 	break // break out of sprite-loop if max tiles found
				// }
			}
		}
		index += 2
	}
}

func (ppu *PPU) getVRAMRemap() uint16 {
	var adr uint16 = ppu.vramPointer
	switch ppu.vramRemapMode {
	case 0:
		return adr
	case 1:
		return (adr & 0xff00) | ((adr & 0xe0) >> 5) | ((adr & 0x1f) << 3)
	case 2:
		return (adr & 0xfe00) | ((adr & 0x1c0) >> 6) | ((adr & 0x3f) << 3)
	case 3:
		return (adr & 0xfc00) | ((adr & 0x380) >> 7) | ((adr & 0x7f) << 3)
	}
	return adr
}

func (ppu *PPU) Read(addr byte) byte {
	switch addr {
	case 0x04, 0x14, 0x24,
		0x05, 0x15, 0x25,
		0x06, 0x16, 0x26,
		0x08, 0x18, 0x28,
		0x09, 0x19, 0x29,
		0x0a, 0x1a, 0x2a:
		return ppu.ppu1OpenBus
	case 0x34, 0x35, 0x36:
		var result int = int(ppu.mode7Matrix[0]) * (int(ppu.mode7Matrix[1]) >> 8)
		ppu.ppu1OpenBus = byte((result >> (8 * (int(addr) - 0x34))) & 0xff)
		return ppu.ppu1OpenBus
	case 0x37:
		// TODO: only when ppulatch is set
		ppu.hCount = ppu.console.hPos / 4
		ppu.vCount = ppu.console.vPos
		ppu.countersLatched = true
		return ppu.console.openBus
	case 0x38:
		var ret byte = 0
		if ppu.oamInHigh {
			if ppu.oamSecondWrite {
				ret = ppu.highOAM[((ppu.oamAddr&0xf)<<1)|1]
			} else {
				ret = ppu.highOAM[((ppu.oamAddr&0xf)<<1)|0]
			}
			if ppu.oamSecondWrite {
				ppu.oamAddr++
				if ppu.oamAddr == 0 {
					ppu.oamInHigh = false
				}
			}
		} else {
			if !ppu.oamSecondWrite {
				ret = byte(ppu.oam[ppu.oamAddr] & 0xff)
			} else {
				ret = byte(ppu.oam[ppu.oamAddr] >> 8)
				ppu.oamAddr++
				if ppu.oamAddr == 0 {
					ppu.oamInHigh = true
				}
			}
		}
		ppu.oamSecondWrite = !ppu.oamSecondWrite
		ppu.ppu1OpenBus = ret
		return ret
	case 0x39:
		var val uint16 = ppu.vramReadBuffer
		if !ppu.vramIncrementOnHigh {
			ppu.vramReadBuffer = ppu.vram[ppu.getVRAMRemap()&0x7fff]
			ppu.vramPointer += ppu.vramIncrement
		}
		ppu.ppu1OpenBus = byte(val & 0xff)
		return byte(val & 0xff)
	case 0x3a:
		var val uint16 = ppu.vramReadBuffer
		if ppu.vramIncrementOnHigh {
			ppu.vramReadBuffer = ppu.vram[ppu.getVRAMRemap()&0x7fff]
			ppu.vramPointer += ppu.vramIncrement
		}
		ppu.ppu1OpenBus = byte(val >> 8)
		return byte(val >> 8)
	case 0x3b:
		var ret byte = 0
		if !ppu.cgramSecondWrite {
			ret = byte(ppu.cgram[ppu.cgramPointer] & 0xff)
		} else {
			ret = byte(((ppu.cgram[ppu.cgramPointer] >> 8) & 0x7f) | (uint16(ppu.ppu2OpenBus) & 0x80))
			ppu.cgramPointer++
		}
		ppu.cgramSecondWrite = !ppu.cgramSecondWrite
		ppu.ppu2OpenBus = ret
		return ret
	case 0x3c:
		var val byte = 0
		if ppu.hCountSecond {
			val = byte(((ppu.hCount >> 8) & 1) | (uint16(ppu.ppu2OpenBus) & 0xfe))
		} else {
			val = byte(ppu.hCount & 0xff)
		}
		ppu.hCountSecond = !ppu.hCountSecond
		ppu.ppu2OpenBus = val
		return val
	case 0x3d:
		var val byte = 0
		if ppu.vCountSecond {
			val = byte(((ppu.vCount >> 8) & 1) | (uint16(ppu.ppu2OpenBus) & 0xfe))
		} else {
			val = byte(ppu.vCount & 0xff)
		}
		ppu.vCountSecond = !ppu.vCountSecond
		ppu.ppu2OpenBus = val
		return val
	case 0x3e:
		var val byte = 0x1 // ppu1 version (4 bit)
		val |= ppu.ppu1OpenBus & 0x10
		if ppu.rangeOver {
			val |= 1 << 6
		}
		if ppu.timeOver {
			val |= 1 << 7
		}
		ppu.ppu1OpenBus = val
		return val
	case 0x3f:
		var val byte = 0x3 // ppu2 version (4 bit), bit 4: ntsc/pal
		val |= ppu.ppu2OpenBus & 0x20
		if ppu.countersLatched {
			val |= 1 << 6
		}
		if ppu.evenFrame {
			val |= 1 << 7
		}
		ppu.countersLatched = false // TODO: only when ppulatch is set
		ppu.hCountSecond = false
		ppu.vCountSecond = false
		ppu.ppu2OpenBus = val
		return val
	}

	return ppu.console.openBus
}

func (ppu *PPU) Write(addr byte, value byte) {
	switch addr {
	case 0x00:
		// TODO: oam address reset when written on first line of vblank, (and when forced blank is disabled?)
		ppu.brightness = value & 0xf
		ppu.forcedBlank = (value & 0x80) > 0
	case 0x01:
		ppu.objSize = value >> 5
		ppu.objTileAddr1 = (uint16(value) & 7) << 13
		ppu.objTileAddr2 = ppu.objTileAddr1 + (((uint16(value) & 0x18) + 8) << 9)
	case 0x02:
		ppu.oamAddr = value
		ppu.oamAddrWritten = ppu.oamAddr
		ppu.oamInHigh = ppu.oamInHighWritten
		ppu.oamSecondWrite = false
	case 0x03:
		ppu.objPriority = (value & 0x80) > 0
		ppu.oamInHigh = (value & 1) > 0
		ppu.oamInHighWritten = ppu.oamInHigh
		ppu.oamAddr = ppu.oamAddrWritten
		ppu.oamSecondWrite = false
	case 0x04:
		if ppu.oamInHigh {
			if ppu.oamSecondWrite {
				ppu.highOAM[((ppu.oamAddr&0xf)<<1)|1] = value
			} else {
				ppu.highOAM[((ppu.oamAddr&0xf)<<1)|0] = value
			}
			if ppu.oamSecondWrite {
				ppu.oamAddr++
				if ppu.oamAddr == 0 {
					ppu.oamInHigh = false
				}
			}
		} else {
			if !ppu.oamSecondWrite {
				ppu.oamBuffer = value
			} else {
				ppu.oam[ppu.oamAddr] = (uint16(value) << 8) | uint16(ppu.oamBuffer)
				ppu.oamAddr++
				if ppu.oamAddr == 0 {
					ppu.oamInHigh = true
				}
			}
		}
		ppu.oamSecondWrite = !ppu.oamSecondWrite
	case 0x05:
		ppu.mode = value & 0x7
		ppu.bg3priority = (value & 0x8) > 0
		ppu.bgLayer[0].bigTiles = (value & 0x10) > 0
		ppu.bgLayer[1].bigTiles = (value & 0x20) > 0
		ppu.bgLayer[2].bigTiles = (value & 0x40) > 0
		ppu.bgLayer[3].bigTiles = (value & 0x80) > 0
	case 0x06:
		// TODO: mosaic line reset specifics
		ppu.bgLayer[0].mosaicEnabled = (value & 0x1) > 0
		ppu.bgLayer[1].mosaicEnabled = (value & 0x2) > 0
		ppu.bgLayer[2].mosaicEnabled = (value & 0x4) > 0
		ppu.bgLayer[3].mosaicEnabled = (value & 0x8) > 0
		ppu.mosaicSize = (value >> 4) + 1
		ppu.mosaicStartLine = byte(ppu.console.vPos)
	case 0x07, 0x08, 0x09, 0x0a:
		ppu.bgLayer[addr-7].tilemapWider = (value & 0x1) > 0
		ppu.bgLayer[addr-7].tilemapHigher = (value & 0x2) > 0
		ppu.bgLayer[addr-7].tilemapAddr = (uint16(value) & 0xfc) << 8
	case 0x0b:
		ppu.bgLayer[0].tileAddr = (uint16(value) & 0xf) << 12
		ppu.bgLayer[1].tileAddr = (uint16(value) & 0xf0) << 8
	case 0x0c:
		ppu.bgLayer[2].tileAddr = (uint16(value) & 0xf) << 12
		ppu.bgLayer[3].tileAddr = (uint16(value) & 0xf0) << 8
	case 0x0d:
		ppu.mode7Matrix[6] = ((int16(value) << 8) | int16(ppu.mode7Prev)) & 0x1fff
		ppu.mode7Prev = value
		// fallthrough to normal layer BG-HOFS
		fallthrough
	case 0x0f, 0x11, 0x13:
		ppu.bgLayer[(addr-0xd)/2].hScroll = ((uint16(value) << 8) | (uint16(ppu.scrollPrev) & 0xf8) | (uint16(ppu.scrollPrev2) & 0x7)) & 0x3FF
		ppu.scrollPrev = value
		ppu.scrollPrev2 = value
	case 0x0e:
		ppu.mode7Matrix[7] = ((int16(value) << 8) | int16(ppu.mode7Prev)) & 0x1fff
		ppu.mode7Prev = value
		// fallthrough to normal layer BG-VOFS
		fallthrough
	case 0x10, 0x12, 0x14:
		ppu.bgLayer[(addr-0xe)/2].vScroll = ((uint16(value) << 8) | uint16(ppu.scrollPrev)) & 0x3ff
		ppu.scrollPrev = value
	case 0x15:
		if (value & 3) == 0 {
			ppu.vramIncrement = 1
		} else if (value & 3) == 1 {
			ppu.vramIncrement = 32
		} else {
			ppu.vramIncrement = 128
		}
		ppu.vramRemapMode = (value & 0xc) >> 2
		ppu.vramIncrementOnHigh = (value & 0x80) > 0
	case 0x16:
		ppu.vramPointer = (ppu.vramPointer & 0xff00) | uint16(value)
		ppu.vramReadBuffer = ppu.vram[ppu.getVRAMRemap()&0x7fff]
	case 0x17:
		ppu.vramPointer = (ppu.vramPointer & 0x00ff) | (uint16(value) << 8)
		ppu.vramReadBuffer = ppu.vram[ppu.getVRAMRemap()&0x7fff]
	case 0x18:
		// TODO: vram access during rendering (also cgram and oam)
		var vramAdr uint16 = ppu.getVRAMRemap()
		ppu.vram[vramAdr&0x7fff] = (ppu.vram[vramAdr&0x7fff] & 0xff00) | uint16(value)
		if !ppu.vramIncrementOnHigh {
			ppu.vramPointer += ppu.vramIncrement
		}
	case 0x19:
		var vramAdr uint16 = ppu.getVRAMRemap()
		ppu.vram[vramAdr&0x7fff] = (ppu.vram[vramAdr&0x7fff] & 0x00ff) | (uint16(value) << 8)
		if ppu.vramIncrementOnHigh {
			ppu.vramPointer += ppu.vramIncrement
		}
	case 0x1a:
		ppu.mode7LargeField = (value & 0x80) > 0
		ppu.mode7CharFill = (value & 0x40) > 0
		ppu.mode7YFlip = (value & 0x2) > 0
		ppu.mode7XFlip = (value & 0x1) > 0
	case 0x1b, 0x1c, 0x1d, 0x1e:
		ppu.mode7Matrix[addr-0x1b] = (int16(value) << 8) | int16(ppu.mode7Prev)
		ppu.mode7Prev = value
	case 0x1f, 0x20:
		ppu.mode7Matrix[addr-0x1b] = ((int16(value) << 8) | int16(ppu.mode7Prev)) & 0x1fff
		ppu.mode7Prev = value
	case 0x21:
		ppu.cgramPointer = value
		ppu.cgramSecondWrite = false
	case 0x22:
		if !ppu.cgramSecondWrite {
			ppu.cgramBuffer = value
		} else {
			ppu.cgram[ppu.cgramPointer] = (uint16(value) << 8) | uint16(ppu.cgramBuffer)
			ppu.cgramPointer++
		}
		ppu.cgramSecondWrite = !ppu.cgramSecondWrite
	case 0x23, 0x24, 0x25:
		ppu.windowLayer[(addr-0x23)*2].window1Inversed = (value & 0x1) > 0
		ppu.windowLayer[(addr-0x23)*2].window1Enabled = (value & 0x2) > 0
		ppu.windowLayer[(addr-0x23)*2].window2Inversed = (value & 0x4) > 0
		ppu.windowLayer[(addr-0x23)*2].window2Enabled = (value & 0x8) > 0
		ppu.windowLayer[(addr-0x23)*2+1].window1Inversed = (value & 0x10) > 0
		ppu.windowLayer[(addr-0x23)*2+1].window1Enabled = (value & 0x20) > 0
		ppu.windowLayer[(addr-0x23)*2+1].window2Inversed = (value & 0x40) > 0
		ppu.windowLayer[(addr-0x23)*2+1].window2Enabled = (value & 0x80) > 0
	case 0x26:
		ppu.window1Left = value
	case 0x27:
		ppu.window1Right = value
	case 0x28:
		ppu.window2Left = value
	case 0x29:
		ppu.window2Right = value
	case 0x2a:
		ppu.windowLayer[0].maskLogic = value & 0x3
		ppu.windowLayer[1].maskLogic = (value >> 2) & 0x3
		ppu.windowLayer[2].maskLogic = (value >> 4) & 0x3
		ppu.windowLayer[3].maskLogic = (value >> 6) & 0x3
	case 0x2b:
		ppu.windowLayer[4].maskLogic = value & 0x3
		ppu.windowLayer[5].maskLogic = (value >> 2) & 0x3
	case 0x2c:
		ppu.layer[0].mainScreenEnabled = (value & 0x1) > 0
		ppu.layer[1].mainScreenEnabled = (value & 0x2) > 0
		ppu.layer[2].mainScreenEnabled = (value & 0x4) > 0
		ppu.layer[3].mainScreenEnabled = (value & 0x8) > 0
		ppu.layer[4].mainScreenEnabled = (value & 0x10) > 0
	case 0x2d:
		ppu.layer[0].subScreenEnabled = (value & 0x1) > 0
		ppu.layer[1].subScreenEnabled = (value & 0x2) > 0
		ppu.layer[2].subScreenEnabled = (value & 0x4) > 0
		ppu.layer[3].subScreenEnabled = (value & 0x8) > 0
		ppu.layer[4].subScreenEnabled = (value & 0x10) > 0
	case 0x2e:
		ppu.layer[0].mainScreenWindowed = (value & 0x1) > 0
		ppu.layer[1].mainScreenWindowed = (value & 0x2) > 0
		ppu.layer[2].mainScreenWindowed = (value & 0x4) > 0
		ppu.layer[3].mainScreenWindowed = (value & 0x8) > 0
		ppu.layer[4].mainScreenWindowed = (value & 0x10) > 0
	case 0x2f:
		ppu.layer[0].subScreenWindowed = (value & 0x1) > 0
		ppu.layer[1].subScreenWindowed = (value & 0x2) > 0
		ppu.layer[2].subScreenWindowed = (value & 0x4) > 0
		ppu.layer[3].subScreenWindowed = (value & 0x8) > 0
		ppu.layer[4].subScreenWindowed = (value & 0x10) > 0
	case 0x30:
		ppu.directColor = (value & 0x1) > 0
		ppu.addSubscreen = (value & 0x2) > 0
		ppu.preventMathMode = (value & 0x30) >> 4
		ppu.clipMode = (value & 0xc0) >> 6
	case 0x31:
		ppu.subtractColor = (value & 0x80) > 0
		ppu.halfColor = (value & 0x40) > 0
		for i := 0; i < 6; i++ {
			ppu.mathEnabled[i] = ((value & (1 << i)) > 0)
		}
	case 0x32:
		if (value & 0x80) > 0 {
			ppu.fixedColorB = value & 0x1f
		}
		if (value & 0x40) > 0 {
			ppu.fixedColorG = value & 0x1f
		}
		if (value & 0x20) > 0 {
			ppu.fixedColorR = value & 0x1f
		}
	case 0x33:
		ppu.interlace = (value & 0x1) > 0
		ppu.objInterlace = (value & 0x2) > 0
		ppu.overscan = (value & 0x4) > 0
		ppu.pseudoHires = (value & 0x8) > 0
		ppu.mode7ExtBG = (value & 0x40) > 0
	}
}

func (ppu *PPU) putPixels(pixels []byte) {
	var maxY int
	if ppu.frameOverscan {
		maxY = 239
	} else {
		maxY = 224
	}

	for y := 0; y < maxY; y++ {
		var dest int
		if ppu.frameOverscan {
			dest = y*2 + 2
		} else {
			dest = y*2 + 16
		}
		var y1 int = y
		var y2 int = y + 239
		if !ppu.frameInterlace {
			if ppu.evenFrame {
				y1 = y + 0
			} else {
				y1 = y + 239
			}
			y2 = y1
		}

		pixelsBase := dest * 2048
		pixelBufferBase := y1 * 2048
		size := 2048
		copy(pixels[pixelsBase:pixelsBase+size], ppu.pixelBuffer[pixelBufferBase:pixelBufferBase+size])
		pixelsBase = (dest + 1) * 2048
		pixelBufferBase = y2 * 2048
		copy(pixels[pixelsBase:pixelsBase+size], ppu.pixelBuffer[pixelBufferBase:pixelBufferBase+size])
	}

	// clear top 2 lines, and following 14 and last 16 lines if not overscanning
	for i := 0; i < (2048 * 2); i++ {
		pixels[i] = 0
	}

	if !ppu.overscan {
		for i := (2 * 2048); i < (2048 * 14); i++ {
			pixels[i] = 0
		}
		for i := (464 * 2048); i < (2048 * 16); i++ {
			pixels[i] = 0
		}
	}
}
