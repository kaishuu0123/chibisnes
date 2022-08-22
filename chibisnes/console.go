package chibisnes

import (
	"errors"
	"fmt"
	"log"
)

const apuCyclesPerMaster float64 = (32040.0 * 32.0) / (1364.0 * 262.0 * 60.0)

type Console struct {
	CPU       *CPU
	PPU       *PPU
	APU       *APU
	DMA       *DMA
	Cartridge *Cartridge

	Controller1 *Controller
	Controller2 *Controller

	RAM     [0x20000]byte
	RAMAddr uint32

	hPos   uint16
	vPos   uint16
	frames uint32

	cpuCyclesLeft    byte
	cpuMemOps        byte
	apuCatchupCycles float64

	hIRQEnabled bool
	vIRQEnabled bool
	nmiEnabled  bool
	hTimer      uint16
	vTimer      uint16

	inNMI    bool
	inIRQ    bool
	inVBlank bool

	// joypad handling
	portAutoRead [4]uint16 // as read by auto-joypad read
	autoJoyRead  bool
	autoJoyTimer uint16

	ppuLatch bool

	multiplyA      byte
	multiplyResult uint16
	divideA        uint16
	divideResult   uint16

	fastMem bool
	openBus byte

	Debug       bool
	RomFilePath string
}

func NewConsole() *Console {
	c := &Console{}
	c.CPU = NewCPU(c)
	c.PPU = NewPPU(c)
	c.APU = NewAPU(c)
	c.DMA = NewDMA(c)
	c.Cartridge = NewCartridge(c)
	c.Controller1 = NewController(c)
	c.Controller2 = NewController(c)

	return c
}

func (console *Console) Reset(hard bool) {
	console.Cartridge.Reset()
	console.CPU.Reset()
	console.PPU.Reset()
	console.APU.Reset()
	console.DMA.Reset()
	console.Controller1.Reset()
	console.Controller2.Reset()
	if hard {
		for i := 0; i < len(console.RAM); i++ {
			console.RAM[i] = 0
		}
	}
	console.RAMAddr = 0
	console.hPos = 0
	console.vPos = 0
	console.frames = 0
	console.cpuCyclesLeft = 52 // 5 reads (8) + 2 IntOp (6)
	console.cpuMemOps = 0
	console.apuCatchupCycles = 0.0
	console.hIRQEnabled = false
	console.vIRQEnabled = false
	console.nmiEnabled = false
	console.hTimer = 0x1FF
	console.vTimer = 0x1FF
	console.inNMI = false
	console.inIRQ = false
	console.inVBlank = false
	console.ppuLatch = false
	console.multiplyA = 0xFF
	console.multiplyResult = 0xFE01
	console.divideA = 0xFFFF
	console.divideResult = 0x0101
	console.fastMem = false
	console.openBus = 0
}

func (console *Console) CPURead(addr uint32) byte {
	console.cpuMemOps++
	console.cpuCyclesLeft += byte(console.getAccessTime(addr))
	return console.Read(addr)
}

func (console *Console) CPUWrite(addr uint32, value byte) {
	console.cpuMemOps++
	console.cpuCyclesLeft += byte(console.getAccessTime(addr))
	console.Write(addr, value)
}

func (console *Console) getAccessTime(addr uint32) int {
	var bank byte = byte(addr >> 16)
	addr &= 0xFFFF
	switch {
	case bank >= 0x40 && bank < 0x80:
		return 8
	case bank >= 0xC0:
		if console.fastMem {
			return 6
		} else {
			return 8
		}
	case addr < 0x2000:
		return 8
	case addr < 0x4000:
		return 6
	case addr < 0x4200:
		return 12
	case addr < 0x6000:
		return 6
	case addr < 0x8000:
		return 8
	}

	if console.fastMem && bank >= 0x80 {
		return 6
	}

	return 8
}

func (console *Console) Read(addr uint32) byte {
	var value byte = console.RRead(addr)
	console.openBus = value
	return value
}

