package chibisnes

type CPU struct {
	console *Console

	a  uint16
	x  uint16
	y  uint16
	sp uint16
	pc uint16
	dp uint16 // direct page (D)
	k  uint8  // program bank (PB)
	db uint8  // data bank (B)

	// flags
	c  byte // Carry
	z  byte // Zero
	i  byte // IRQ disable
	d  byte // Decimal
	xf byte // Index register size (native mode only), (0 = 16-bit, 1 = 8-bit)
	mf byte // Accumulator register size (native mode only), (0 = 16-bit, 1 = 8-bit)
	v  byte // Overflow
	n  byte // Negative
	e  byte // 6502 emulation mode

	// interrupts
	irqWanted bool
	nmiWanted bool

	// power state (WAI/STP)
	waiting bool
	stopped bool

	// internal use
	cyclesUsed   uint8 // indicates how many cycles an opcode used
	cycleCounter uint64
}

const (
	CPUFlagsCarry                  = 0x01
	CPUFlagsZero                   = 0x02
	CPUFlagsInterrupt              = 0x04
	CPUFlagsDecimal                = 0x08
	CPUFlagsIndexRegisterSize      = 0x10
	CPUFlagsAccumulateRegisterSize = 0x20
	CPUFlagsOverflow               = 0x40
	CPUFlagsNegative               = 0x80
)

var cyclesPerCPUOpcode [256]int = [256]int{
	7, 6, 7, 4, 5, 3, 5, 6, 3, 2, 2, 4, 6, 4, 6, 5,
	2, 5, 5, 7, 5, 4, 6, 6, 2, 4, 2, 2, 6, 4, 7, 5,
	6, 6, 8, 4, 3, 3, 5, 6, 4, 2, 2, 5, 4, 4, 6, 5,
	2, 5, 5, 7, 4, 4, 6, 6, 2, 4, 2, 2, 4, 4, 7, 5,
	6, 6, 2, 4, 7, 3, 5, 6, 3, 2, 2, 3, 3, 4, 6, 5,
	2, 5, 5, 7, 7, 4, 6, 6, 2, 4, 3, 2, 4, 4, 7, 5,
	6, 6, 6, 4, 3, 3, 5, 6, 4, 2, 2, 6, 5, 4, 6, 5,
	2, 5, 5, 7, 4, 4, 6, 6, 2, 4, 4, 2, 6, 4, 7, 5,
	3, 6, 4, 4, 3, 3, 3, 6, 2, 2, 2, 3, 4, 4, 4, 5,
	2, 6, 5, 7, 4, 4, 4, 6, 2, 5, 2, 2, 4, 5, 5, 5,
	2, 6, 2, 4, 3, 3, 3, 6, 2, 2, 2, 4, 4, 4, 4, 5,
	2, 5, 5, 7, 4, 4, 4, 6, 2, 4, 2, 2, 4, 4, 4, 5,
	2, 6, 3, 4, 3, 3, 5, 6, 2, 2, 2, 3, 4, 4, 6, 5,
	2, 5, 5, 7, 6, 4, 6, 6, 2, 4, 3, 3, 6, 4, 7, 5,
	2, 6, 3, 4, 3, 3, 5, 6, 2, 2, 2, 3, 4, 4, 6, 5,
	2, 5, 5, 7, 5, 4, 6, 6, 2, 4, 4, 2, 8, 4, 7, 5,
}

func NewCPU(c *Console) *CPU {
	return &CPU{
		console: c,
	}
}

func (cpu *CPU) Reset() {
	cpu.sp = 0x100
	cpu.pc = uint16(cpu.Read(0xFFFC)) | (uint16(cpu.Read(0xFFFD)) << 8)

	// cpu.SetFlags(CPUFlagsInterrupt | CPUFlagsIndexRegisterSize | CPUFlagsAccumulateRegisterSize)
	cpu.i = 0x01
	cpu.xf = 0x01
	cpu.mf = 0x01
	cpu.e = 0x01
}

