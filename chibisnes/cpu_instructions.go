package chibisnes

func (cpu *CPU) doOpcode(opcode byte) {
	switch opcode {
	case 0x00:
		// brk imp
		cpu.pushByte(cpu.k)
		cpu.pushWord(cpu.pc + 1)
		cpu.pushByte(cpu.Flags())
		cpu.cyclesUsed++ // native mode: 1 extra cycle
		cpu.SetFlags(CPUFlagsInterrupt)
		cpu.ClearFlags(CPUFlagsDecimal)
		cpu.k = 0
		cpu.pc = cpu.ReadWord(0xffe6, 0xffe7)
	case 0x01:
		// ora idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.ora(low, high)
	case 0x02:
		// cop imm(s)
		cpu.readOpcode()
		cpu.pushByte(cpu.k)
		cpu.pushByte(byte(cpu.pc))
		cpu.pushByte(cpu.Flags())
		cpu.cyclesUsed++ // native mode: 1 extra cycle
		cpu.SetFlags(CPUFlagsInterrupt)
		cpu.ClearFlags(CPUFlagsDecimal)
		cpu.k = 0
		cpu.pc = cpu.ReadWord(0xffe4, 0xffe5)
	case 0x03:
		// ora sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.ora(low, high)
	case 0x04:
		// tsb dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.tsb(low, high)
	case 0x05:
		// ora dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.ora(low, high)
	case 0x06:
		// asl dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.asl(low, high)
	case 0x07:
		// ora idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.ora(low, high)
	case 0x08:
		// php imp
		cpu.pushByte(cpu.Flags())
	case 0x09:
		// ora imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.ora(low, high)
	case 0x0a:
		// asla imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			if (cpu.a & 0x80) > 0 {
				cpu.SetFlags(CPUFlagsCarry)
			} else {
				cpu.ClearFlags(CPUFlagsCarry)
			}
			cpu.a = (cpu.a & 0xff00) | ((cpu.a << 1) & 0xff)
		} else {
			if (cpu.a & 0x8000) > 0 {
				cpu.SetFlags(CPUFlagsCarry)
			} else {
				cpu.ClearFlags(CPUFlagsCarry)
			}
			cpu.a <<= 1
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x0b:
		// phd imp
		cpu.pushWord(cpu.dp)
	case 0x0c:
		// tsb abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.tsb(low, high)
	case 0x0d:
		// ora abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.ora(low, high)
	case 0x0e:
		// asl abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.asl(low, high)
	case 0x0f:
		// ora abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.ora(low, high)
	case 0x10:
		// bpl rel
		cpu.branch(cpu.readOpcode(), !cpu.CheckFlag(CPUFlagsNegative))
	case 0x11:
		// ora idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.ora(low, high)
	case 0x12:
		// ora idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.ora(low, high)
	case 0x13:
		// ora isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.ora(low, high)
	case 0x14:
		// trb dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.trb(low, high)
	case 0x15:
		// ora dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.ora(low, high)
	case 0x16:
		// asl dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.asl(low, high)
	case 0x17:
		// ora ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.ora(low, high)
	case 0x18:
		// clc imp
		cpu.ClearFlags(CPUFlagsCarry)
	case 0x19:
		// ora aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.ora(low, high)
	case 0x1a:
		// inca imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | ((cpu.a + 1) & 0xff)
		} else {
			cpu.a++
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x1b:
		// tcs imp
		cpu.sp = cpu.a
	case 0x1c:
		// trb abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.trb(low, high)
	case 0x1d:
		// ora abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.ora(low, high)
	case 0x1e:
		// asl abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.asl(low, high)
	case 0x1f:
		// ora alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.ora(low, high)
	case 0x20:
		// jsr abs
		var value uint16 = cpu.readOpcodeWord()
		cpu.pushWord(cpu.pc - 1)
		cpu.pc = value
	case 0x21:
		// and idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.and(low, high)
	case 0x22:
		// jsl abl
		var value uint16 = cpu.readOpcodeWord()
		var newK byte = cpu.readOpcode()
		cpu.pushByte(cpu.k)
		cpu.pushWord(cpu.pc - 1)
		cpu.pc = value
		cpu.k = newK
	case 0x23:
		// and sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.and(low, high)
	case 0x24:
		// bit dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.bit(low, high)
	case 0x25:
		// and dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.and(low, high)
	case 0x26:
		// rol dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.rol(low, high)
	case 0x27:
		// and idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.and(low, high)
	case 0x28:
		// plp imp
		cpu.SetAllFlags(cpu.pullByte())
	case 0x29:
		// and imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.and(low, high)
	case 0x2a:
		// rola imp
		var result int = (int(cpu.a) << 1) | int(cpu.c)

		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			if (result & 0x0100) > 0 {
				cpu.SetFlags(CPUFlagsCarry)
			} else {
				cpu.ClearFlags(CPUFlagsCarry)
			}
			cpu.a = (cpu.a & 0xff00) | (uint16(result) & 0xff)
		} else {
			if (result & 0x10000) > 0 {
				cpu.SetFlags(CPUFlagsCarry)
			} else {
				cpu.ClearFlags(CPUFlagsCarry)
			}
			cpu.a = uint16(result)
		}

		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x2b:
		// pld imp
		cpu.dp = cpu.pullWord()
		cpu.setZN(cpu.dp, false)
	case 0x2c:
		// bit abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.bit(low, high)
	case 0x2d:
		// and abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.and(low, high)
	case 0x2e:
		// rol abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.rol(low, high)
	case 0x2f:
		// and abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.and(low, high)
	case 0x30:
		// bmi rel
		cpu.branch(cpu.readOpcode(), cpu.CheckFlag(CPUFlagsNegative))
	case 0x31:
		// and idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.and(low, high)
	case 0x32:
		// and idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.and(low, high)
	case 0x33:
		// and isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.and(low, high)
	case 0x34:
		// bit dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.bit(low, high)
	case 0x35:
		// and dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.and(low, high)
	case 0x36:
		// rol dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.rol(low, high)
	case 0x37:
		// and ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.and(low, high)
	case 0x38:
		// sec imp
		cpu.SetFlags(CPUFlagsCarry)
	case 0x39:
		// and aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.and(low, high)
	case 0x3a:
		// deca imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | ((cpu.a - 1) & 0xff)
		} else {
			cpu.a--
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x3b:
		// tsc imp
		cpu.a = cpu.sp
		cpu.setZN(cpu.a, false)
	case 0x3c:
		// bit abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.bit(low, high)
	case 0x3d:
		// and abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.and(low, high)
	case 0x3e:
		// rol abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.rol(low, high)
	case 0x3f:
		// and alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.and(low, high)
	case 0x40:
		// rti imp
		cpu.SetAllFlags(cpu.pullByte())
		cpu.cyclesUsed++ // native mode: 1 extra cycle
		cpu.pc = cpu.pullWord()
		cpu.k = cpu.pullByte()
	case 0x41:
		// eor idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.eor(low, high)
	case 0x42:
		// wdm imm(s)
		cpu.readOpcode()
	case 0x43:
		// eor sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.eor(low, high)
	case 0x44:
		// mvp bm
		var dest byte = cpu.readOpcode()
		var src byte = cpu.readOpcode()
		cpu.db = dest
		cpu.Write((uint32(dest)<<16)|uint32(cpu.y), cpu.Read((uint32(src)<<16)|uint32(cpu.x)))
		cpu.a--
		cpu.x--
		cpu.y--
		if cpu.a != 0xffff {
			cpu.pc -= 3
		}
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x &= 0xff
			cpu.y &= 0xff
		}
	case 0x45:
		// eor dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.eor(low, high)
	case 0x46:
		// lsr dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.lsr(low, high)
	case 0x47:
		// eor idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.eor(low, high)
	case 0x48:
		// pha imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.pushByte(byte(cpu.a))
		} else {
			cpu.cyclesUsed++ // m = 0: 1 extra cycle
			cpu.pushWord(cpu.a)
		}
	case 0x49:
		// eor imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.eor(low, high)
	case 0x4a:
		// lsra imp
		if (cpu.a & 1) > 0 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | ((cpu.a >> 1) & 0x7f)
		} else {
			cpu.a >>= 1
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x4b:
		// phk imp
		cpu.pushByte(cpu.k)
	case 0x4c:
		// jmp abs
		cpu.pc = cpu.readOpcodeWord()
	case 0x4d:
		// eor abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.eor(low, high)
	case 0x4e:
		// lsr abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.lsr(low, high)
	case 0x4f:
		// eor abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.eor(low, high)
	case 0x50:
		// bvc rel
		cpu.branch(cpu.readOpcode(), !cpu.CheckFlag(CPUFlagsOverflow))
	case 0x51:
		// eor idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.eor(low, high)
	case 0x52:
		// eor idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.eor(low, high)
	case 0x53:
		// eor isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.eor(low, high)
	case 0x54:
		// mvn bm
		var dest byte = cpu.readOpcode()
		var src byte = cpu.readOpcode()
		cpu.db = dest
		cpu.Write((uint32(dest)<<16)|uint32(cpu.y), cpu.Read((uint32(src)<<16)|uint32(cpu.x)))
		cpu.a--
		cpu.x++
		cpu.y++
		if cpu.a != 0xffff {
			cpu.pc -= 3
		}
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x &= 0xff
			cpu.y &= 0xff
		}
	case 0x55:
		// eor dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.eor(low, high)
	case 0x56:
		// lsr dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.lsr(low, high)
	case 0x57:
		// eor ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.eor(low, high)
	case 0x58:
		// cli imp
		cpu.ClearFlags(CPUFlagsInterrupt)
	case 0x59:
		// eor aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.eor(low, high)
	case 0x5a:
		// phy imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.pushByte(byte(cpu.y))
		} else {
			cpu.cyclesUsed++ // m = 0: 1 extra cycle
			cpu.pushWord(cpu.y)
		}
	case 0x5b:
		// tcd imp
		cpu.dp = cpu.a
		cpu.setZN(cpu.dp, false)
	case 0x5c:
		// jml abl
		var value uint16 = cpu.readOpcodeWord()
		cpu.k = cpu.readOpcode()
		cpu.pc = value
	case 0x5d:
		// eor abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.eor(low, high)
	case 0x5e:
		// lsr abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.lsr(low, high)
	case 0x5f:
		// eor alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.eor(low, high)
	case 0x60:
		// rts imp
		cpu.pc = cpu.pullWord() + 1
	case 0x61:
		// adc idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.adc(low, high)
	case 0x62:
		// per rll
		var value uint16 = cpu.readOpcodeWord()
		cpu.pushWord(cpu.pc + uint16(value))
	case 0x63:
		// adc sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.adc(low, high)
	case 0x64:
		// stz dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.stz(low, high)
	case 0x65:
		// adc dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.adc(low, high)
	case 0x66:
		// ror dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.ror(low, high)
	case 0x67:
		// adc idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.adc(low, high)
	case 0x68:
		// pla imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | uint16(cpu.pullByte())
		} else {
			cpu.cyclesUsed++ // 16-bit m: 1 extra cycle
			cpu.a = cpu.pullWord()
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x69:
		// adc imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.adc(low, high)
	case 0x6a:
		// rora imp
		var carry bool = (cpu.a & 1) > 0
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | ((cpu.a >> 1) & 0x7f) | (uint16(cpu.c) << 7)
		} else {
			cpu.a = (cpu.a >> 1) | (uint16(cpu.c) << 15)
		}
		if carry {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x6b:
		// rtl imp
		cpu.pc = cpu.pullWord() + 1
		cpu.k = cpu.pullByte()
	case 0x6c:
		// jmp ind
		var addr uint32 = uint32(cpu.readOpcodeWord())
		cpu.pc = cpu.ReadWord(addr, (addr+1)&0xffff)
	case 0x6d:
		// adc abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.adc(low, high)
	case 0x6e:
		// ror abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.ror(low, high)
	case 0x6f:
		// adc abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.adc(low, high)
	case 0x70:
		// bvs rel
		cpu.branch(cpu.readOpcode(), cpu.CheckFlag(CPUFlagsOverflow))
	case 0x71:
		// adc idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.adc(low, high)
	case 0x72:
		// adc idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.adc(low, high)
	case 0x73:
		// adc isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.adc(low, high)
	case 0x74:
		// stz dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.stz(low, high)
	case 0x75:
		// adc dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.adc(low, high)
	case 0x76:
		// ror dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.ror(low, high)
	case 0x77:
		// adc ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.adc(low, high)
	case 0x78:
		// sei imp
		cpu.SetFlags(CPUFlagsInterrupt)
	case 0x79:
		// adc aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.adc(low, high)
	case 0x7a:
		// ply imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.y = uint16(cpu.pullByte())
		} else {
			cpu.cyclesUsed++ // 16-bit x: 1 extra cycle
			cpu.y = uint16(cpu.pullWord())
		}
		cpu.setZN(cpu.y, cpu.xf == 1)
	case 0x7b:
		// tdc imp
		cpu.a = cpu.dp
		cpu.setZN(cpu.a, false)
	case 0x7c:
		// jmp iax
		cpu.pc = cpu.addrIax()
	case 0x7d:
		// adc abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.adc(low, high)
	case 0x7e:
		// ror abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.ror(low, high)
	case 0x7f:
		// adc alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.adc(low, high)
	case 0x80:
		// bra rel
		// XXX: signed
		cpu.pc = uint16(int16(cpu.pc) + int16(int8(cpu.readOpcode())))
	case 0x81:
		// sta idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.sta(low, high)
	case 0x82:
		// brl rll
		cpu.pc += uint16(cpu.readOpcodeWord())
	case 0x83:
		// sta sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.sta(low, high)
	case 0x84:
		// sty dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.sty(low, high)
	case 0x85:
		// sta dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.sta(low, high)
	case 0x86:
		// stx dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.stx(low, high)
	case 0x87:
		// sta idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.sta(low, high)
	case 0x88:
		// dey imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.y = (cpu.y - 1) & 0xff
		} else {
			cpu.y--
		}
		cpu.setZN(cpu.y, cpu.xf == 1)
	case 0x89:
		// biti imm(m)
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			var result byte = byte((cpu.a & 0xff) & uint16(cpu.readOpcode()))
			if result == 0 {
				cpu.SetFlags(CPUFlagsZero)
			} else {
				cpu.ClearFlags(CPUFlagsZero)
			}
		} else {
			cpu.cyclesUsed++ // m = 0: 1 extra cycle
			var result uint16 = cpu.a & cpu.readOpcodeWord()
			if result == 0 {
				cpu.SetFlags(CPUFlagsZero)
			} else {
				cpu.ClearFlags(CPUFlagsZero)
			}
		}
	case 0x8a:
		// txa imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | (cpu.x & 0xff)
		} else {
			cpu.a = cpu.x
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x8b:
		// phb imp
		cpu.pushByte(cpu.db)
	case 0x8c:
		// sty abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.sty(low, high)
	case 0x8d:
		// sta abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.sta(low, high)
	case 0x8e:
		// stx abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.stx(low, high)
	case 0x8f:
		// sta abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.sta(low, high)
	case 0x90:
		// bcc rel
		cpu.branch(cpu.readOpcode(), !cpu.CheckFlag(CPUFlagsCarry))
	case 0x91:
		// sta idy
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, true)
		cpu.sta(low, high)
	case 0x92:
		// sta idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.sta(low, high)
	case 0x93:
		// sta isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.sta(low, high)
	case 0x94:
		// sty dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.sty(low, high)
	case 0x95:
		// sta dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.sta(low, high)
	case 0x96:
		// stx dpy
		var low uint32 = 0
		var high uint32 = cpu.addrDpy(&low)
		cpu.stx(low, high)
	case 0x97:
		// sta ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.sta(low, high)
	case 0x98:
		// tya imp
		if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
			cpu.a = (cpu.a & 0xff00) | (cpu.y & 0xff)
		} else {
			cpu.a = cpu.y
		}
		cpu.setZN(cpu.a, cpu.mf == 1)
	case 0x99:
		// sta aby
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, true)
		cpu.sta(low, high)
	case 0x9a:
		// txs imp
		cpu.sp = cpu.x
	case 0x9b:
		// txy imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.y = cpu.x & 0xff
		} else {
			cpu.y = cpu.x
		}
		cpu.setZN(cpu.y, cpu.xf == 1)
	case 0x9c:
		// stz abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.stz(low, high)
	case 0x9d:
		// sta abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.sta(low, high)
	case 0x9e:
		// stz abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.stz(low, high)
	case 0x9f:
		// sta alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.sta(low, high)
	case 0xa0:
		// ldy imm(x)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, true)
		cpu.ldy(low, high)
	case 0xa1:
		// lda idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.lda(low, high)
	case 0xa2:
		// ldx imm(x)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, true)
		cpu.ldx(low, high)
	case 0xa3:
		// lda sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.lda(low, high)
	case 0xa4:
		// ldy dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.ldy(low, high)
	case 0xa5:
		// lda dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.lda(low, high)
	case 0xa6:
		// ldx dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.ldx(low, high)
	case 0xa7:
		// lda idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.lda(low, high)
	case 0xa8:
		// tay imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.y = cpu.a & 0xff
		} else {
			cpu.y = cpu.a
		}
		cpu.setZN(cpu.y, cpu.xf == 1)
	case 0xa9:
		// lda imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.lda(low, high)
	case 0xaa:
		// tax imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = cpu.a & 0xff
		} else {
			cpu.x = cpu.a
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xab:
		// plb imp
		cpu.db = cpu.pullByte()
		cpu.setZN(uint16(cpu.db), true)
	case 0xac:
		// ldy abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.ldy(low, high)
	case 0xad:
		// lda abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.lda(low, high)
	case 0xae:
		// ldx abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.ldx(low, high)
	case 0xaf:
		// lda abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.lda(low, high)
	case 0xb0:
		// bcs rel
		cpu.branch(cpu.readOpcode(), cpu.CheckFlag(CPUFlagsCarry))
	case 0xb1:
		// lda idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.lda(low, high)
	case 0xb2:
		// lda idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.lda(low, high)
	case 0xb3:
		// lda isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.lda(low, high)
	case 0xb4:
		// ldy dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.ldy(low, high)
	case 0xb5:
		// lda dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.lda(low, high)
	case 0xb6:
		// ldx dpy
		var low uint32 = 0
		var high uint32 = cpu.addrDpy(&low)
		cpu.ldx(low, high)
	case 0xb7:
		// lda ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.lda(low, high)
	case 0xb8:
		// clv imp
		cpu.ClearFlags(CPUFlagsOverflow)
	case 0xb9:
		// lda aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.lda(low, high)
	case 0xba:
		// tsx imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = cpu.sp & 0xff
		} else {
			cpu.x = cpu.sp
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xbb:
		// tyx imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = cpu.y & 0xff
		} else {
			cpu.x = cpu.y
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xbc:
		// ldy abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.ldy(low, high)
	case 0xbd:
		// lda abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.lda(low, high)
	case 0xbe:
		// ldx aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.ldx(low, high)
	case 0xbf:
		// lda alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.lda(low, high)
	case 0xc0:
		// cpy imm(x)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, true)
		cpu.cpy(low, high)
	case 0xc1:
		// cmp idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.cmp(low, high)
	case 0xc2:
		// rep imm(s)
		v := cpu.readOpcode()
		cpu.SetAllFlags(cpu.Flags() & ^v)
	case 0xc3:
		// cmp sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.cmp(low, high)
	case 0xc4:
		// cpy dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.cpy(low, high)
	case 0xc5:
		// cmp dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.cmp(low, high)
	case 0xc6:
		// dec dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.dec(low, high)
	case 0xc7:
		// cmp idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.cmp(low, high)
	case 0xc8:
		// iny imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.y = (cpu.y + 1) & 0xff
		} else {
			cpu.y++
		}
		cpu.setZN(cpu.y, cpu.xf == 1)
	case 0xc9:
		// cmp imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.cmp(low, high)
	case 0xca:
		// dex imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = (cpu.x - 1) & 0xff
		} else {
			cpu.x--
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xcb:
		// wai imp
		cpu.waiting = true
	case 0xcc:
		// cpy abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.cpy(low, high)
	case 0xcd:
		// cmp abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.cmp(low, high)
	case 0xce:
		// dec abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.dec(low, high)
	case 0xcf:
		// cmp abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.cmp(low, high)
	case 0xd0:
		// bne rel
		cpu.branch(cpu.readOpcode(), !cpu.CheckFlag(CPUFlagsZero))
	case 0xd1:
		// cmp idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.cmp(low, high)
	case 0xd2:
		// cmp idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.cmp(low, high)
	case 0xd3:
		// cmp isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.cmp(low, high)
	case 0xd4:
		// pei dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.pushWord(cpu.ReadWord(low, high))
	case 0xd5:
		// cmp dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.cmp(low, high)
	case 0xd6:
		// dec dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.dec(low, high)
	case 0xd7:
		// cmp ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.cmp(low, high)
	case 0xd8:
		// cld imp
		cpu.ClearFlags(CPUFlagsDecimal)
	case 0xd9:
		// cmp aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.cmp(low, high)
	case 0xda:
		// phx imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.pushByte(byte(cpu.x))
		} else {
			cpu.cyclesUsed++ // m = 0: 1 extra cycle
			cpu.pushWord(cpu.x)
		}
	case 0xdb:
		// stp imp
		cpu.stopped = true
	case 0xdc:
		// jml ial
		var addr uint32 = uint32(cpu.readOpcodeWord())
		cpu.pc = cpu.ReadWord(addr, (addr+1)&0xffff)
		cpu.k = cpu.Read((addr + 2) & 0xffff)
	case 0xdd:
		// cmp abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.cmp(low, high)
	case 0xde:
		// dec abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.dec(low, high)
	case 0xdf:
		// cmp alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.cmp(low, high)
	case 0xe0:
		// cpx imm(x)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, true)
		cpu.cpx(low, high)
	case 0xe1:
		// sbc idx
		var low uint32 = 0
		var high uint32 = cpu.addrIdx(&low)
		cpu.sbc(low, high)
	case 0xe2:
		// sep imm(s)
		cpu.SetAllFlags(cpu.Flags() | cpu.readOpcode())
	case 0xe3:
		// sbc sr
		var low uint32 = 0
		var high uint32 = cpu.addrSr(&low)
		cpu.sbc(low, high)
	case 0xe4:
		// cpx dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.cpx(low, high)
	case 0xe5:
		// sbc dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.sbc(low, high)
	case 0xe6:
		// inc dp
		var low uint32 = 0
		var high uint32 = cpu.addrDp(&low)
		cpu.inc(low, high)
	case 0xe7:
		// sbc idl
		var low uint32 = 0
		var high uint32 = cpu.addrIdl(&low)
		cpu.sbc(low, high)
	case 0xe8:
		// inx imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = (cpu.x + 1) & 0xff
		} else {
			cpu.x++
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xe9:
		// sbc imm(m)
		var low uint32 = 0
		var high uint32 = cpu.addrImm(&low, false)
		cpu.sbc(low, high)
	case 0xea:
		// nop imp
		// no operation
	case 0xeb:
		// xba imp
		var low uint16 = cpu.a & 0xff
		var high uint16 = cpu.a >> 8
		cpu.a = (low << 8) | high
		cpu.setZN(high, true)
	case 0xec:
		// cpx abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.cpx(low, high)
	case 0xed:
		// sbc abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.sbc(low, high)
	case 0xee:
		// inc abs
		var low uint32 = 0
		var high uint32 = cpu.addrAbs(&low)
		cpu.inc(low, high)
	case 0xef:
		// sbc abl
		var low uint32 = 0
		var high uint32 = cpu.addrAbl(&low)
		cpu.sbc(low, high)
	case 0xf0:
		// beq rel
		cpu.branch(cpu.readOpcode(), cpu.CheckFlag(CPUFlagsZero))
	case 0xf1:
		// sbc idy(r)
		var low uint32 = 0
		var high uint32 = cpu.addrIdy(&low, false)
		cpu.sbc(low, high)
	case 0xf2:
		// sbc idp
		var low uint32 = 0
		var high uint32 = cpu.addrIdp(&low)
		cpu.sbc(low, high)
	case 0xf3:
		// sbc isy
		var low uint32 = 0
		var high uint32 = cpu.addrIsy(&low)
		cpu.sbc(low, high)
	case 0xf4:
		// pea imm(l)
		cpu.pushWord(cpu.readOpcodeWord())
	case 0xf5:
		// sbc dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.sbc(low, high)
	case 0xf6:
		// inc dpx
		var low uint32 = 0
		var high uint32 = cpu.addrDpx(&low)
		cpu.inc(low, high)
	case 0xf7:
		// sbc ily
		var low uint32 = 0
		var high uint32 = cpu.addrIly(&low)
		cpu.sbc(low, high)
	case 0xf8:
		// sed imp
		cpu.SetFlags(CPUFlagsDecimal)
	case 0xf9:
		// sbc aby(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAby(&low, false)
		cpu.sbc(low, high)
	case 0xfa:
		// plx imp
		if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
			cpu.x = uint16(cpu.pullByte())
		} else {
			cpu.cyclesUsed++ // 16-bit x: 1 extra cycle
			cpu.x = cpu.pullWord()
		}
		cpu.setZN(cpu.x, cpu.xf == 1)
	case 0xfb:
		// xce imp
		var temp byte = cpu.c
		cpu.c = cpu.e
		cpu.e = temp
		cpu.SetAllFlags(cpu.Flags()) // updates x and m flags, clears upper half of x and y if needed
	case 0xfc:
		// jsr iax
		var value uint16 = cpu.addrIax()
		cpu.pushWord(cpu.pc - 1)
		cpu.pc = value
	case 0xfd:
		// sbc abx(r)
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, false)
		cpu.sbc(low, high)
	case 0xfe:
		// inc abx
		var low uint32 = 0
		var high uint32 = cpu.addrAbx(&low, true)
		cpu.inc(low, high)
	case 0xff:
		// sbc alx
		var low uint32 = 0
		var high uint32 = cpu.addrAlx(&low)
		cpu.sbc(low, high)
	}
}