func (console *Console) RRead(addr uint32) byte {
	var bank byte = byte(addr >> 16)
	addr &= 0xFFFF

	if bank == 0x7E || bank == 0x7F {
		return console.RAM[((uint32(bank)&1)<<16)|addr]
	}

	if bank < 0x40 || (bank >= 0x80 && bank < 0xC0) {
		switch {
		case addr < 0x2000:
			return console.RAM[addr]
		case addr >= 0x2100 && addr < 0x2200:
			return console.ReadBBus(byte(addr & 0xFF))
		case addr == 0x4016:
			return console.Controller1.Read() | (console.openBus & 0xFC)
		case addr == 0x4017:
			return console.Controller2.Read() | (console.openBus & 0xE0) | 0x1C
		case addr >= 0x4200 && addr < 0x4220:
			return console.ReadReg(uint16(addr))
		case addr >= 0x4300 && addr < 0x4380:
			return console.DMA.Read(uint16(addr))
		}
	}

	return console.Cartridge.Read(bank, uint16(addr))
}

func (console *Console) Write(addr uint32, value byte) {
	console.openBus = value

	var bank byte = byte(addr >> 16)
	addr &= 0xFFFF

	if bank == 0x7E || bank == 0x7F {
		console.RAM[((uint32(bank)&1)<<16)|addr] = value
	}

	if bank < 0x40 || (bank >= 0x80 && bank < 0xC0) {
		switch {
		case addr < 0x2000:
			console.RAM[addr] = value
		case addr >= 0x2100 && addr < 0x2200:
			console.WriteBBus(byte(addr&0xFF), value)
		case addr == 0x4016:
			console.Controller1.latchLine = (value & 0x01) > 0
			console.Controller2.latchLine = (value & 0x01) > 0
		case addr >= 0x4200 && addr < 0x4220:
			console.WriteReg(uint16(addr), value)
		case addr >= 0x4300 && addr < 0x4380:
			console.DMA.Write(uint16(addr), value)
		}
	}

	console.Cartridge.Write(bank, uint16(addr), value)
}

func (console *Console) ReadBBus(addr byte) byte {
	switch {
	case addr < 0x40:
		return console.PPU.Read(addr)
	case addr < 0x80:
		console.catchupAPU()
		return console.APU.outPorts[addr&0x03]
	case addr == 0x80:
		var ret uint8 = console.RAM[console.RAMAddr]
		console.RAMAddr++
		console.RAMAddr &= 0x1FFFF
		return ret
	}

	return console.openBus
}

func (console *Console) WriteBBus(addr byte, value byte) {
	switch {
	case addr < 0x40:
		console.PPU.Write(addr, value)
		return
	case addr < 0x80:
		console.catchupAPU()
		console.APU.inPorts[addr&0x3] = value
		return
	}

	switch addr {
	case 0x80:
		console.RAM[console.RAMAddr] = value
		console.RAMAddr++
		console.RAMAddr &= 0x1FFFF
	case 0x81:
		console.RAMAddr = (console.RAMAddr & 0x1FF00) | uint32(value)
	case 0x82:
		console.RAMAddr = (console.RAMAddr & 0x100FF) | (uint32(value) << 8)
	case 0x83:
		console.RAMAddr = (console.RAMAddr & 0x0FFFF) | ((uint32(value) & 1) << 16)
	}
}

func (console *Console) ReadReg(addr uint16) byte {
	switch addr {
	case 0x4210:
		var val byte = 0x2 // CPU version (4 bit)
		if console.inNMI {
			val |= 1 << 7
		}

		// XXX: really fix?
		//   ex.) Chrono Trigger
		if console.CPU.nmiWanted {
			// nothing done
		} else {
			console.inNMI = false
		}
		return val | (console.openBus & 0x70)
	case 0x4211:
		var val byte
		if console.inIRQ {
			val |= 1 << 7
		}
		console.inIRQ = false
		console.CPU.irqWanted = false
		return val | (console.openBus & 0x7f)
	case 0x4212:
		var val byte
		if console.autoJoyTimer > 0 {
			val = 0x01
		} else {
			val = 0x00
		}
		if console.hPos >= 1024 {
			val |= 1 << 6
		}
		if console.inVBlank {
			val |= 1 << 7
		}
		return val | (console.openBus & 0x3e)
	case 0x4213:
		// IO-port
		if console.ppuLatch {
			return 1 << 7
		}
		return 0
	case 0x4214:
		return byte(console.divideResult & 0xFF)
	case 0x4215:
		return byte(console.divideResult >> 8)
	case 0x4216:
		return byte(console.multiplyResult & 0xFF)
	case 0x4217:
		return byte(console.multiplyResult >> 8)
	case 0x4218, 0x421a, 0x421c, 0x421e:
		return byte(console.portAutoRead[(addr-0x4218)/2] & 0xff)
	case 0x4219, 0x421b, 0x421d, 0x421f:
		return byte(console.portAutoRead[(addr-0x4219)/2] >> 8)
	}

	return console.openBus
}

