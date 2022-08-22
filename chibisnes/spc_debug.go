package chibisnes

import (
	"fmt"
)

var opcodeNamesSpc [256]string = [256]string{
	"nop              ", "tcall 0          ", "set1 $%02x.0       ", "bbs $%02x.0, $%04x ", "or a, $%02x        ", "or a, $%04x      ", "or a, [X]        ", "or a, [$%02x+x]    ", "or a, #$%02x       ", "or $%02x, $%02x      ", "or1 c, $%04x.%01x   ", "asl $%02x          ", "asl $%04x        ", "push p           ", "tset $%04x       ", "brk              ",
	"bpl $%04x        ", "tcall 1          ", "clr1 $%02x.0       ", "bbc $%02x.0, $%04x ", "or a, $%02x+x      ", "or a, $%04x+x    ", "or a, $%04x+y    ", "or a, [$%02x]+y    ", "or $%02x, #$%02x     ", "or [X], [Y]      ", "decw $%02x         ", "asl $%02x+x        ", "asl a            ", "dec x            ", "cmp x, $%04x     ", "jmp [$%04x+x]    ",
	"clrp             ", "tcall 2          ", "set1 $%02x.1       ", "bbs $%02x.1, $%04x ", "and a, $%02x       ", "and a, $%04x     ", "and a, [X]       ", "and a, [$%02x+x]   ", "and a, #$%02x      ", "and $%02x, $%02x     ", "or1 c, /$%04x.%01x  ", "rol $%02x          ", "rol $%04x        ", "push a           ", "cbne $%02x, $%04x  ", "bra $%04x        ",
	"bmi $%04x        ", "tcall 3          ", "clr1 $%02x.1       ", "bbc $%02x.1, $%04x ", "and a, $%02x+x     ", "and a, $%04x+x   ", "and a, $%04x+y   ", "and a, [$%02x]+y   ", "and $%02x, #$%02x    ", "and [X], [Y]     ", "incw $%02x         ", "rol $%02x+x        ", "rol a            ", "inc x            ", "cmp x, $%02x       ", "call $%04x       ",
	"setp             ", "tcall 4          ", "set1 $%02x.2       ", "bbs $%02x.2, $%04x ", "eor a, $%02x       ", "eor a, $%04x     ", "eor a, [X]       ", "eor a, [$%02x+x]   ", "eor a, #$%02x      ", "eor $%02x, $%02x     ", "and1 c, $%04x.%01x  ", "lsr $%02x          ", "lsr $%04x        ", "push x           ", "tclr $%04x       ", "pcall $%02x        ",
	"bvc $%04x        ", "tcall 5          ", "clr1 $%02x.2       ", "bbc $%02x.2, $%04x ", "eor a, $%02x+x     ", "eor a, $%04x+x   ", "eor a, $%04x+y   ", "eor a, [$%02x]+y   ", "eor $%02x, #$%02x    ", "eor [X], [Y]     ", "cmpw ya, $%02x     ", "lsr $%02x+x        ", "lsr a            ", "mov x, a         ", "cmp y, $%04x     ", "jmp $%04x        ",
	"clrc             ", "tcall 6          ", "set1 $%02x.3       ", "bbs $%02x.3, $%04x ", "cmp a, $%02x       ", "cmp a, $%04x     ", "cmp a, [X]       ", "cmp a, [$%02x+x]   ", "cmp a, #$%02x      ", "cmp $%02x, $%02x     ", "and1 c, /$%04x.%01x ", "ror $%02x          ", "ror $%04x        ", "push y           ", "dbnz $%02x, $%04x  ", "ret              ",
	"bvs $%04x        ", "tcall 7          ", "clr1 $%02x.3       ", "bbc $%02x.3, $%04x ", "cmp a, $%02x+x     ", "cmp a, $%04x+x   ", "cmp a, $%04x+y   ", "cmp a, [$%02x]+y   ", "cmp $%02x, #$%02x    ", "cmp [X], [Y]     ", "addw ya, $%02x     ", "ror $%02x+x        ", "ror a            ", "mov a, x         ", "cmp y, $%02x       ", "reti             ",
	"setc             ", "tcall 8          ", "set1 $%02x.4       ", "bbs $%02x.4, $%04x ", "adc a, $%02x       ", "adc a, $%04x     ", "adc a, [X]       ", "adc a, [$%02x+x]   ", "adc a, #$%02x      ", "adc $%02x, $%02x     ", "eor1 c, $%04x.%01x  ", "dec $%02x          ", "dec $%04x        ", "mov y, #$%02x      ", "pop p            ", "mov $%02x, #$%02x    ",
	"bcc $%04x        ", "tcall 9          ", "clr1 $%02x.4       ", "bbc $%02x.4, $%04x ", "adc a, $%02x+x     ", "adc a, $%04x+x   ", "adc a, $%04x+y   ", "adc a, [$%02x]+y   ", "adc $%02x, #$%02x    ", "adc [X], [Y]     ", "subw ya, $%02x     ", "dec $%02x+x        ", "dec a            ", "mov x, sp        ", "div ya, x        ", "xcn a            ",
	"ei               ", "tcall 10         ", "set1 $%02x.5       ", "bbs $%02x.5, $%04x ", "sbc a, $%02x       ", "sbc a, $%04x     ", "sbc a, [X]       ", "sbc a, [$%02x+x]   ", "sbc a, #$%02x      ", "sbc $%02x, $%02x     ", "mov1 c, $%04x.%01x  ", "inc $%02x          ", "inc $%04x        ", "cmp y, #$%02x      ", "pop a            ", "mov [x+], a      ",
	"bcs $%04x        ", "tcall 11         ", "clr1 $%02x.5       ", "bbc $%02x.5, $%04x ", "sbc a, $%02x+x     ", "sbc a, $%04x+x   ", "sbc a, $%04x+y   ", "sbc a, [$%02x]+y   ", "sbc $%02x, #$%02x    ", "sbc [X], [Y]     ", "movw ya, $%02x     ", "inc $%02x+x        ", "inc a            ", "mov sp, x        ", "das a            ", "mov a, [x+]      ",
	"di               ", "tcall 12         ", "set1 $%02x.6       ", "bbs $%02x.6, $%04x ", "mov $%02x, a       ", "mov $%04x, a     ", "mov [X], a       ", "mov [$%02x+x], a   ", "cmp x, #$%02x      ", "mov $%04x, x     ", "mov1 $%04x.%01x, c  ", "mov $%02x, y       ", "mov $%04x, y     ", "mov x, #$%02x      ", "pop x            ", "mul ya           ",
	"bne $%04x        ", "tcall 13         ", "clr1 $%02x.6       ", "bbc $%02x.6, $%04x ", "mov $%02x+x, a     ", "mov $%04x+x, a   ", "mov $%04x+y, a   ", "mov [$%02x]+y, a   ", "mov $%02x, x       ", "mov $%02x+y, x     ", "movw $%02x, ya     ", "mov $%02x+x, y     ", "dec y            ", "mov a, y         ", "cbne $%02x+x, $%04x", "daa a            ",
	"clrv             ", "tcall 14         ", "set1 $%02x.7       ", "bbs $%02x.7, $%04x ", "mov a, $%02x       ", "mov a, $%04x     ", "mov a, [X]       ", "mov a, [$%02x+x]   ", "mov a, #$%02x      ", "mov x, $%04x     ", "not1 $%04x.%01x     ", "mov y, $%02x       ", "mov y, $%04x     ", "notc             ", "pop y            ", "sleep            ",
	"beq $%04x        ", "tcall 15         ", "clr1 $%02x.7       ", "bbc $%02x.7, $%04x ", "mov a, $%02x+x     ", "mov a, $%04x+x   ", "mov a, $%04x+y   ", "mov a, [$%02x]+y   ", "mov x, $%02x       ", "mov x, $%02x+y     ", "mov $%02x, $%02x     ", "mov y, $%02x+x     ", "inc y            ", "mov y, a         ", "dbnz y, $%04x    ", "stop             ",
}