// opcode functions
func (cpu *CPU) and(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		cpu.a = (cpu.a & 0xFF00) | ((cpu.a & uint16(value)) & 0xFF)
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high)
		cpu.a &= value
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) ora(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		cpu.a = (cpu.a & 0xFF00) | ((cpu.a | uint16(value)) & 0xFF)
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high)
		cpu.a |= value
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) eor(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		cpu.a = (cpu.a & 0xFF00) | ((cpu.a ^ uint16(value)) & 0xFF)
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high)
		cpu.a ^= value
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) adc(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		var result int = 0
		if cpu.CheckFlag(CPUFlagsDecimal) {
			result = (int(cpu.a) & 0xF) + (int(value) & 0xF) + int(cpu.c)
			if result > 0x9 {
				result = ((result + 0x6) & 0xF) + 0x10
			}
			result = (int(cpu.a) & 0xF0) + (int(value) & 0xF0) + result
		} else {
			result = (int(cpu.a) & 0xFF) + int(value) + int(cpu.c)
		}
		if ((cpu.a & 0x80) == (uint16(value) & 0x80)) && ((int(value) & 0x80) != (result & 0x80)) {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
		if cpu.CheckFlag(CPUFlagsDecimal) && result > 0x9F {
			result += 0x60
		}
		if result > 0xFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.a = (cpu.a & 0xFF00) | (uint16(result) & 0x00FF)
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high)
		var result int = 0
		if cpu.CheckFlag(CPUFlagsDecimal) {
			result = (int(cpu.a) & 0xF) + (int(value) & 0xF) + int(cpu.c)
			if result > 0x9 {
				result = ((result + 0x6) & 0xF) + 0x10
			}
			result = (int(cpu.a) & 0xF0) + (int(value) & 0xF0) + result
			if result > 0x9f {
				result = ((result + 0x60) & 0xFF) + 0x100
			}
			result = (int(cpu.a) & 0xF00) + (int(value) & 0xF00) + result
			if result > 0x9FF {
				result = ((result + 0x600) & 0xFFF) + 0x1000
			}
			result = (int(cpu.a) & 0xF000) + (int(value) & 0xF000) + result
		} else {
			result = int(cpu.a) + int(value) + int(cpu.c)
		}
		if ((cpu.a & 0x8000) == (value & 0x8000)) && ((int(value) & 0x8000) != (result & 0x8000)) {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
		if cpu.CheckFlag(CPUFlagsDecimal) && result > 0x9FFF {
			result += 0x6000
		}
		if result > 0xFFFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.a = uint16(result)
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) sbc(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low) ^ 0xFF
		var result int = 0
		if cpu.CheckFlag(CPUFlagsDecimal) {
			result = (int(cpu.a) & 0xF) + (int(value) & 0xF) + int(cpu.c)
			if result < 0x10 {
				if result-0x6 < 0 {
					result = (result - 0x6) & 0xF
				} else {
					result = (result - 0x6) & 0x1F
				}
			}
			result = (int(cpu.a) & 0xF0) + (int(value) & 0xF0) + result
		} else {
			result = (int(cpu.a) & 0xff) + int(value) + int(cpu.c)
		}
		if ((cpu.a & 0x80) == (uint16(value) & 0x80)) && ((int(value) & 0x80) != (result & 0x80)) {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
		if cpu.CheckFlag(CPUFlagsDecimal) && result < 0x0100 {
			result -= 0x60
		}
		if result > 0xFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.a = (cpu.a & 0xFF00) | (uint16(result) & 0x00FF)
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high) ^ 0xFFFF
		var result int = 0
		if cpu.CheckFlag(CPUFlagsDecimal) {
			result = (int(cpu.a) & 0xF) + (int(value) & 0xF) + int(cpu.c)
			if result < 0x10 {
				if result-0x6 < 0 {
					result = (result - 0x6) & 0xF
				} else {
					result = (result - 0x6) & 0x1F
				}
			}
			result = (int(cpu.a) & 0xf0) + (int(value) & 0xf0) + result
			if result < 0x100 {
				if result-0x60 < 0 {
					result = (result - 0x60) & 0xFF
				} else {
					result = (result - 0x60) & 0x1FF
				}
			}
			result = (int(cpu.a) & 0xf00) + (int(value) & 0xf00) + result
			if result < 0x1000 {
				if result-0x600 < 0 {
					result = (result - 0x600) & 0xFFF
				} else {
					result = (result - 0x600) & 0x1FFF
				}
			}
			result = (int(cpu.a) & 0xf000) + (int(value) & 0xf000) + result
		} else {
			result = int(cpu.a) + int(value) + int(cpu.c)
		}
		if ((cpu.a & 0x8000) == (value & 0x8000)) && ((int(value) & 0x8000) != (result & 0x8000)) {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
		if cpu.CheckFlag(CPUFlagsDecimal) && result < 0x10000 {
			result -= 0x6000
		}
		if result > 0xFFFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.a = uint16(result)
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) cmp(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		v := cpu.console.CPU.Read(low)
		value := v ^ 0xFF
		result = (int(cpu.a) & 0xFF) + int(value) + 1
		if result > 0xFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high) ^ 0xFFFF
		result = int(cpu.a) + int(value) + 1
		if result > 0xFFFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) cpx(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		value := cpu.Read(low) ^ 0xFF
		result = (int(cpu.x) & 0xFF) + int(value) + 1
		if result > 0xFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high) ^ 0xFFFF
		result = int(cpu.x) + int(value) + 1
		if result > 0xFFFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	}
	cpu.setZN(uint16(result), cpu.xf == 1)
}

