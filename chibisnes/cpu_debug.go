package chibisnes

import (
	"fmt"
)

var opcodeNames [256]string = [256]string{
	"brk          ", "ora ($%02x,x)  ", "cop #$%02x     ", "ora $%02x,s    ", "tsb $%02x      ", "ora $%02x      ", "asl $%02x      ", "ora [$%02x]    ", "php          ", "ora #$%04x   ", "asl          ", "phd          ", "tsb $%04x    ", "ora $%04x    ", "asl $%04x    ", "ora $%06x  ",
	"bpl $%04x    ", "ora ($%02x),y  ", "ora ($%02x)    ", "ora ($%02x,s),y", "trb $%02x      ", "ora $%02x,x    ", "asl $%02x,x    ", "ora [$%02x],y  ", "clc          ", "ora $%04x,y  ", "inc          ", "tcs          ", "trb $%04x    ", "ora $%04x,x  ", "asl $%04x,x  ", "ora $%06x,x",
	"jsr $%04x    ", "and ($%02x,x)  ", "jsl $%06x  ", "and $%02x,s    ", "bit $%02x      ", "and $%02x      ", "rol $%02x      ", "and [$%02x]    ", "plp          ", "and #$%04x   ", "rol          ", "pld          ", "bit $%04x    ", "and $%04x    ", "rol $%04x    ", "and $%06x  ",
	"bmi $%04x    ", "and ($%02x),y  ", "and ($%02x)    ", "and ($%02x,s),y", "bit $%02x,x    ", "and $%02x,x    ", "rol $%02x,x    ", "and [$%02x],y  ", "sec          ", "and $%04x,y  ", "dec          ", "tsc          ", "bit $%04x,x  ", "and $%04x,x  ", "rol $%04x,x  ", "and $%06x,x",
	"rti          ", "eor ($%02x,x)  ", "wdm #$%02x     ", "eor $%02x,s    ", "mvp $%02x, $%02x ", "eor $%02x      ", "lsr $%02x      ", "eor [$%02x]    ", "pha          ", "eor #$%04x   ", "lsr          ", "phk          ", "jmp $%04x    ", "eor $%04x    ", "lsr $%04x    ", "eor $%06x  ",
	"bvc $%04x    ", "eor ($%02x),y  ", "eor ($%02x)    ", "eor ($%02x,s),y", "mvn $%02x, $%02x ", "eor $%02x,x    ", "lsr $%02x,x    ", "eor [$%02x],y  ", "cli          ", "eor $%04x,y  ", "phy          ", "tcd          ", "jml $%06x  ", "eor $%04x,x  ", "lsr $%04x,x  ", "eor $%06x,x",
	"rts          ", "adc ($%02x,x)  ", "per $%04x    ", "adc $%02x,s    ", "stz $%02x      ", "adc $%02x      ", "ror $%02x      ", "adc [$%02x]    ", "pla          ", "adc #$%04x   ", "ror          ", "rtl          ", "jmp ($%04x)  ", "adc $%04x    ", "ror $%04x    ", "adc $%06x  ",
	"bvs $%04x    ", "adc ($%02x),y  ", "adc ($%02x)    ", "adc ($%02x,s),y", "stz $%02x,x    ", "adc $%02x,x    ", "ror $%02x,x    ", "adc [$%02x],y  ", "sei          ", "adc $%04x,y  ", "ply          ", "tdc          ", "jmp ($%04x,x)", "adc $%04x,x  ", "ror $%04x,x  ", "adc $%06x,x",
	"bra $%04x    ", "sta ($%02x,x)  ", "brl $%04x    ", "sta $%02x,s    ", "sty $%02x      ", "sta $%02x      ", "stx $%02x      ", "sta [$%02x]    ", "dey          ", "bit #$%04x   ", "txa          ", "phb          ", "sty $%04x    ", "sta $%04x    ", "stx $%04x    ", "sta $%06x  ",
	"bcc $%04x    ", "sta ($%02x),y  ", "sta ($%02x)    ", "sta ($%02x,s),y", "sty $%02x,x    ", "sta $%02x,x    ", "stx $%02x,y    ", "sta [$%02x],y  ", "tya          ", "sta $%04x,y  ", "txs          ", "txy          ", "stz $%04x    ", "sta $%04x,x  ", "stz $%04x,x  ", "sta $%06x,x",
	"ldy #$%04x   ", "lda ($%02x,x)  ", "ldx #$%04x   ", "lda $%02x,s    ", "ldy $%02x      ", "lda $%02x      ", "ldx $%02x      ", "lda [$%02x]    ", "tay          ", "lda #$%04x   ", "tax          ", "plb          ", "ldy $%04x    ", "lda $%04x    ", "ldx $%04x    ", "lda $%06x  ",
	"bcs $%04x    ", "lda ($%02x),y  ", "lda ($%02x)    ", "lda ($%02x,s),y", "ldy $%02x,x    ", "lda $%02x,x    ", "ldx $%02x,y    ", "lda [$%02x],y  ", "clv          ", "lda $%04x,y  ", "tsx          ", "tyx          ", "ldy $%04x,x  ", "lda $%04x,x  ", "ldx $%04x,y  ", "lda $%06x,x",
	"cpy #$%04x   ", "cmp ($%02x,x)  ", "rep #$%02x     ", "cmp $%02x,s    ", "cpy $%02x      ", "cmp $%02x      ", "dec $%02x      ", "cmp [$%02x]    ", "iny          ", "cmp #$%04x   ", "dex          ", "wai          ", "cpy $%04x    ", "cmp $%04x    ", "dec $%04x    ", "cmp $%06x  ",
	"bne $%04x    ", "cmp ($%02x),y  ", "cmp ($%02x)    ", "cmp ($%02x,s),y", "pei $%02x      ", "cmp $%02x,x    ", "dec $%02x,x    ", "cmp [$%02x],y  ", "cld          ", "cmp $%04x,y  ", "phx          ", "stp          ", "jml [$%04x]  ", "cmp $%04x,x  ", "dec $%04x,x  ", "cmp $%06x,x",
	"cpx #$%04x   ", "sbc ($%02x,x)  ", "sep #$%02x     ", "sbc $%02x,s    ", "cpx $%02x      ", "sbc $%02x      ", "inc $%02x      ", "sbc [$%02x]    ", "inx          ", "sbc #$%04x   ", "nop          ", "xba          ", "cpx $%04x    ", "sbc $%04x    ", "inc $%04x    ", "sbc $%06x  ",
	"beq $%04x    ", "sbc ($%02x),y  ", "sbc ($%02x)    ", "sbc ($%02x,s),y", "pea #$%04x   ", "sbc $%02x,x    ", "inc $%02x,x    ", "sbc [$%02x],y  ", "sed          ", "sbc $%04x,y  ", "plx          ", "xce          ", "jsr ($%04x,x)", "sbc $%04x,x  ", "inc $%04x,x  ", "sbc $%06x,x",
}