// address types for each opcode, for spc
var opcodeTypeSpc [256]int = [256]int{
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 0, 2, 0,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 2, 2,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 0, 5, 3,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 1, 2,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 0, 2, 1,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 2, 2,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 0, 5, 0,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 1, 0,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 1, 0, 4,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 0, 0,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 4, 6, 1, 2, 1, 0, 0,
	3, 0, 1, 5, 1, 2, 2, 1, 4, 0, 1, 1, 0, 0, 0, 0,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 2, 6, 1, 2, 1, 0, 0,
	3, 0, 1, 5, 1, 2, 2, 1, 1, 1, 1, 1, 0, 0, 5, 0,
	0, 0, 1, 5, 1, 2, 0, 1, 1, 2, 6, 1, 2, 0, 0, 0,
	3, 0, 1, 5, 1, 2, 2, 1, 1, 1, 4, 1, 0, 0, 3, 0,
}

func (spc *SPC) getProcessorStateSPC() string {
	var nChar, vChar, pChar, bChar, hChar, iChar, zChar, cChar string
	if spc.n == 0x01 {
		nChar = "N"
	} else {
		nChar = "n"
	}
	if spc.v == 0x01 {
		vChar = "V"
	} else {
		vChar = "v"
	}
	if spc.p == 0x01 {
		pChar = "P"
	} else {
		pChar = "p"
	}
	if spc.b == 0x01 {
		bChar = "B"
	} else {
		bChar = "b"
	}
	if spc.h == 0x01 {
		hChar = "H"
	} else {
		hChar = "h"
	}
	if spc.i == 0x01 {
		iChar = "I"
	} else {
		iChar = "i"
	}
	if spc.z == 0x01 {
		zChar = "Z"
	} else {
		zChar = "z"
	}
	if spc.c == 0x01 {
		cChar = "C"
	} else {
		cChar = "c"
	}

	disLine := spc.getDisassemblySPC()

	return fmt.Sprintf(
		"SPC %04x %s A:%02x X:%02x Y:%02x SP:%02x %s%s%s%s%s%s%s%s",
		spc.pc, disLine, spc.a, spc.x, spc.y, spc.sp,
		nChar,
		vChar,
		pChar,
		bChar,
		hChar,
		iChar,
		zChar,
		cChar,
	)
}