func (console *Console) WriteReg(addr uint16, value byte) {
	switch addr {
	case 0x4200:
		console.autoJoyRead = (value & 0x1) > 0
		if !console.autoJoyRead {
			console.autoJoyTimer = 0
		}
		if (value & 0x10) > 0 {
			console.hIRQEnabled = true
		} else {
			console.hIRQEnabled = false
		}
		if (value & 0x20) > 0 {
			console.vIRQEnabled = true
		} else {
			console.vIRQEnabled = false
		}
		if (value & 0x80) > 0 {
			console.nmiEnabled = true
		} else {
			console.nmiEnabled = false
		}
		if !console.hIRQEnabled && !console.vIRQEnabled {
			console.inIRQ = false
			console.CPU.irqWanted = false
		}
		// TODO: enabling nmi during vblank with inNmi still set generates nmi
		//   enabling virq (and not h) on the vPos that vTimer is at generates irq (?)
	case 0x4201:
		if !((value & 0x80) > 0) && console.ppuLatch {
			// latch the ppu
			console.PPU.Read(0x37)
		}
		if (value & 0x80) > 0 {
			console.ppuLatch = true
		} else {
			console.ppuLatch = false
		}
	case 0x4202:
		console.multiplyA = value
	case 0x4203:
		console.multiplyResult = uint16(console.multiplyA) * uint16(value)
	case 0x4204:
		console.divideA = (uint16(console.divideA) & 0xff00) | uint16(value)
	case 0x4205:
		console.divideA = (uint16(console.divideA) & 0x00ff) | (uint16(value) << 8)
	case 0x4206:
		if value == 0 {
			console.divideResult = 0xffff
			console.multiplyResult = console.divideA
		} else {
			console.divideResult = uint16(console.divideA) / uint16(value)
			console.multiplyResult = uint16(console.divideA) % uint16(value)
		}
	case 0x4207:
		console.hTimer = (console.hTimer & 0x100) | uint16(value)
	case 0x4208:
		console.hTimer = (console.hTimer & 0x0ff) | ((uint16(value) & 1) << 8)
	case 0x4209:
		console.vTimer = (console.vTimer & 0x100) | uint16(value)
	case 0x420a:
		console.vTimer = (console.vTimer & 0x0ff) | ((uint16(value) & 1) << 8)
	case 0x420b:
		console.DMA.StartDMA(value, false)
	case 0x420c:
		console.DMA.StartDMA(value, true)
	case 0x420d:
		if (value & 0x1) > 0 {
			console.fastMem = true
		} else {
			console.fastMem = false
		}
	}
}

func (console *Console) catchupAPU() {
	var catchupCycles int = int(console.apuCatchupCycles)
	for i := 0; i < catchupCycles; i++ {
		console.APU.Cycle()
	}
	console.apuCatchupCycles -= float64(catchupCycles)
}

func (console *Console) RunFrame() {
	console.runCycle()
	for !(console.hPos == 0 && console.vPos == 0) {
		console.runCycle()
	}
}