func (cpu *CPU) SetAllFlags(flags byte) {
	cpu.c = (flags >> 0) & 1
	cpu.z = (flags >> 1) & 1
	cpu.i = (flags >> 2) & 1
	cpu.d = (flags >> 3) & 1
	cpu.xf = (flags >> 4) & 1
	cpu.mf = (flags >> 5) & 1
	cpu.v = (flags >> 6) & 1
	cpu.n = (flags >> 7) & 1

	if cpu.CheckEmulationFlag() {
		cpu.mf = 0x01
		cpu.xf = 0x01
		cpu.sp = (cpu.sp & 0xFF) | 0x0100
	}

	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		cpu.x &= 0xFF
		cpu.y &= 0xFF
	}
}

func (cpu *CPU) SetFlags(flags byte) {
	cpu.SetAllFlags(cpu.Flags() | flags)
}

func (cpu *CPU) ClearFlags(flags byte) {
	cpu.SetAllFlags(cpu.Flags() & ^flags)
}

func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.c << 0
	flags |= cpu.z << 1
	flags |= cpu.i << 2
	flags |= cpu.d << 3
	flags |= cpu.xf << 4
	flags |= cpu.mf << 5
	flags |= cpu.v << 6
	flags |= cpu.n << 7

	return flags
}

func (cpu *CPU) CheckFlag(flag byte) bool {
	return (cpu.Flags() & flag) == flag
}

func (cpu *CPU) SetEmulationFlag(value byte) {
	cpu.e = value & 0x01
	if cpu.e == 0x01 {
		cpu.SetFlags(CPUFlagsAccumulateRegisterSize | CPUFlagsIndexRegisterSize)
		cpu.sp = (cpu.sp & 0xFF) | 0x0100
	}
}

func (cpu *CPU) CheckEmulationFlag() bool {
	if cpu.e == 0x01 {
		return true
	}
	return false
}

func (cpu *CPU) runOpcode() int {
	cpu.cyclesUsed = 0

	if cpu.stopped {
		return 1
	}
	if cpu.waiting {
		if cpu.irqWanted || cpu.nmiWanted {
			cpu.waiting = false
		}
		return 1
	}

	var opcode uint8 = cpu.readOpcode()
	cpu.cyclesUsed = uint8(cyclesPerCPUOpcode[opcode])
	cpu.doOpcode(opcode)

	if (!cpu.CheckFlag(CPUFlagsInterrupt) && cpu.irqWanted) || cpu.nmiWanted {
		cpu.cyclesUsed = 7
		if cpu.nmiWanted {
			cpu.nmiWanted = false
			cpu.doInterrupt(false)
		} else {
			cpu.doInterrupt(true)
		}
	}

	return int(cpu.cyclesUsed)
}

func (cpu *CPU) readOpcode() byte {
	v := cpu.Read((uint32(cpu.k) << 16) | uint32(cpu.pc))
	cpu.pc++
	return v
}

func (cpu *CPU) readOpcodeWord() uint16 {
	var low byte = cpu.readOpcode()
	v := uint16(low) | (uint16(cpu.readOpcode()) << 8)
	return v
}