func (spc *SPC) getDisassemblySPC() string {
	var addr uint16 = spc.pc
	// read 3 bytes
	// TODO: this can have side effects, implement and use peaking
	var opcode byte = spc.apu.Read(addr)
	var byte1 byte = spc.apu.Read((addr + 1) & 0xffff)
	var byte2 byte = spc.apu.Read((addr + 2) & 0xffff)
	var word uint16 = (uint16(byte2) << 8) | uint16(byte1)
	var rel uint16 = uint16(int16(spc.pc) + 2 + int16(int8(byte1)))
	var rel2 uint16 = uint16(int16(spc.pc) + 2 + int16(int8(byte2)))
	var wordb uint16 = word & 0x1fff
	var bit byte = byte(word >> 13)
	// switch on type
	switch opcodeTypeSpc[opcode] {
	case 0:
		return fmt.Sprintf("%s", opcodeNamesSpc[opcode])
	case 1:
		return fmt.Sprintf(opcodeNamesSpc[opcode], byte1)
	case 2:
		return fmt.Sprintf(opcodeNamesSpc[opcode], word)
	case 3:
		return fmt.Sprintf(opcodeNamesSpc[opcode], rel)
	case 4:
		return fmt.Sprintf(opcodeNamesSpc[opcode], byte2, byte1)
	case 5:
		return fmt.Sprintf(opcodeNamesSpc[opcode], byte1, rel2)
	case 6:
		return fmt.Sprintf(opcodeNamesSpc[opcode], wordb, bit)
	}

	return ""
}