func (console *Console) runCycle() {
	console.apuCatchupCycles += apuCyclesPerMaster * 2.0
	console.Controller1.Cycle()
	console.Controller2.Cycle()
	// if not in dram refresh, if we are busy with hdma/dma, do that, else do cpu cycle
	if console.hPos < 536 || console.hPos >= 576 {
		if !console.DMA.Cycle() {
			console.runCPU()
		}
	}

	// check for h/v timer irq's
	if console.vIRQEnabled && console.hIRQEnabled {
		if console.vPos == console.vTimer && console.hPos == (4*console.hTimer) {
			console.inIRQ = true
			console.CPU.irqWanted = true // request IRQ on CPU
		}
	} else if console.vIRQEnabled && !console.hIRQEnabled {
		if console.vPos == console.vTimer && console.hPos == 0 {
			console.inIRQ = true
			console.CPU.irqWanted = true // request IRQ on CPU
		}
	} else if !console.vIRQEnabled && console.hIRQEnabled {
		if console.hPos == (4 * console.hTimer) {
			console.inIRQ = true
			console.CPU.irqWanted = true // request IRQ on CPU
		}
	}

	// handle positional stuff
	if console.hPos == 512 {
		// render the line halfway of the screen for better compatibility
		if !console.inVBlank {
			console.PPU.runLine(int(console.vPos))
		}
	} else if console.hPos == 1024 {
		// start of hblank
		if !console.inVBlank {
			console.DMA.doHDMA()
		}
	}

	// handle autoJoyRead-timer
	if console.autoJoyTimer > 0 {
		console.autoJoyTimer -= 2
	}

	// increment position
	// exact frame timing line 240 on odd frame is 4 cycles shorter. (1360)
	console.hPos += 2
	if console.hPos == 1364 || (!console.PPU.interlace && !console.PPU.evenFrame && console.vPos == 240 && console.hPos == 1360) {
		console.hPos = 0
		console.vPos++

		var endVPos uint16
		// even frames in interlace is 1 extra line
		if console.PPU.interlace && console.PPU.evenFrame {
			endVPos = 262
		} else {
			endVPos = 261
		}

		if console.vPos == (endVPos + 1) {
			console.vPos = 0
			console.frames++
			console.catchupAPU() // catch up the apu at the end of the frame
		}
	}

	// TODO: better timing? (especially Hpos)
	if console.hPos == 0 {
		// end of hblank, do most vPos-tests
		var startingVblank bool = false
		if console.vPos == 0 {
			// end of vblank
			console.inVBlank = false
			console.inNMI = false
			console.DMA.initHDMA()
		} else if console.vPos == 225 {
			// ask the ppu if we start vblank now or at vPos 240 (overscan)
			startingVblank = !console.PPU.checkOverscan()
		} else if console.vPos == 240 {
			// if we are not yet in vblank, we had an overscan frame, set startingVblank
			if !console.inVBlank {
				startingVblank = console.PPU.checkOverscan()
			}
		}

		if startingVblank {
			// if we are starting vblank
			console.PPU.handleVBlank()
			console.inVBlank = true
			console.inNMI = true
			if console.autoJoyRead {
				// TODO: this starts a little after start of vblank
				console.autoJoyTimer = 4224
				console.doAutoJoypad()
			}
			if console.nmiEnabled {
				console.CPU.nmiWanted = true // request NMI on CPU
			}
		}
	}

	if console.Debug {
		console.debugPrint()
	}
}

func (console *Console) runCPU() {
	if console.cpuCyclesLeft == 0 {
		console.cpuMemOps = 0
		var cycles int = console.CPU.runOpcode()
		console.CPU.cycleCounter += uint64(cycles)
		console.cpuCyclesLeft += byte((cycles - int(console.cpuMemOps)) * 6)

	}
	console.cpuCyclesLeft -= 2
}