func (cpu *CPU) cpy(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		value := cpu.Read(low) ^ 0xFF
		result = (int(cpu.y) & 0xFF) + int(value) + 1
		if result > 0xFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high) ^ 0xFFFF
		result = int(cpu.y) + int(value) + 1
		if result > 0xFFFF {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
	}
	cpu.setZN(uint16(result), cpu.xf == 1)
}

func (cpu *CPU) bit(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		result := (int(cpu.a) & 0xFF) & int(value)
		if result == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		if (value & 0x80) > 0 {
			cpu.SetFlags(CPUFlagsNegative)
		} else {
			cpu.ClearFlags(CPUFlagsNegative)
		}
		if (value & 0x40) > 0 {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
	} else {
		cpu.cyclesUsed++
		value := cpu.ReadWord(low, high)
		result := int(cpu.a) & int(value)
		if result == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		if (value & 0x8000) > 0 {
			cpu.SetFlags(CPUFlagsNegative)
		} else {
			cpu.ClearFlags(CPUFlagsNegative)
		}
		if (value & 0x4000) > 0 {
			cpu.SetFlags(CPUFlagsOverflow)
		} else {
			cpu.ClearFlags(CPUFlagsOverflow)
		}
	}
}

func (cpu *CPU) lda(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		v := cpu.Read(low)
		cpu.a = (cpu.a & 0xFF00) | uint16(v)
	} else {
		cpu.cyclesUsed++
		cpu.a = cpu.ReadWord(low, high)
	}
	cpu.setZN(cpu.a, cpu.mf == 1)
}

