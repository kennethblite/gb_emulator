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
	for i := 0; i < 0xffffff; i++{
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
	PC++
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

type AND8 struct{
}

func (a AND8) Execute(b ...byte) {
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst =  dst&src
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

type XOR8 struct{}
func (x XOR8) Execute(b ...byte){
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst =  dst^src
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

type OR8 struct{}
func (o OR8) Execute(b ...byte){
	dst := get_reg8(0x7) //A
	src := get_reg8(b[0] & 0x7)
	dst =  dst|src
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

type CP8 struct{}
func (c CP8) Execute(b ...byte){
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

func set_statusflags() {

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
func (r reg) Set_Byte(a hreg) {
	memory[int(r)] = byte(a)
}
