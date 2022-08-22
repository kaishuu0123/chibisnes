package chibisnes

type DMAChannel struct {
	bAddr      byte
	aAddr      uint16
	aBank      byte
	size       uint16 // also indirect HDMA addr
	indBank    byte   // HDMA
	tableAddr  uint16 // HDMA
	repCount   byte   // HDMA
	unusedByte byte
	dmaActive  bool
	hdmaActive bool
	mode       byte
	fixed      bool
	decrement  bool
	indirect   bool // hdma
	fromB      bool
	unusedBit  bool
	doTransfer bool // hdma
	terminated bool // hdma
	offIndex   byte
}

type DMA struct {
	console   *Console
	channels  [8]DMAChannel
	hdmaTimer uint16
	dmaTimer  uint32
	dmaBusy   bool
}

var bAddrOffsets [8][4]int = [8][4]int{
	{0, 0, 0, 0},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{0, 0, 1, 1},
	{0, 1, 2, 3},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{0, 0, 1, 1},
}

var transferLength [8]int = [8]int{
	1, 2, 2, 4, 4, 4, 2, 4,
}

func NewDMA(console *Console) *DMA {
	return &DMA{
		console: console,
	}
}

func (dma *DMA) Reset() {
	for i := 0; i < len(dma.channels); i++ {
		dma.channels[i].bAddr = 0xFF
		dma.channels[i].aAddr = 0xFFFF
		dma.channels[i].aBank = 0xFF
		dma.channels[i].size = 0xFFFF
		dma.channels[i].indBank = 0xFF
		dma.channels[i].tableAddr = 0xFFFF
		dma.channels[i].repCount = 0xFF
		dma.channels[i].unusedByte = 0xFF
		dma.channels[i].dmaActive = false
		dma.channels[i].hdmaActive = false
		dma.channels[i].mode = 7
		dma.channels[i].fixed = true
		dma.channels[i].decrement = true
		dma.channels[i].indirect = true
		dma.channels[i].fromB = true
		dma.channels[i].unusedBit = true
		dma.channels[i].doTransfer = false
		dma.channels[i].terminated = false
		dma.channels[i].offIndex = 0
	}
	dma.hdmaTimer = 0
	dma.dmaTimer = 0
	dma.dmaBusy = false
}

func (dma *DMA) Read(addr uint16) byte {
	var c byte = byte((addr & 0x70) >> 4)
	switch addr & 0xf {
	case 0x0:
		var val byte = dma.channels[c].mode
		if dma.channels[c].fixed {
			val |= 1 << 3
		}
		if dma.channels[c].decrement {
			val |= 1 << 4
		}
		if dma.channels[c].unusedBit {
			val |= 1 << 5
		}
		if dma.channels[c].indirect {
			val |= 1 << 6
		}
		if dma.channels[c].fromB {
			val |= 1 << 7
		}
		return val
	case 0x1:
		return dma.channels[c].bAddr
	case 0x2:
		return byte(dma.channels[c].aAddr & 0xff)
	case 0x3:
		return byte(dma.channels[c].aAddr >> 8)
	case 0x4:
		return dma.channels[c].aBank
	case 0x5:
		return byte(dma.channels[c].size & 0xff)
	case 0x6:
		return byte(dma.channels[c].size >> 8)
	case 0x7:
		return dma.channels[c].indBank
	case 0x8:
		return byte(dma.channels[c].tableAddr & 0xff)
	case 0x9:
		return byte(dma.channels[c].tableAddr >> 8)
	case 0xa:
		return dma.channels[c].repCount
	case 0xb, 0xf:
		return dma.channels[c].unusedByte
	default:
		return dma.console.openBus
	}
}

