package main

import (
	"fmt"
)

// can i create some sort of writing based thing?
var BC, DE, HL, AF, SP, PC reg
var reg_numbers = 6
var ZERO uint8 = 0x80
var NEG uint8 = 0x40
var HCARRY uint8 = 0x20
var CARRY uint8 = 0x10

type reg uint16
type hreg uint8

var memory [0x10000]byte

type ins interface {
	Execute([]byte)
}

func PrintRegs() {
	fmt.Printf("%04x,%04x,%04x,%04x,%04x,%04x\n", BC, DE, HL, AF, SP, PC)
}
func main() {
	fmt.Println("hello World")
	for i := 0; i < 0xffffff; i++ {
		INC8{}.Execute(byte(0x4))
		//PrintRegs()
	}
	//}
}

type NOP struct{}

func (i NOP) Execute(b ...byte) {
	PC++
}

type INC8 struct {
}

func (i INC8) Execute(b ...byte) {
	val := get_reg8(b[0] >> 3)
	val++
	if val == 0 {
		set_flagbit(ZERO)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(ZERO)
		unset_flagbit(HCARRY)
	}
	unset_flagbit(NEG)
	set_reg8(b[0]>>3, val)
	PC++
}

type INC16 struct{}

func (i INC16) Execute(b ...byte) {
	val := get_reg16(b[0] >> 4)
	val++
	set_reg16(b[0]>>4, val)
	PC++
}

type DEC8 struct{}