var opcodeNamesSp [256]string = [256]string{
	"", "", "", "", "", "", "", "", "", "ora #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "and #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "eor #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "adc #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "bit #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"ldy #$%02x     ", "", "ldx #$%02x     ", "", "", "", "", "", "", "lda #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"cpy #$%02x     ", "", "", "", "", "", "", "", "", "cmp #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"cpx #$%02x     ", "", "", "", "", "", "", "", "", "sbc #$%02x     ", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
}

var opcodeType [256]int = [256]int{
	0, 1, 1, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	2, 1, 3, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	0, 1, 1, 1, 8, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 8, 1, 1, 1, 0, 2, 0, 0, 3, 2, 2, 3,
	0, 1, 7, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	6, 1, 7, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	5, 1, 5, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	5, 1, 1, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 1, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
	5, 1, 1, 1, 1, 1, 1, 1, 0, 4, 0, 0, 2, 2, 2, 3,
	6, 1, 1, 1, 2, 1, 1, 1, 0, 2, 0, 0, 2, 2, 2, 3,
}

func (cpu *CPU) getProcessorStateCPU() string {
	var eChar, nChar, vChar, mfChar, xfChar, dChar, iChar, zChar, cChar string
	if cpu.e == 0x01 {
		eChar = "E"
	} else {
		eChar = "e"
	}
	if cpu.n == 0x01 {
		nChar = "N"
	} else {
		nChar = "n"
	}
	if cpu.v == 0x01 {
		vChar = "V"
	} else {
		vChar = "v"
	}
	if cpu.mf == 0x01 {
		mfChar = "M"
	} else {
		mfChar = "m"
	}
	if cpu.xf == 0x01 {
		xfChar = "X"
	} else {
		xfChar = "x"
	}
	if cpu.d == 0x01 {
		dChar = "D"
	} else {
		dChar = "d"
	}
	if cpu.i == 0x01 {
		iChar = "I"
	} else {
		iChar = "i"
	}
	if cpu.z == 0x01 {
		zChar = "Z"
	} else {
		zChar = "z"
	}
	if cpu.c == 0x01 {
		cChar = "C"
	} else {
		cChar = "c"
	}

	disLine := cpu.getDisassemblyCPU()

	return fmt.Sprintf(
		"CPU %02x:%04x %-11s A:%04x X:%04x Y:%04x SP:%04x DP:%04x DB:%02x %s %s%s%s%s%s%s%s%s",
		cpu.k, cpu.pc, disLine, cpu.a, cpu.x, cpu.y,
		cpu.sp, cpu.dp, cpu.db,
		eChar,
		nChar,
		vChar,
		mfChar,
		xfChar,
		dChar,
		iChar,
		zChar,
		cChar,
	)
}

func (cpu *CPU) getDisassemblyCPU() string {
	var addr uint32 = uint32(cpu.pc) | (uint32(cpu.k) << 16)
	opcode := cpu.console.Read(addr)
	byte1 := cpu.console.Read((addr + 1) & 0xFFFFFF)
	byte2 := cpu.console.Read((addr + 2) & 0xFFFFFF)
	word := uint16(byte2)<<8 | uint16(byte1)
	longv := (uint32(cpu.console.Read(((addr+3)&0xFFFFFF)))<<16 | uint32(word))
	rel := uint16(int16(cpu.pc) + 2 + int16(int8(byte1)))
	rell := uint16(int16(cpu.pc) + 3 + int16(word))

	switch opcodeType[opcode] {
	case 0:
		return fmt.Sprintf("%s", opcodeNames[opcode])
	case 1:
		return fmt.Sprintf(opcodeNames[opcode], byte1)
	case 2:
		return fmt.Sprintf(opcodeNames[opcode], word)
	case 3:
		return fmt.Sprintf(opcodeNames[opcode], longv)
	case 4:
		if cpu.mf == 0x01 {
			return fmt.Sprintf(opcodeNamesSp[opcode], byte1)
		} else {
			return fmt.Sprintf(opcodeNames[opcode], word)
		}
	case 5:
		if cpu.xf == 0x01 {
			return fmt.Sprintf(opcodeNamesSp[opcode], byte1)
		} else {
			return fmt.Sprintf(opcodeNames[opcode], word)
		}
	case 6:
		return fmt.Sprintf(opcodeNames[opcode], rel)
	case 7:
		return fmt.Sprintf(opcodeNames[opcode], rell)
	case 8:
		return fmt.Sprintf(opcodeNames[opcode], byte2, byte1)
	}

	return ""
}
