package cpu

import (
	"fmt"
)

// Status used by P (status) registers
type Status uint8

const (
	// C carry flag
	C Status = 1 << iota
	// Z zero flag
	Z
	// I interrupt disable
	I
	// D decimal mode
	D
	// B break command
	B
	// U UNUSED
	U
	// V overflow flag
	V
	// N negative flag
	N
)

// Registers are 6502 CPU registers
type Registers struct {
	A  uint8  // accumulator
	X  uint8  // index register X
	Y  uint8  // index register Y
	P  Status // processor status
	SP uint8  // stack pointer
	PC uint16 // program counter
}

// NewRegisters creates a new set of registers, all set to 0
func NewRegisters() (reg Registers) {
	reg = Registers{}
	reg.Reset()
	return
}

// Reset resets all registers. P is initially only I bit, SP is initially 0xfd,
// PC is initially 0xfffc (reset vector) and all others are 0
func (reg *Registers) Reset() {
	reg.A = 0
	reg.X = 0
	reg.Y = 0
	reg.P = I | U
	reg.SP = 0xfd
	reg.PC = 0xfffc
}

type Interrupt uint8

const (
	// Irq interrupt request
	Irq Interrupt = iota
	// Nmi non-maskable interrupt
	Nmi
	// Rst reset
	Rst
)

type Index uint8

const (
	X Index = iota
	Y
)

// CPU represents a 6502 CPU
type CPU struct {
	Registers    Registers
	Memory       Memory
	Instructions InstructionTable
	decimalMode  bool
	breakError   bool
	Irq          bool
	Nmi          bool
	Rst          bool
}

func NewCPU(mem Memory) *CPU {
	instructions := NewInstructionTable()
	instructions.InitInstructions()

	return &CPU{
		Registers:    NewRegisters(),
		Memory:       mem,
		Instructions: instructions,
		decimalMode:  true,
		breakError:   false,
	}
}

func (cpu *CPU) Reset() {
	cpu.Registers.Reset()
	cpu.Memory.Reset()
	cpu.ExecuteRst()
}

func (cpu *CPU) IndexToRegister(which Index) uint8 {
	var index uint8
	switch which {
	case X:
		index = cpu.Registers.X
	case Y:
		index = cpu.Registers.Y
	}
	return index
}

type BadOpCodeError OpCode

func (b BadOpCodeError) Error() string {
	return fmt.Sprintf("No such opcode %#02x", b)
}

type BrkOpCodeError OpCode

func (b BrkOpCodeError) Error() string {
	return fmt.Sprintf("Executed BRK opcode")
}

// Execute takes instruction of PC and executes it in the number
// of cycles as returned by the instruction's Exec function.
// Returns the number of cycles executed and any error, if any.
func (cpu *CPU) Execute() (cycles uint16, err error) {
	cycles += cpu.ExecuteInterrupt()

	// fetch
	opcode := OpCode(cpu.Memory.Read(cpu.Registers.PC))
	inst := cpu.Instructions.opcodes[opcode]
	if inst == nil {
		return 0, BadOpCodeError(opcode)
	}

	// execute
	cpu.Registers.PC++
	cycles += cpu.Instructions.Execute(cpu, opcode)

	if cpu.breakError && opcode == 0x00 {
		return cycles, BrkOpCodeError(opcode)
	}

	return cycles, nil
}

// Run executes instruction until Execute() returns an error
func (cpu *CPU) Run() (err error) {
	for {
		if _, err = cpu.Execute(); err != nil {
			return
		}
	}
}

// Interrupts
// ==========

func (cpu *CPU) ExecuteInterrupt() (cycles uint16) {
	cycles = 7
	switch {
	case cpu.Irq && cpu.Registers.P&I == 0:
		cpu.ExecuteIrq()
		cpu.Irq = false
	case cpu.Nmi:
		cpu.ExecuteNmi()
		cpu.Nmi = false
	case cpu.Rst:
		cpu.ExecuteRst()
		cpu.Rst = false
	default:
		cycles = 0
	}
	return
}