// setZ sets the zero flag if the argument is zero
func (cpu *CPU) setZ(value uint16, isByte bool) {
	if isByte {
		if (value & 0xFF) == 0 {
			cpu.z = 1
		} else {
			cpu.z = 0
		}
	} else {
		if value == 0 {
			cpu.z = 1
		} else {
			cpu.z = 0
		}
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func (cpu *CPU) setN(value uint16, isByte bool) {
	if isByte {
		if (value & CPUFlagsNegative) == CPUFlagsNegative {
			cpu.n = 1
		} else {
			cpu.n = 0
		}
	} else {
		if (value & (CPUFlagsNegative << 8)) == (CPUFlagsNegative << 8) {
			cpu.n = 1
		} else {
			cpu.n = 0
		}
	}
}

func (cpu *CPU) setZN(value uint16, isByte bool) {
	cpu.setZ(value, isByte)
	cpu.setN(value, isByte)
}

func (cpu *CPU) branch(value byte, check bool) {
	if check {
		cpu.cyclesUsed++
		// XXX: signed
		cpu.pc = uint16(int16(cpu.pc) + int16(int8(value)))
	}
}

func (cpu *CPU) pullByte() byte {
	cpu.sp++
	if cpu.CheckEmulationFlag() {
		cpu.sp = (cpu.sp & 0xFF) | 0x0100
	}
	return cpu.Read(uint32(cpu.sp))
}

func (cpu *CPU) pushByte(value byte) {
	cpu.Write(uint32(cpu.sp), value)
	cpu.sp--
	if cpu.CheckEmulationFlag() {
		cpu.sp = (cpu.sp & 0xFF) | 0x0100
	}
}

func (cpu *CPU) pullWord() uint16 {
	value := cpu.pullByte()
	return uint16(value) | (uint16(cpu.pullByte()) << 8)
}

func (cpu *CPU) pushWord(value uint16) {
	cpu.pushByte(byte(value >> 8))
	cpu.pushByte(byte(value & 0x00FF))
}

func (cpu *CPU) Read(addr uint32) byte {
	return cpu.console.CPURead(addr)
}

func (cpu *CPU) Write(addr uint32, value byte) {
	cpu.console.CPUWrite(addr, value)
}

func (cpu *CPU) ReadWord(addrLow uint32, addrHi uint32) uint16 {
	value := cpu.Read(addrLow)
	return uint16(value) | (uint16(cpu.Read(addrHi)) << 8)
}

func (cpu *CPU) WriteWord(addrLow uint32, addrHi uint32, value uint16, reversed bool) {
	if reversed {
		cpu.Write(addrHi, byte(value>>8))
		cpu.Write(addrLow, byte(value&0x00FF))
	} else {
		cpu.Write(addrLow, byte(value&0xFF))
		cpu.Write(addrHi, byte(value>>8))
	}
}

func (cpu *CPU) doInterrupt(irq bool) {
	cpu.pushByte(cpu.k)
	cpu.pushWord(cpu.pc)
	cpu.pushByte(cpu.Flags())
	cpu.cyclesUsed++
	cpu.SetFlags(CPUFlagsInterrupt)
	cpu.ClearFlags(CPUFlagsDecimal)
	cpu.k = 0

	if irq {
		cpu.pc = cpu.ReadWord(0xFFEE, 0xFFEF)
	} else {
		// NMI
		cpu.pc = cpu.ReadWord(0xFFEA, 0xFFEB)
	}
}

// addressing modes

func (cpu *CPU) addrImm(low *uint32, xFlag bool) uint32 {
	if (xFlag && cpu.CheckFlag(CPUFlagsIndexRegisterSize)) || (!xFlag && cpu.CheckFlag(CPUFlagsAccumulateRegisterSize)) {
		*low = (uint32(cpu.k) << 16) | uint32(cpu.pc)
		cpu.pc++
		return 0
	} else {
		*low = (uint32(cpu.k) << 16) | uint32(cpu.pc)
		cpu.pc++
		v := (uint32(cpu.k) << 16) | uint32(cpu.pc)
		cpu.pc++
		return v
	}
}

func (cpu *CPU) addrDp(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	*low = uint32((cpu.dp + uint16(addr)) & 0xFFFF)
	return uint32((cpu.dp + uint16(addr) + 1) & 0xFFFF)
}

func (cpu *CPU) addrDpx(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := cpu.dp + uint16(addr) + cpu.x
	*low = uint32(base & 0xFFFF)
	return uint32((base + 1) & 0xFFFF)
}

func (cpu *CPU) addrDpy(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := cpu.dp + uint16(addr) + cpu.y
	*low = uint32(base & 0xFFFF)
	return uint32((base + 1) & 0xFFFF)
}

func (cpu *CPU) addrIdp(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := (cpu.dp + uint16(addr))
	pointer := cpu.ReadWord(uint32(base&0xFFFF), uint32((base+1)&0xFFFF))
	*low = (uint32(cpu.db) << 16) + uint32(pointer)
	return ((uint32(cpu.db) << 16) + uint32(pointer) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrIdx(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := cpu.dp + uint16(addr) + cpu.x
	pointer := cpu.ReadWord(uint32((base)&0xFFFF), uint32((base+1)&0xFFFF))
	*low = (uint32(cpu.db) << 16) + uint32(pointer)
	return ((uint32(cpu.db) << 16) + uint32(pointer) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrIdy(low *uint32, write bool) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := cpu.dp + uint16(addr)
	pointer := cpu.ReadWord(uint32(base&0xFFFF), uint32((base+1)&0xFFFF))
	if write && (!cpu.CheckFlag(CPUFlagsIndexRegisterSize) || ((pointer >> 8) != ((pointer + cpu.y) >> 8))) {
		cpu.cyclesUsed++
	}
	// x = 0 or page crossed, with writing opcode: 1 extra cycle
	*low = (uint32(cpu.db) << 16) + uint32(pointer) + uint32(cpu.y)&0xFFFFFFFF
	return ((uint32(cpu.db) << 16) + uint32(pointer) + uint32(cpu.y) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrIdl(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := (cpu.dp + uint16(addr))
	pointer := uint32(cpu.ReadWord(uint32(base&0xFFFF), uint32((base+1)&0xFFFF)))
	pointer |= uint32(cpu.Read(uint32((base+2)&0xFFFF))) << 16
	*low = pointer
	return (pointer + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrIly(low *uint32) uint32 {
	addr := cpu.readOpcode()
	if (cpu.dp & 0xFF) > 0 {
		cpu.cyclesUsed++
	}
	base := (cpu.dp + uint16(addr))
	pointer := uint32(cpu.ReadWord(uint32(base&0xFFFF), uint32((base+1)&0xFFFF)))
	pointer |= uint32(cpu.Read(uint32((base+2)&0xFFFF))) << 16
	*low = (pointer + uint32(cpu.y)) & 0xFFFFFFFF
	return (pointer + uint32(cpu.y) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrSr(low *uint32) uint32 {
	addr := cpu.readOpcode()
	base := cpu.sp + uint16(addr)
	*low = uint32((base) & 0xFFFF)
	return uint32((base + 1) & 0xFFFF)
}

func (cpu *CPU) addrIsy(low *uint32) uint32 {
	addr := cpu.readOpcode()
	base := cpu.sp + uint16(addr)
	pointer := cpu.ReadWord(uint32(base&0xFFFF), uint32((base+1)&0xFFFF))
	*low = ((uint32(cpu.db) << 16) + uint32(pointer) + uint32(cpu.y)) & 0xFFFFFFFF
	return ((uint32(cpu.db) << 16) + uint32(pointer) + uint32(cpu.y) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrAbs(low *uint32) uint32 {
	addr := cpu.readOpcodeWord()
	*low = (uint32(cpu.db) << 16) + uint32(addr)
	return ((uint32(cpu.db) << 16) + uint32(addr) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrAbx(low *uint32, write bool) uint32 {
	addr := cpu.readOpcodeWord()
	if write && (!cpu.CheckFlag(CPUFlagsIndexRegisterSize) || ((addr >> 8) != ((addr + cpu.x) >> 8))) {
		cpu.cyclesUsed++
	}
	base := (uint32(cpu.db) << 16) + uint32(addr) + uint32(cpu.x)
	*low = (base & 0xFFFFFFFF)
	return (base + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrAby(low *uint32, write bool) uint32 {
	addr := cpu.readOpcodeWord()
	if write && (!cpu.CheckFlag(CPUFlagsIndexRegisterSize) || ((addr >> 8) != ((addr + cpu.y) >> 8))) {
		cpu.cyclesUsed++
	}
	base := (uint32(cpu.db) << 16) + uint32(addr) + uint32(cpu.y)
	*low = (base & 0xFFFFFFFF)
	return (base + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrAbl(low *uint32) uint32 {
	addr := uint32(cpu.readOpcodeWord())
	addr |= uint32(cpu.readOpcode()) << 16
	*low = addr
	return (addr + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrAlx(low *uint32) uint32 {
	addr := uint32(cpu.readOpcodeWord())
	addr |= uint32(cpu.readOpcode()) << 16
	*low = addr + uint32(cpu.x)&0xFFFFFFFF
	return (addr + uint32(cpu.x) + 1) & 0xFFFFFFFF
}

func (cpu *CPU) addrIax() uint16 {
	addr := cpu.readOpcodeWord()
	baseLow := uint32(addr) + uint32(cpu.x)
	baseHigh := uint32(cpu.k) << 16
	v := cpu.ReadWord((baseHigh | (baseLow & 0xFFFF)), (baseHigh | ((baseLow + 1) & 0xFFFF)))
	return v
}