func (f DEC8) Execute(b ...byte) {
	val := get_reg8(b[0] >> 3)
	val--
	if val == 0 {
		set_flagbit(ZERO)
	}
	if val == 0xff {
		set_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	set_reg8(b[0]>>3, val)
	PC++
}

type DEC16 struct{}

func (i DEC16) Execute(b ...byte) {
	val := get_reg16(b[0] >> 4)
	val--
	set_reg16(b[0]>>4, val)
	PC++
}

type LDRreg struct{}

func (l LDRreg) Execute(b ...byte) {
	dst, src := (b[0]&0x38)>>3, b[0]&0xf
	val := get_reg8(src)
	set_reg8(dst, val)
	PC++
}

type LDRu8 struct{}

func (l LDRu8) Execute(b ...byte) {
	dst := b[0] >> 3
	set_reg8(dst, hreg(b[1]))
	PC += 2
}

type LDRu16 struct{}

func (l LDRu16) Execute(b ...byte) {
	dst := b[0] >> 4
	set_reg16(dst, reg(b[1])<<8+reg(b[2]))
	PC += 3
}

type LDRintoreg struct{}

func (l LDRintoreg) Execute(b ...byte) {
	a := get_reg8(0x7)

	val := b[0] >> 4
	if val == 0x2 {
		HL.Set_Byte(a)
		HL++
	} else if val == 0x3 {
		HL.Set_Byte(a)
		HL--
	} else {
		k := get_reg16(val)
		k.Set_Byte(a)
	}
	PC++
}

type ADD16 struct{}

func (a ADD16) Execute(b ...byte) {
	r := b[0] >> 4
	src := get_reg16(r)
	dst := get_reg16(0x2) //HL
	dst_tmp := dst
	dst += src
	set_reg16(0x2, dst)
	if dst_tmp > dst {
		set_flagbit(CARRY)
	} else {
		unset_flagbit(CARRY)
	}
	unset_flagbit(HCARRY)
	unset_flagbit(NEG)
	PC++
}

type ADD8 struct {
}

func (a ADD8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst_tmp := dst
	dst += src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp > dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	unset_flagbit(NEG)
	PC++
}

type ADDu8 struct{}

func (a ADDu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst_tmp := dst
	dst += src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp > dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	unset_flagbit(NEG)
	PC = PC + 2
}

type ADC8 struct {
}

func (a ADC8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst_tmp := dst
	dst += src + hreg(get_flagbit(HCARRY))
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	unset_flagbit(NEG)
	PC += 2
}

type ADCu8 struct {
}

func (a ADCu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst_tmp := dst
	dst += src + hreg(get_flagbit(HCARRY))
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	unset_flagbit(NEG)
	PC += 2
}

type SUB8 struct{}

func (s SUB8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst_tmp := dst
	dst -= src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC++
}

type SUBu8 struct{}

func (s SUBu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst_tmp := dst
	dst -= src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC += 2
}

type SBC8 struct{}

func (s SBC8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst_tmp := dst
	dst -= src + hreg(get_flagbit(HCARRY))
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp > dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC++
}

type SBCu8 struct{}

func (s SBCu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst_tmp := dst
	dst -= src + hreg(get_flagbit(HCARRY))
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp > dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC += 2
}

type AND8 struct {
}

func (a AND8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst = dst & src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	set_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC++
}

type ANDu8 struct {
}

func (a ANDu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst = dst & src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	set_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC += 2
}

type XOR8 struct{}

func (x XOR8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst = dst ^ src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	unset_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC++
}

type XORu8 struct{}

func (x XORu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst = dst ^ src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	unset_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC += 2
}

type OR8 struct{}

func (o OR8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst = dst | src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	unset_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC++
}

type ORu8 struct{}

func (o ORu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst = dst | src
	set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	unset_flagbit(HCARRY)
	unset_flagbit(CARRY)
	unset_flagbit(NEG)
	PC += 2
}

type CP8 struct{}

func (c CP8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst_tmp := dst
	dst -= src
	//set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC++
}

type CPu8 struct{}

func (c CPu8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := hreg(b[1])
	dst_tmp := dst
	dst -= src
	//set_reg8(0x7, dst)
	if dst == 0 {
		set_flagbit(ZERO)
	} else {
		unset_flagbit(ZERO)
	}
	if dst_tmp < dst {
		set_flagbit(CARRY)
		set_flagbit(HCARRY)
	} else {
		unset_flagbit(CARRY)
		unset_flagbit(HCARRY)
	}
	set_flagbit(NEG)
	PC += 2
}

type PUSH struct{}

func (p PUSH) Execute(b ...byte) {
	val := reg(0)
	if (b[0]&0x30)>>4 != 0x3 {
		val = get_reg16((b[0] & 0x30) >> 4)
	} else {
		val = AF
	}
	SP = SP - 2
	SP.Set_Word(val)
	PC++
}

type POP struct{}

func (p POP) Execute(b ...byte) {
	val := SP.Get_Word()
	SP = SP + 2
	if (b[0]&0x30)>>4 != 0x3 {
		set_reg16((b[0]&0x30)>>4, val)
	} else {
		AF = val
	}
	PC++
}

type CALL struct{}

func (c CALL) Execute(b ...byte) {
	if b[0]&1 == 1 || jmp_on_flag((b[0]&0x18)>>3) {
		SP = SP - 2
		SP.Set_Word(PC + 3)
		PC = reg(b[1])<<8 + reg(b[2])
	} else {
		PC = PC + 3
	}
}

type JP struct{}

func (j JP) Execute(b ...byte) {
	if b[0]&0x21 == 0x21 {
		PC = HL
	} else if b[0]&0x21 == 0x1 || jmp_on_flag((b[0]&0x18)>>3) {
		PC = reg(b[1])<<8 + reg(b[2])
	} else {
		PC = PC + 3
	}
}

type JR struct{}

func (j JR) Execute(b ...byte) {
	if b[0]&0x20 == 0x0 || jmp_on_flag((b[0] & 0x18)) {
		if int8(b[1]) < 0 {
			PC -= reg((b[1] ^ 0xff) + 1)
		} else {
			PC += reg(b[1])
		}
	} else {
		PC += 2
	}
}

type RET struct{}

// This handles all return types except for RETI, which will implemented later after the emulation period.
func (r RET) Execute(b ...byte) {
	if b[0]&0x1 == 1 || jmp_on_flag(b[0]&0x18) {
		PC = SP.Get_Word()
		SP = SP + 2
	} else {
		PC++
	}
}

type RETI struct{}

func (r RETI) Execute(b ...byte) {
	panic("NOT IMPLEMENTED")
}

type RCCA struct{}

func (r RCCA) Execute(b ...byte) {
	a := get_reg8(0x7)
	c := a & 0x1
	a = (a >> 1) + c<<7
	set_reg8(0x7, a)
	unset_flagbit(ZERO)
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	if c == 0 {
		unset_flagbit(CARRY)
	} else {
		set_flagbit(CARRY)
	}
	PC++
}

type RRA struct{}

func (r RRA) Execute(b ...byte) {
	c := hreg(get_flagbit(CARRY))
	a := get_reg8(0x7)
	c_new := a & 0x1
	a = (a >> 1) + c<<7
	set_reg8(0x7, a)
	unset_flagbit(ZERO)
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	if c_new == 0 {
		unset_flagbit(CARRY)
	} else {
		set_flagbit(CARRY)
	}
	PC++

}

type RLCA struct{}

func (r RLCA) Execute(b ...byte) {
	a := get_reg8(0x7)
	c := (a & 0x80) >> 7
	a = (a << 1) + c
	set_reg8(0x7, a)
	unset_flagbit(ZERO)
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	if c == 0 {
		unset_flagbit(CARRY)
	} else {
		set_flagbit(CARRY)
	}
	PC++
}

type RLA struct{}

func (r RLA) Execute(b ...byte) {
	c := hreg(get_flagbit(CARRY))
	a := get_reg8(0x7)
	c_new := (a & 0x80) >> 7
	a = (a << 1) + c
	set_reg8(0x7, a)
	unset_flagbit(ZERO)
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	if c_new == 0 {
		unset_flagbit(CARRY)
	} else {
		set_flagbit(CARRY)
	}
	PC++

}

type HALT struct{}

func (h HALT) Execute(b ...byte) {
	panic("NOT IMPLEMENTED YET")
}

type DAA struct{}

func (d DAA) Execute(b ...byte) {
	a := get_reg8(0x7)
	if !is_flagbit(NEG) { // after an addition, adjust if (half-)carry occurred or if result is out of bounds
		if is_flagbit(CARRY) || a > 0x99 {
			a += 0x60
		}
		if is_flagbit(HCARRY) || (a&0x0f) > 0x09 {
			a += 0x6
		}
	} else { // after a subtraction, only adjust if (half-)carry occurred
		if is_flagbit(CARRY) {
			a -= 0x60
		}
		if is_flagbit(HCARRY) {
			a -= 0x6
		}
	}
	set_reg8(0x7, a)
	set_flagbit(CARRY)
	set_flagbit(NEG)
	PC++
}

type CPL struct{}

func (c CPL) Execute(b ...byte) {
	a := get_reg8(0x7)
	a = a ^ 0xff
	set_reg8(0x7, a)
	set_flagbit(HCARRY)
	set_flagbit(NEG)
	PC++
}

type SCF struct{}

func (s SCF) Execute(b ...byte) {
	set_flagbit(CARRY)
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	PC++
}

type CCF struct{}

func (c CCF) Execute(b ...byte) {
	if is_flagbit(CARRY) {
		unset_flagbit(CARRY)
	} else {
		set_flagbit(CARRY)
	}
	unset_flagbit(NEG)
	unset_flagbit(HCARRY)
	PC++
}

type RST struct{}

//This is basically the same as CALL except just one byte, and a few locations
func (r RST) Execute(b ...byte) {
	addr := reg(b[0] & 0x38)
	SP.Set_Word(PC + 1)
	SP = SP - 2
	PC = addr
}

type DI struct{}

func (d DI) Execute(b ...byte) {
	panic("SHOULD DISABLE INTERRUPTS")
}

type STOP struct{}

func (s STOP) Execute(b ...byte) {
	panic("STOP")
}

type EI struct{}

func (e EI) Execute(b ...byte) {
	panic("SHOULD ENABLE INTERRUPTS")
}

func jmp_on_flag(index uint8) bool {
	switch index {
	case 0:
		return !is_flagbit(ZERO)
	case 1:
		return is_flagbit(ZERO)
	case 2:
		return !is_flagbit(CARRY)
	case 3:
		return is_flagbit(CARRY)
	default:
		panic("JMP ON FLAG PANIC")
	}

}

func get_reg8(index uint8) hreg {
	ret_byte := reg(0)
	switch index {
	case 0:
		ret_byte = BC >> 8
	case 1:
		ret_byte = BC & 0xff
	case 2:
		ret_byte = DE >> 8
	case 3:
		ret_byte = DE & 0xff
	case 4:
		ret_byte = HL >> 8
	case 5:
		ret_byte = HL & 0xff
	case 6:
		ret_byte = reg(HL.Get_Byte())
	case 7:
		ret_byte = AF >> 8
	default:
		panic("OH SHIT")
	}
	return hreg(ret_byte)
}

func set_reg8(index uint8, v hreg) {
	val := reg(v)
	switch index {
	case 0:
		BC = BC&0xff + val<<8
	case 1:
		BC = BC&0xff00 + val
	case 2:
		DE = DE&0xff + val<<8
	case 3:
		DE = DE&0xff00 + val
	case 4:
		HL = HL&0xff + val<<8
	case 5:
		HL = HL&0xff00 + val
	case 6:
		HL.Set_Byte(hreg(val))
	case 7:
		AF = AF&0xff + val<<8
	default:
		panic("OH SHIT")
	}
}

func set_flag(flag hreg) {
	AF = AF&0xff00 + reg(flag)
}
func set_flagbit(i uint8) {
	AF = AF | reg(i)
}
func unset_flagbit(i uint8) {
	AF = AF & (reg(i) ^ 0xffff)
}
func get_flagbit(i uint8) uint8 {
	if AF&reg(i) >= 1 {
		return 1
	}
	return 0
}

func is_flagbit(i uint8) bool {
	return get_flagbit(i) == 1
}

func get_flag() hreg {
	return hreg(AF & 0xff)
}

func get_reg16(index uint8) reg {
	switch index {
	case 0:
		return BC
	case 1:
		return DE
	case 2:
		return HL
	case 3:
		return SP
	default:
		panic("OH SHIT")
	}
}
func set_reg16(index uint8, val reg) {
	switch index {
	case 0:
		BC = val
	case 1:
		DE = val
	case 2:
		HL = val
	case 3:
		SP = val
	default:
		panic("OH SHIT")
	}
}

// Stubbed until later
func (r reg) Get_Byte() hreg {
	return hreg(memory[int(r)])
}
func (r reg) Get_Word() reg {
	val := reg(memory[int(r)])<<8 + reg(memory[int(r)+1])
	return val
}
func (r reg) Set_Byte(a hreg) {
	memory[int(r)] = byte(a)
}
func (r reg) Set_Word(a reg) {
	memory[int(r)] = byte((a & 0xff00) >> 8)
	memory[int(r)+1] = byte((a & 0xff))
}
