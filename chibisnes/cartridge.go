package chibisnes

import "path/filepath"

type CartridgeHeader struct {
	// normal header
	headerVersion byte     // 1, 2, 3
	name          [22]byte // $ffc0-$ffd4 (max 21 bytes + \0), $ffd4=$00: header V2
	speed         byte     // $ffd5.7-4 (always 2 or 3)

	mode               byte   // $ffd5.3-0
	coprocessor        byte   // $ffd6.7-4
	chips              byte   // $ffd6.3-0
	romSize            uint32 // $ffd7 (0x400 << x)
	ramSize            uint32 // $ffd8 (0x400 << x)
	region             byte   // $ffd9 (also NTSC/PAL)
	maker              byte   // $ffda ($33: header V3)
	version            byte   // $ffdb
	checksumComplement uint16 // $ffdc,$ffdd
	checksum           uint16 // $ffde,$ffdf
	// v2/v3 (v2 only exCoprocessor)
	makerCode      [3]byte // $ffb0,$ffb1: (2 chars + \0)
	gameCode       [5]byte // $ffb2-$ffb5: (4 chars + \0)
	flashSize      uint32  // $ffbc (0x400 << x)
	exRamSize      uint32  // $ffbd (0x400 << x) (used for GSU?)
	specialVersion byte    // $ffbe
	exCoprocessor  byte    // $ffbf (if coprocessor = $f)
	// calculated stuff
	score    int16 // score for header, to see which mapping is most likely
	pal      bool  // if this is a rom for PAL regions instead of NTSC
	cartType byte  // calculated type
}

type Cartridge struct {
	console *Console

	coprocessor *CPU
	cartType    byte
	rom         []byte
	romSize     uint32
	// ram      []byte
	ram     *SRAM
	ramSize uint32
}

func NewCartridge(console *Console) *Cartridge {
	return &Cartridge{
		console: console,
	}
}

func (cartridge *Cartridge) Reset() {
	// nothing done
}

func (cartridge *Cartridge) Load(cartType int, rom []byte, romSize int, ramSize int, coprocessor byte) {
	// XXX: correct? (byte cast)
	cartridge.cartType = byte(cartType)

	cartridge.rom = make([]byte, romSize)
	cartridge.romSize = uint32(romSize)
	for i := 0; i < len(rom); i++ {
		cartridge.rom[i] = rom[i]
	}

	if ramSize > 0 {
		// cartridge.ram = make([]byte, ramSize)
		// cartridge.ramSize = uint32(ramSize)
		// for i := 0; i < len(cartridge.ram); i++ {
		// 	cartridge.ram[i] = 0
		// }

		cartridge.ramSize = uint32(ramSize)

		filePath := cartridge.console.RomFilePath
		ramFileName := getFileNameWithoutExtension(filePath)
		ramFileDir := filepath.Dir(filepath.Clean(filePath))
		ramFilePath := filepath.Join(ramFileDir, ramFileName+`.srm`)

		cartridge.ram = NewSRAM(ramFilePath, ramSize)
	} else {
		cartridge.ram = nil
		cartridge.ramSize = 0
	}

	if coprocessor == 3 {
		cartridge.coprocessor = NewCPU(cartridge.console)
	}
}

func (cartridge *Cartridge) Read(bank byte, addr uint16) byte {
	switch cartridge.cartType {
	case 0:
		return cartridge.console.openBus
	case 1:
		return cartridge.readLoROM(bank, addr)
	case 2:
		return cartridge.readHiROM(bank, addr)
	}

	return cartridge.console.openBus
}

func (cartridge *Cartridge) Write(bank byte, addr uint16, value byte) {
	switch cartridge.cartType {
	case 0:
		// nothing done
	case 1:
		cartridge.writeLoROM(bank, addr, value)
	case 2:
		cartridge.writeHiROM(bank, addr, value)
	}
}

func (cartridge *Cartridge) readLoROM(bank byte, addr uint16) byte {
	if ((bank >= 0x70 && bank < 0x7e) || bank >= 0xf0) && addr < 0x8000 && cartridge.ramSize > 0 {
		// banks 70-7e and f0-ff, adr 0000-7fff
		// return cartridge.ram[(((uint32(bank)&0xf)<<15)|uint32(addr))&(uint32(cartridge.ramSize)-1)]
		return cartridge.ram.Read((((uint32(bank) & 0xf) << 15) | uint32(addr)) & (uint32(cartridge.ramSize) - 1))
	}
	bank &= 0x7f
	if addr >= 0x8000 || bank >= 0x40 {
		// adr 8000-ffff in all banks or all addresses in banks 40-7f and c0-ff
		return cartridge.rom[((uint32(bank)<<15)|(uint32(addr)&0x7fff))&(cartridge.romSize-1)]
	}
	return cartridge.console.openBus
}

func (cartridge *Cartridge) writeLoROM(bank byte, addr uint16, value byte) {
	if ((bank >= 0x70 && bank < 0x7e) || bank > 0xf0) && addr < 0x8000 && cartridge.ramSize > 0 {
		// banks 70-7e and f0-ff, adr 0000-7fff
		// cartridge.ram[(((uint32(bank)&0xf)<<15)|uint32(addr))&(uint32(cartridge.ramSize)-1)] = value
		cartridge.ram.Write((((uint32(bank)&0xf)<<15)|uint32(addr))&(uint32(cartridge.ramSize)-1), value)
	}
}

func (cartridge *Cartridge) readHiROM(bank byte, addr uint16) byte {
	bank &= 0x7f
	if bank < 0x40 && addr >= 0x6000 && addr < 0x8000 && cartridge.ramSize > 0 {
		// banks 00-3f and 80-bf, adr 6000-7fff
		// return cartridge.ram[(((uint32(bank)&0x3f)<<13)|(uint32(addr)&0x1fff))&(uint32(cartridge.ramSize)-1)]
		return cartridge.ram.Read((((uint32(bank) & 0x3f) << 13) | (uint32(addr) & 0x1fff)) & (uint32(cartridge.ramSize) - 1))
	}
	if addr >= 0x8000 || bank >= 0x40 {
		// addr 8000-ffff in all banks or all addresses in banks 40-7f and c0-ff
		return cartridge.rom[(((uint32(bank)&0x3f)<<16)|uint32(addr))&(uint32(cartridge.romSize)-1)]
	}
	return cartridge.console.openBus
}

func (cartridge *Cartridge) writeHiROM(bank byte, addr uint16, value byte) {
	bank &= 0x7f
	if bank < 0x40 && addr >= 0x6000 && addr < 0x8000 && cartridge.ramSize > 0 {
		// banks 00-3f and 80-bf, adr 6000-7fff
		// cartridge.ram[(((uint32(bank)&0x3f)<<13)|(uint32(addr)&0x1fff))&(uint32(cartridge.ramSize)-1)] = value
		cartridge.ram.Write((((uint32(bank)&0x3f)<<13)|(uint32(addr)&0x1fff))&(uint32(cartridge.ramSize)-1), value)
	}
}

func (cartridge *Cartridge) Close() {
	cartridge.ram.Close()
}