func (dma *DMA) Write(addr uint16, value byte) {
	var c byte = byte((addr & 0x70) >> 4)
	switch addr & 0xf {
	case 0x0:
		dma.channels[c].mode = value & 0x7
		dma.channels[c].fixed = (value & 0x8) > 0
		dma.channels[c].decrement = (value & 0x10) > 0
		dma.channels[c].unusedBit = (value & 0x20) > 0
		dma.channels[c].indirect = (value & 0x40) > 0
		dma.channels[c].fromB = (value & 0x80) > 0
	case 0x1:
		dma.channels[c].bAddr = value
	case 0x2:
		dma.channels[c].aAddr = (dma.channels[c].aAddr & 0xff00) | uint16(value)
	case 0x3:
		dma.channels[c].aAddr = (dma.channels[c].aAddr & 0xff) | (uint16(value) << 8)
	case 0x4:
		dma.channels[c].aBank = value
	case 0x5:
		dma.channels[c].size = (dma.channels[c].size & 0xff00) | uint16(value)
	case 0x6:
		dma.channels[c].size = (dma.channels[c].size & 0xff) | (uint16(value) << 8)
	case 0x7:
		dma.channels[c].indBank = value
	case 0x8:
		dma.channels[c].tableAddr = (dma.channels[c].tableAddr & 0xff00) | uint16(value)
	case 0x9:
		dma.channels[c].tableAddr = (dma.channels[c].tableAddr & 0xff) | (uint16(value) << 8)
	case 0xa:
		dma.channels[c].repCount = value
	case 0xb, 0xf:
		dma.channels[c].unusedByte = value
	default:
		break
	}
}

func (dma *DMA) doDMA() {
	if dma.dmaTimer > 0 {
		dma.dmaTimer -= 2
		return
	}
	// figure out first channel that is active
	var i int = 0
	for i = 0; i < len(dma.channels); i++ {
		if dma.channels[i].dmaActive {
			break
		}
	}
	if i == 8 {
		// no active channels
		dma.dmaBusy = false
		return
	}
	// do channel i
	dma.transferByte(
		dma.channels[i].aAddr, dma.channels[i].aBank,
		dma.channels[i].bAddr+byte(bAddrOffsets[dma.channels[i].mode][dma.channels[i].offIndex]),
		dma.channels[i].fromB)
	dma.channels[i].offIndex++
	dma.channels[i].offIndex &= 3

	dma.dmaTimer += 6 // 8 cycles for each byte taken, -2 for this cycle
	if !dma.channels[i].fixed {
		if dma.channels[i].decrement {
			dma.channels[i].aAddr -= 1
		} else {
			dma.channels[i].aAddr += 1
		}
	}

	dma.channels[i].size--
	if dma.channels[i].size == 0 {
		dma.channels[i].offIndex = 0 // reset offset index
		dma.channels[i].dmaActive = false
		dma.dmaTimer += 8 // 8 cycle overhead per channel
	}
}

func (dma *DMA) initHDMA() {
	dma.hdmaTimer = 0
	var hdmaHappened bool = false
	for i := 0; i < len(dma.channels); i++ {
		if dma.channels[i].hdmaActive {
			hdmaHappened = true
			// terminate any dma
			dma.channels[i].dmaActive = false
			dma.channels[i].offIndex = 0
			// load address, repCount, and indirect address if needed
			dma.channels[i].tableAddr = dma.channels[i].aAddr
			dma.channels[i].repCount = dma.console.Read((uint32(dma.channels[i].aBank) << 16) | uint32(dma.channels[i].tableAddr))
			dma.channels[i].tableAddr++
			dma.hdmaTimer += 8 // 8 cycle overhead for each active channel
			if dma.channels[i].indirect {
				dma.channels[i].size = uint16(dma.console.Read((uint32(dma.channels[i].aBank) << 16) | uint32(dma.channels[i].tableAddr)))
				dma.channels[i].tableAddr++
				dma.channels[i].size |= uint16(dma.console.Read((uint32(dma.channels[i].aBank)<<16)|uint32(dma.channels[i].tableAddr))) << 8
				dma.channels[i].tableAddr++
				dma.hdmaTimer += 16 // another 16 cycles for indirect (total 24)
			}
			dma.channels[i].doTransfer = true
		} else {
			dma.channels[i].doTransfer = false
		}
		dma.channels[i].terminated = false
	}
	if hdmaHappened {
		dma.hdmaTimer += 16 // 18 cycles overhead, -2 for this cycle
	}
}