func (cpu *CPU) ldx(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		cpu.x = uint16(cpu.Read(low))
	} else {
		cpu.cyclesUsed++
		cpu.x = cpu.ReadWord(low, high)
	}
	cpu.setZN(cpu.x, cpu.xf == 1)
}

func (cpu *CPU) ldy(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		cpu.y = uint16(cpu.Read(low))
	} else {
		cpu.cyclesUsed++
		cpu.y = cpu.ReadWord(low, high)
	}
	cpu.setZN(cpu.y, cpu.xf == 1)
}

func (cpu *CPU) sta(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		cpu.Write(low, byte(cpu.a))
	} else {
		cpu.cyclesUsed++
		cpu.WriteWord(low, high, cpu.a, false)
	}
}

func (cpu *CPU) stx(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		cpu.Write(low, byte(cpu.x))
	} else {
		cpu.cyclesUsed++
		cpu.WriteWord(low, high, cpu.x, false)
	}
}

func (cpu *CPU) sty(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsIndexRegisterSize) {
		cpu.Write(low, byte(cpu.y))
	} else {
		cpu.cyclesUsed++
		cpu.WriteWord(low, high, cpu.y, false)
	}
}

func (cpu *CPU) stz(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		cpu.Write(low, 0)
	} else {
		cpu.cyclesUsed++
		cpu.WriteWord(low, high, 0, false)
	}
}

