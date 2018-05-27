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

// Prints the values of each register to os.Stderr
func (reg *Registers) String() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		reg.A, reg.X, reg.Y, reg.P, reg.SP)
}

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
	debug        decode // for debugging
}

type decode struct {
	enabled     bool
	pc          uint16
	opcode      OpCode
	args        string
	mneumonic   string
	decodedArgs string
	registers   string
	ticks       uint64
}

func (d *decode) String() string {
	return fmt.Sprintf("%04X  %02X %-5s %4s %-26s  %25s",
		d.pc, d.opcode, d.args, d.mneumonic, d.decodedArgs, d.registers)
}

func NewCPU(mem Memory) *CPU {
	instructions := NewInstructionTable()
	instructions.InitInstructions()

	return &CPU{
		debug:        decode{},
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
	// fetch
	opcode := OpCode(cpu.Memory.Read(cpu.Registers.PC))
	inst := cpu.Instructions.opcodes[opcode]
	if inst == nil {
		return 0, BadOpCodeError(opcode)
	}

	// execute
	if cpu.debug.enabled {
		cpu.debug.pc = cpu.Registers.PC
		cpu.debug.opcode = opcode
		cpu.debug.args = ""
		cpu.debug.mneumonic = inst.Mneumonic
		cpu.debug.decodedArgs = ""
		cpu.debug.registers = cpu.Registers.String()
	}
	cpu.Registers.PC++
	cycles += cpu.Instructions.Execute(cpu, opcode)
	if cpu.debug.enabled {
		fmt.Println(cpu.debug.String())
		return cycles, BrkOpCodeError(opcode)
	}

	return cycles, nil
}