func (dma *DMA) doHDMA() {
	dma.hdmaTimer = 0
	var hdmaHappened bool = false
	for i := 0; i < len(dma.channels); i++ {
		if dma.channels[i].hdmaActive && !dma.channels[i].terminated {
			hdmaHappened = true
			// terminate any dma
			dma.channels[i].dmaActive = false
			dma.channels[i].offIndex = 0
			// do the hdma
			dma.hdmaTimer += 8 // 8 cycles overhead for each active channel
			if dma.channels[i].doTransfer {
				for j := 0; j < transferLength[dma.channels[i].mode]; j++ {
					dma.hdmaTimer += 8 // 8 cycles for each byte transferred
					if dma.channels[i].indirect {
						dma.transferByte(
							dma.channels[i].size, dma.channels[i].indBank,
							dma.channels[i].bAddr+byte(bAddrOffsets[dma.channels[i].mode][j]),
							dma.channels[i].fromB)
						dma.channels[i].size++
					} else {
						dma.transferByte(
							dma.channels[i].tableAddr, dma.channels[i].aBank,
							dma.channels[i].bAddr+byte(bAddrOffsets[dma.channels[i].mode][j]),
							dma.channels[i].fromB)
						dma.channels[i].tableAddr++
					}
				}
			}
			dma.channels[i].repCount--
			dma.channels[i].doTransfer = (dma.channels[i].repCount & 0x80) > 0
			if (dma.channels[i].repCount & 0x7f) == 0 {
				dma.channels[i].repCount = dma.console.Read((uint32(dma.channels[i].aBank) << 16) | uint32(dma.channels[i].tableAddr))
				dma.channels[i].tableAddr++
				if dma.channels[i].indirect {
					// TODO: oddness with not fetching high byte if last active channel and reCount is 0
					dma.channels[i].size = uint16(dma.console.Read((uint32(dma.channels[i].aBank) << 16) | uint32(dma.channels[i].tableAddr)))
					dma.channels[i].tableAddr++
					dma.channels[i].size |= uint16(dma.console.Read((uint32(dma.channels[i].aBank)<<16)|uint32(dma.channels[i].tableAddr))) << 8
					dma.channels[i].tableAddr++
					dma.hdmaTimer += 16 // 16 cycles for new indirect address
				}
				if dma.channels[i].repCount == 0 {
					dma.channels[i].terminated = true
				}
				dma.channels[i].doTransfer = true
			}
		}
	}
	if hdmaHappened {
		dma.hdmaTimer += 16 // 18 cycles overhead, -2 for this cycle
	}
}

func (dma *DMA) transferByte(aAddr uint16, aBank byte, bAddr byte, fromB bool) {
	// TODO: invalid writes:
	//   accesing b-bus via a-bus gives open bus,
	//   $2180-$2183 while accessing ram via a-bus open busses $2180-$2183
	//   cannot access $4300-$437f (dma regs), or $420b / $420c
	if fromB {
		dma.console.Write((uint32(aBank)<<16)|uint32(aAddr), dma.console.ReadBBus(bAddr))
	} else {
		dma.console.WriteBBus(bAddr, dma.console.Read((uint32(aBank)<<16)|uint32(aAddr)))
	}
}

func (dma *DMA) Cycle() bool {
	if dma.hdmaTimer > 0 {
		dma.hdmaTimer -= 2
		return true
	} else if dma.dmaBusy {
		dma.doDMA()
		return true
	}
	return false
}

func (dma *DMA) StartDMA(value byte, hdma bool) {
	for i := 0; i < len(dma.channels); i++ {
		if hdma {
			dma.channels[i].hdmaActive = (value & (1 << i)) > 0
		} else {
			dma.channels[i].dmaActive = (value & (1 << i)) > 0
		}
	}
	if !hdma {
		dma.dmaBusy = (value > 0)
		if dma.dmaBusy {
			dma.dmaTimer += 16
		} else {
			dma.dmaTimer += 0
		}
	}
}