func (cpu *CPU) ExecuteIrq() {
	cpu.push16(cpu.Registers.PC)
	cpu.push8(uint8((cpu.Registers.P | U) & ^B))
	cpu.Registers.P |= I

	low := cpu.Memory.Read(0xfffe)
	high := cpu.Memory.Read(0xffff)
	cpu.Registers.PC = (uint16(high)<<8 | uint16(low))
}

func (cpu *CPU) ExecuteNmi() {
	cpu.push16(cpu.Registers.PC)
	cpu.push8(uint8((cpu.Registers.P | U) & ^B))
	cpu.Registers.P |= I

	low := cpu.Memory.Read(0xfffa)
	high := cpu.Memory.Read(0xfffb)
	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
}

func (cpu *CPU) ExecuteRst() {
	low := cpu.Memory.Read(0xfffc)
	high := cpu.Memory.Read(0xfffd)
	cpu.Registers.PC = (uint16(high) << 8) | uint16(low)
}

// Addressing Modes
// ================

func (cpu *CPU) aluAddress(opcode OpCode, status *InstructionStatus) (address uint16) {
	// alu opcodes end with 01
	if opcode&0x10 == 0 {
		switch (opcode >> 2) & 0x03 {
		case 0x00:
			address = cpu.indexedIndirectAddress()
		case 0x01:
			address = cpu.zeroPageAddress()
		case 0x02:
			address = cpu.immediateAddress()
		case 0x03:
			address = cpu.absoluteAddress()
		}

	} else {
		switch (opcode >> 2) & 0x03 {
		case 0x00:
			address = cpu.indirectIndexedAddress(status)
		case 0x01:
			address = cpu.indexedZeroPageAddress(X)
		case 0x02:
			address = cpu.indexedAbsoluteAddress(Y, status)
		case 0x03:
			address = cpu.indexedAbsoluteAddress(X, status)
		}
	}

	return
}

// read-modify-write instructions
func (cpu *CPU) rmwAddress(opcode OpCode, status *InstructionStatus) (address uint16) {
	// rmw opcodes end with 10
	if opcode&0x10 == 0 {
		switch (opcode >> 2) & 0x03 {
		case 0x00:
			address = cpu.immediateAddress()
		case 0x01:
			address = cpu.zeroPageAddress()
		case 0x02:
			address = 0 // UNUSED
		case 0x03:
			address = cpu.absoluteAddress()
		}
	} else {
		switch (opcode >> 2) & 0x03 {
		case 0x00:
			address = 0 // UNUSED
		case 0x01:
			switch opcode & 0xf0 {
			case 0x90, 0xb0:
				address = cpu.indexedZeroPageAddress(Y)
			default:
				address = cpu.indexedZeroPageAddress(X)
			}
		case 0x02:
			address = 0 // UNUSED
		case 0x03:
			switch opcode & 0xf0 {
			case 0x90, 0xb0:
				address = cpu.indexedAbsoluteAddress(Y, status)
			default:
				address = cpu.indexedAbsoluteAddress(X, status)
			}
		}
	}
	return
}

// E.1
func (cpu *CPU) zeroPageAddress() (result uint16) {
	result = uint16(cpu.Memory.Read(cpu.Registers.PC))
	cpu.Registers.PC++
	return
}

// E.2
func (cpu *CPU) indexedZeroPageAddress(index Index) (result uint16) {
	value := cpu.Memory.Read(cpu.Registers.PC)
	result = uint16(value + cpu.IndexToRegister(index))
	cpu.Registers.PC++
	return
}

// E.3
func (cpu *CPU) absoluteAddress() (result uint16) {
	low := cpu.Memory.Read(cpu.Registers.PC)
	high := cpu.Memory.Read(cpu.Registers.PC + 1)
	result = (uint16(high) << 8) | uint16(low)
	cpu.Registers.PC += 2
	return
}