func (cpu *CPU) ror(low uint32, high uint32) {
	var carry bool = false
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		if (value & 1) == 0x1 {
			carry = true
		} else {
			carry = false
		}
		result = (int(value) >> 1) | (int(cpu.c) << 7)
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		value := cpu.ReadWord(low, high)
		if (value & 1) == 0x01 {
			carry = true
		} else {
			carry = false
		}
		result = (int(value) >> 1) | (int(cpu.c) << 15)
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
	if carry {
		cpu.SetFlags(CPUFlagsCarry)
	} else {
		cpu.ClearFlags(CPUFlagsCarry)
	}
}

func (cpu *CPU) rol(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		result = (int(cpu.Read(low)) << 1) | int(cpu.c)
		if (result & 0x0100) > 0 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		result = (int(cpu.ReadWord(low, high))<<1 | int(cpu.c))
		if (result & 0x10000) > 0 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) lsr(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		if (value & 0x01) == 0x01 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		result = int(value >> 1)
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		value := cpu.ReadWord(low, high)
		if (value & 0x01) == 0x01 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		result = int(value >> 1)
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) asl(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		result = int(cpu.Read(low)) << 1
		if (result & 0x0100) > 0 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		result = int(cpu.ReadWord(low, high)) << 1
		if (result & 0x10000) > 0 {
			cpu.SetFlags(CPUFlagsCarry)
		} else {
			cpu.ClearFlags(CPUFlagsCarry)
		}
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) inc(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		result = int(cpu.Read(low)) + 1
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		result = int(cpu.ReadWord(low, high)) + 1
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) dec(low uint32, high uint32) {
	var result int = 0
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		result = int(cpu.Read(low)) - 1
		cpu.Write(low, byte(result))
	} else {
		cpu.cyclesUsed += 2
		result = int(cpu.ReadWord(low, high)) - 1
		cpu.WriteWord(low, high, uint16(result), true)
	}
	cpu.setZN(uint16(result), cpu.mf == 1)
}

func (cpu *CPU) tsb(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		if ((cpu.a & 0xFF) & uint16(value)) == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		cpu.Write(low, byte(uint16(value)|(cpu.a&0xFF)))
	} else {
		cpu.cyclesUsed += 2
		value := cpu.ReadWord(low, high)
		if (cpu.a & value) == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		cpu.WriteWord(low, high, value|cpu.a, true)
	}
}

func (cpu *CPU) trb(low uint32, high uint32) {
	if cpu.CheckFlag(CPUFlagsAccumulateRegisterSize) {
		value := cpu.Read(low)
		if ((cpu.a & 0xFF) & uint16(value)) == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		cpu.Write(low, byte(uint16(value) & ^(cpu.a&0xFF)))
	} else {
		cpu.cyclesUsed += 2
		value := cpu.ReadWord(low, high)
		if (cpu.a & value) == 0 {
			cpu.SetFlags(CPUFlagsZero)
		} else {
			cpu.ClearFlags(CPUFlagsZero)
		}
		cpu.WriteWord(low, high, value & ^(cpu.a), true)
	}
}