func (console *Console) LoadROM(romFilePath string, data []byte, dataLen int) error {
	console.RomFilePath = romFilePath

	if len(data) < 0x8000 {
		// if smaller than smallest possible, don't load
		msg := fmt.Sprintf("Failed to load rom: rom to small (%d bytes)\n", len(data))
		return errors.New(msg)
	}

	// check headers
	var headers [4]CartridgeHeader
	for i := 0; i < len(headers); i++ {
		headers[i].score = -50
	}
	if dataLen >= 0x8000 {
		console.readHeader(data, 0x7fc0, &headers[0])
	}
	if dataLen >= 0x8200 {
		console.readHeader(data, 0x81c0, &headers[1])
	}
	if dataLen >= 0x10000 {
		console.readHeader(data, 0xffc0, &headers[2])
	}
	if dataLen >= 0x10200 {
		console.readHeader(data, 0x101c0, &headers[3])
	}
	// see which it is
	var max int = 0
	var used int = 0
	for i := 0; i < len(headers); i++ {
		if int(headers[i].score) > max {
			max = int(headers[i].score)
			used = i
		}
	}
	if (used & 1) > 0 {
		// odd-numbered ones are for headered roms
		// data += 0x200    // move pointer past header
		copy(data[:], data[0x200:])
		dataLen -= 0x200 // and subtract from size
	}
	// check if we can load it
	if headers[used].cartType > 2 {
		msg := fmt.Sprintf("Failed to load rom: unsupported type (%d)\n", headers[used].cartType)
		return errors.New(msg)
	}
	// expand to a power of 2
	var newLength int = 0x8000
	for true {
		if dataLen <= newLength {
			break
		}
		newLength *= 2
	}
	var newData []byte = make([]byte, newLength)
	for i := 0; i < dataLen; i++ {
		newData[i] = data[i]
	}
	var test int = 1
	for dataLen != newLength {
		if (dataLen & test) > 0 {
			copy(newData[dataLen:dataLen+test], newData[dataLen-test:(dataLen-test)+test])
			dataLen += test
		}
		test *= 2
	}

	// load it
	if headers[used].cartType == 2 {
		log.Printf("ROM: Loaded %s rom \"%s\"\n", "HiROM", headers[used].name)
	} else {
		log.Printf("ROM: Loaded %s rom \"%s\"\n", "LoROM", headers[used].name)
	}

	var ramSize int
	if headers[used].chips > 0 {
		ramSize = int(headers[used].ramSize)
	} else {
		ramSize = 0
	}
	console.Cartridge.Load(int(headers[used].cartType), newData, newLength, ramSize, headers[used].coprocessor)

	log.Printf("ROM: Coprocessor Type: %d\n", headers[used].coprocessor)

	console.Reset(true) // reset after loading

	return nil
}

func (console *Console) readHeader(data []byte, offset int, header *CartridgeHeader) {
	// read name, TODO: non-ASCII names?
	for i := 0; i < 21; i++ {
		var ch byte = data[offset+i]
		if ch >= 0x20 && ch < 0x7f {
			header.name[i] = ch
		} else {
			header.name[i] = '.'
		}
	}
	header.name[21] = 0
	// read rest
	header.speed = data[offset+0x15] >> 4
	header.mode = data[offset+0x15] & 0xf
	header.coprocessor = data[offset+0x16] >> 4
	header.chips = data[offset+0x16] & 0xf
	header.romSize = 0x400 << data[offset+0x17]
	header.ramSize = 0x400 << data[offset+0x18]
	header.region = data[offset+0x19]
	header.maker = data[offset+0x1a]
	header.version = data[offset+0x1b]
	header.checksumComplement = (uint16(data[offset+0x1d]) << 8) + uint16(data[offset+0x1c])
	header.checksum = (uint16(data[offset+0x1f]) << 8) + uint16(data[offset+0x1e])

	// read v3 and/or v2
	header.headerVersion = 1
	if header.maker == 0x33 {
		header.headerVersion = 3
		// maker code
		for i := 0; i < 2; i++ {
			var ch byte = data[offset-0x10+i]
			if ch >= 0x20 && ch < 0x7f {
				header.makerCode[i] = ch
			} else {
				header.makerCode[i] = '.'
			}
		}
		header.makerCode[2] = 0
		// game code
		for i := 0; i < 4; i++ {
			base := offset - 0xE + i
			var ch byte = data[base]
			if ch >= 0x20 && ch < 0x7f {
				header.gameCode[i] = ch
			} else {
				header.gameCode[i] = '.'
			}
		}
		header.gameCode[4] = 0
		header.flashSize = 0x400 << data[offset-4]
		header.exRamSize = 0x400 << data[offset-3]
		header.specialVersion = data[offset-2]
		header.exCoprocessor = data[offset-1]
	} else if data[offset+0x14] == 0 {
		header.headerVersion = 2
		header.exCoprocessor = data[offset-1]
	}

	// get region
	header.pal = (header.region >= 0x2 && header.region <= 0xc) || header.region == 0x11

	// set cartType
	if offset < 0x9000 {
		header.cartType = 1
	} else {
		header.cartType = 2
	}

	// get score
	header.score = int16(console.calculateScore(data, offset, header))
}