// E.4
func (cpu *CPU) indexedAbsoluteAddress(index Index, status *InstructionStatus) (result uint16) {

	low := cpu.Memory.Read(cpu.Registers.PC)
	high := cpu.Memory.Read(cpu.Registers.PC + 1)

	address := (uint16(high) << 8) | uint16(low)
	result = address + uint16(cpu.IndexToRegister(index))

	cpu.Registers.PC += 2

	if status != nil && !SamePage(address, result) {
		*status |= PageCross
	}

	return
}

// E.5
func (cpu *CPU) indirectAddress() (result uint16) {
	low := cpu.Memory.Read(cpu.Registers.PC)
	high := cpu.Memory.Read(cpu.Registers.PC + 1)
	cpu.Registers.PC += 2
	// 6502 had a bug where it incremented only the high byte instead
	// of the whole 16bit address when computing the address.
	high = cpu.Memory.Read((uint16(high) << 8) | uint16(low+1))
	low = cpu.Memory.Read((uint16(high) << 8) | uint16(low))
	result = (uint16(high) << 8) | uint16(low)
	return
}

// E.6 Implied (CLD, NOOP)

// E.7 Accumulator Arithmetic shift left, logical shift right, rotate left,
// rotate right

// E.8
func (cpu *CPU) immediateAddress() (result uint16) {
	result = cpu.Registers.PC
	cpu.Registers.PC++
	return
}

// E.9
func (cpu *CPU) relativeAddress() (result uint16) {
	value := uint16(cpu.Memory.Read(cpu.Registers.PC))
	cpu.Registers.PC++

	var offset uint16
	if value > 0x7f {
		offset = -(0x0100 - value)
	} else {
		offset = value
	}
	result = cpu.Registers.PC + offset
	return
}

// E.10 aka pre-indexed
func (cpu *CPU) indexedIndirectAddress() (result uint16) {
	value := cpu.Memory.Read(cpu.Registers.PC)
	address := uint16(value + cpu.Registers.X)
	low := cpu.Memory.Read(address)
	high := cpu.Memory.Read((address + 1) & 0x00ff)
	result = (uint16(high) << 8) | uint16(low)
	cpu.Registers.PC++
	return
}

// E.11 aka post-indexed
func (cpu *CPU) indirectIndexedAddress(status *InstructionStatus) (result uint16) {
	value := cpu.Memory.Read(cpu.Registers.PC)
	address := uint16(value)
	cpu.Registers.PC++
	low := cpu.Memory.Read(address)
	high := cpu.Memory.Read((address + 1) & 0x00ff)

	address = (uint16(high) << 8) | uint16(low)
	result = address + uint16(cpu.Registers.Y)

	if status != nil && !SamePage(address, result) {
		*status |= PageCross
	}
	return
}

// Helpers
// =======

func (cpu *CPU) push8(value uint8) {
	cpu.Memory.Write(0x0100|uint16(cpu.Registers.SP), value)
	cpu.Registers.SP--
}

func (cpu *CPU) push16(value uint16) {
	cpu.push8(uint8(value >> 8))
	cpu.push8(uint8(value))
}

func (cpu *CPU) setZFlag(value uint8) uint8 {
	if value == 0 {
		cpu.Registers.P |= Z
	} else {
		cpu.Registers.P &= ^Z
	}
	return value
}

func (cpu *CPU) setNFlag(value uint8) uint8 {
	cpu.Registers.P = (cpu.Registers.P & ^N) | Status(value&uint8(N))
	return value
}

func (cpu *CPU) setZNFlags(value uint8) uint8 {
	cpu.setZFlag(value)
	cpu.setNFlag(value)
	return value
}

// CPU Instructions
// ================

// Lda loads a byte of memory into A, setting Z and N flags as required
func (cpu *CPU) Lda(address uint16) {
	cpu.Registers.A = cpu.setZNFlags(cpu.Memory.Read(address))
}
