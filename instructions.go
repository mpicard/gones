package cpu

// OpCode for 6502 CPU
type OpCode uint8

// InstructionStatus
type InstructionStatus uint16

// Instruction implements the instructions for the 6502 CPU
// The Exec implements the instruction and returns the total clock
// cycles to be consumed by the instruction
type Instruction struct {
	Mneumonic string
	OpCode    OpCode
	Exec      func(*CPU) (status InstructionStatus)
}

// InstructionTable maps OpCodes to an instruction
type InstructionTable struct {
	opcodes         []*Instruction
	cycles          []uint16
	cyclesPageCross []uint16
}

const (
	PageCross InstructionStatus = 1 << iota
	Branched
)

// NewInstructionTable returns a new InstructionTable
func NewInstructionTable() InstructionTable {
	instructions := InstructionTable{
		opcodes: make([]*Instruction, 0x100),
		cycles: []uint16{
			7, 6, 0, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
			2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
			2, 6, 0, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
			2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
			2, 5, 0, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
			2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
			2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
			2, 5, 0, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
		},
		cyclesPageCross: []uint16{
			7, 6, 0, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
			6, 6, 0, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
			2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
			3, 6, 0, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
			2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
			3, 6, 0, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
			2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
			2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
			3, 6, 0, 8, 4, 4, 6, 6, 2, 5, 2, 7, 5, 5, 7, 7,
		},
	}
	return instructions
}

// Execute takes instruction from table and executes, returning number of cycles
func (instructions InstructionTable) Execute(cpu *CPU, opcode OpCode) (cycles uint16) {
	inst := instructions.opcodes[opcode]
	if inst == nil {
		return
	}

	status := inst.Exec(cpu)
	if status&PageCross == 0 {
		cycles = instructions.cycles[opcode]
	} else {
		cycles = instructions.cyclesPageCross[opcode]
	}

	if status&Branched != 0 {
		cycles++
	}

	return
}

func (instructions InstructionTable) AddInstruction(inst *Instruction) {
	instructions.opcodes[inst.OpCode] = inst
}

func (instructions InstructionTable) InitInstructions() {
	// http://www.thealmightyguru.com/Games/Hacking/Wiki/index.php?title=6502_Opcodes

	// Storage
	// =======

	// LDA
	for _, o := range []OpCode{0xa1, 0xa5, 0xa9, 0xad, 0xb1, 0xb5, 0xb9, 0xbd} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "LDA",
			OpCode:    opcode,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Lda(cpu.aluAddress(opcode, &status))
				return
			}})
	}

	// LDX
	for _, o := range []OpCode{0xa2, 0xa6, 0xb6, 0xae, 0xbe} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "LDX",
			OpCode:    o,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Ldx(cpu.rmwAddress(opcode, &status))
				return
			}})
	}

	// LDY
	for _, o := range []OpCode{0xa0, 0xa4, 0xb4, 0xac, 0xbc} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "LDY",
			OpCode:    o,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Ldy(cpu.controlAddress(opcode, &status))
				return
			}})
	}

	// STA
	for _, o := range []OpCode{0x85, 0x95, 0x8d, 0x9d, 0x99, 0x81, 0x91} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "STA",
			OpCode:    opcode,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Sta(cpu.aluAddress(opcode, &status))
				return
			}})
	}

	for _, o := range []OpCode{0x86, 0x8e, 0x96} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "STX",
			OpCode:    opcode,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Stx(cpu.rmwAddress(opcode, &status))
				return
			}})
	}

	for _, o := range []OpCode{0x84, 0x94, 0x8c} {
		opcode := o
		instructions.AddInstruction(&Instruction{
			Mneumonic: "STY",
			OpCode:    opcode,
			Exec: func(cpu *CPU) (status InstructionStatus) {
				cpu.Sty(cpu.controlAddress(opcode, &status))
				return
			}})
	}
}