func (console *Console) calculateScore(data []byte, offset int, header *CartridgeHeader) int {
	// get score
	// TODO: check name, maker/game-codes (if V3) for ASCII, more vectors,
	//   more first opcode, rom-sizes (matches?), type (matches header offset?)
	var score int = 0
	if header.speed == 2 || header.speed == 3 {
		score += 5
	} else {
		score += -4
	}
	if header.mode <= 3 || header.mode == 5 {
		score += 5
	} else {
		score += -2
	}
	if header.coprocessor <= 5 || header.coprocessor >= 0xe {
		score += 5
	} else {
		score += -2
	}
	if header.chips <= 6 || header.chips == 9 || header.chips == 0xa {
		score += 5
	} else {
		score += -2
	}
	if header.region <= 0x14 {
		score += 5
	} else {
		score += -2
	}
	if header.checksum+header.checksumComplement == 0xffff {
		score += 8
	} else {
		score += -6
	}
	var resetVector uint16 = uint16(data[offset+0x3c]) | (uint16(data[offset+0x3d]) << 8)
	if resetVector >= 0x8000 {
		score += 8
	} else {
		score += -20
	}
	// check first opcode after reset
	var opcode byte = data[offset+0x40-0x8000+(int(resetVector)&0x7fff)]
	if opcode == 0x78 || opcode == 0x18 {
		// sei, clc (for clc:xce)
		score += 6
	}
	if opcode == 0x4c || opcode == 0x5c || opcode == 0x9c {
		// jmp abs, jml abl, stz abs
		score += 3
	}
	if opcode == 0x00 || opcode == 0xff || opcode == 0xdb {
		// brk, sbc alx, stp
		score -= 6
	}

	return score
}

func (console *Console) doAutoJoypad() {
	// TODO: improve? (now calls input_cycle)
	for i := 0; i < len(console.portAutoRead); i++ {
		console.portAutoRead[i] = 0
	}
	console.Controller1.latchLine = true
	console.Controller2.latchLine = true
	console.Controller1.Cycle() // latches the controllers
	console.Controller2.Cycle()
	console.Controller1.latchLine = false
	console.Controller2.latchLine = false
	for i := 0; i < 16; i++ {
		var val byte = console.Controller1.Read()
		console.portAutoRead[0] |= ((uint16(val) & 1) << (15 - i))
		console.portAutoRead[2] |= (((uint16(val) >> 1) & 1) << (15 - i))
		val = console.Controller2.Read()
		console.portAutoRead[1] |= ((uint16(val) & 1) << (15 - i))
		console.portAutoRead[3] |= (((uint16(val) >> 1) & 1) << (15 - i))
	}
}

func (console *Console) SetPixels(pixelData []byte) {
	console.PPU.putPixels(pixelData)
}

func (console *Console) SetAudioSamples(sampleData []int16, samplesPerFrame int) {
	// size is 2 (int16) * 2 (stereo) * samplesPerFrame
	// sets samples in the sampleData
	console.APU.dsp.getSamples(sampleData, samplesPerFrame)
}

func (console *Console) SetButtonState(player int, button int, pressed bool) {
	// set key in constroller
	switch player {
	case 1:
		if pressed {
			console.Controller1.currentState |= 1 << button
		} else {
			console.Controller1.currentState &= ^(1 << button)
		}
	case 2:
		if pressed {
			console.Controller2.currentState |= 1 << button
		} else {
			console.Controller2.currentState &= ^(1 << button)
		}
	}
}

func (console *Console) debugPrint() {
	console.catchupAPU()
	if console.DMA.hdmaTimer > 0 || console.DMA.dmaBusy {
		// nothing done
	} else {
		if console.cpuCyclesLeft == 0 {
			fmt.Printf("%s\n", console.CPU.getProcessorStateCPU())
		}
	}
	if console.apuCatchupCycles+(apuCyclesPerMaster*2.0) >= 1.0 {
		// we will run a apu cycle next call, see if it also starts a opcode
		if console.APU.cpuCyclesLeft == 0 {
			fmt.Printf("%s\n", console.APU.spc.getProcessorStateSPC())
		}
	} else {
		// nothing done
	}
}

func (console *Console) Close() {
	console.Cartridge.Close()
}
