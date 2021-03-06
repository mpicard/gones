package cpu

import (
	"testing"
	"time"
)

const rate time.Duration = 46 * time.Nanosecond // 21.477272MHz
const divisor = 12

var cpu *CPU

func Setup() {
	cpu = NewCPU(NewBasicMemory(DEFAULT_MEMORY_SIZE))
	cpu.Reset()
	cpu.breakError = true
}

func Teardown() {

}

func TestBadOpcodeError(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x02)

	_, err := cpu.Execute()

	if err == nil {
		t.Error("No error returned")
	}

	if _, ok := err.(BadOpCodeError); !ok {
		t.Error("Did not receive expected error type BadOpCodeError")
	}

	Teardown()
}

func TestLdaImmediate(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa9)
	cpu.Memory.Write(0x0101, 0xff)

	cpu.Execute()

	if cpu.Registers.A != 0xff {
		t.Errorf("Register A 0xff != %#x", cpu.Registers.A)
	}

	Teardown()
}

func TestLdaZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa5)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.A != 0xff {
		t.Error("Register A is not 0xff")
	}

	Teardown()
}

func TestLdaZeroPageX(t *testing.T) {
	Setup()

	cpu.Registers.X = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xb5)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0085, 0xff)

	cpu.Execute()

	if cpu.Registers.A != 0xff {
		t.Error("Register A is not 0xff")
	}

	Teardown()
}

func TestLdaAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xad)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.A != 0xff {
		t.Error("Register A is not 0xff")
	}

	Teardown()
}

func TestLdaAbsoluteX(t *testing.T) {
	Setup()

	cpu.Registers.X = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xbd)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0085, 0xff)

	cycles, _ := cpu.Execute()

	if cycles != 4 {
		t.Error("Cycles is not 4")
	}

	if cpu.Registers.A != 0xff {
		t.Error("Register A is not 0xff")
	}

	cpu.Registers.X = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xbd)
	cpu.Memory.Write(0x0101, 0xff)
	cpu.Memory.Write(0x0102, 0x02)
	cpu.Memory.Write(0x0300, 0xff)

	cycles, _ = cpu.Execute()

	if cycles != 5 {
		t.Errorf("Cycles is %v not 5", cycles)
	}

	Teardown()
}

func TestLdaAbsoluteY(t *testing.T) {
	Setup()

	cpu.Registers.Y = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xb9)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0085, 0xff)

	cycles, _ := cpu.Execute()

	if cycles != 4 {
		t.Error("Cycles is not 4")
	}

	if cpu.Registers.A != 0xff {
		t.Error("Register A is not 0xff")
	}

	cpu.Registers.Y = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xb9)
	cpu.Memory.Write(0x0101, 0xff)
	cpu.Memory.Write(0x0102, 0x02)
	cpu.Memory.Write(0x0300, 0xff)

	cycles, _ = cpu.Execute()

	if cycles != 5 {
		t.Error("Cycles is not 5")
	}

	Teardown()
}

func TestLdxImmediate(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa2)
	cpu.Memory.Write(0x0101, 0xff)

	cpu.Execute()

	if cpu.Registers.X != 0xff {
		t.Error("Register X is not 0xff")
	}

	Teardown()
}

func TestLdxZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa6)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.X != 0xff {
		t.Error("Register X is not 0xff")
	}

	Teardown()
}

func TestLdxZeroPageY(t *testing.T) {
	Setup()

	cpu.Registers.Y = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xb6)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0085, 0xff)

	cpu.Execute()

	if cpu.Registers.X != 0xff {
		t.Error("Register X is not 0xff")
	}

	Teardown()
}

func TestLdxAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xae)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.X != 0xff {
		t.Error("Register X is not 0xff")
	}

	Teardown()
}

func TestLdxAbsoluteY(t *testing.T) {
	Setup()

	cpu.Registers.Y = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xbe)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0085, 0xff)

	cpu.Execute()

	if cpu.Registers.X != 0xff {
		t.Error("Register X is not 0xff")
	}

	Teardown()
}

func TestLdxZFlagSet(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa2)
	cpu.Memory.Write(0x0101, 0x00)

	cpu.Execute()

	if cpu.Registers.P&Z == 0 {
		t.Error("Z flag is not set")
	}

	Teardown()
}

func TestLdxZFlagUnset(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa2)
	cpu.Memory.Write(0x0101, 0x01)

	cpu.Execute()

	if cpu.Registers.P&Z != 0 {
		t.Error("Z flag is set")
	}

	Teardown()
}

func TestLdxNFlagSet(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa2)
	cpu.Memory.Write(0x0101, 0x81)

	cpu.Execute()

	if cpu.Registers.P&N == 0 {
		t.Error("N flag is not set")
	}

	Teardown()
}

func TestLdxNFlagUnset(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa2)
	cpu.Memory.Write(0x0101, 0x01)

	cpu.Execute()

	if cpu.Registers.P&N != 0 {
		t.Error("N flag is set")
	}

	Teardown()
}

func TestLdyImmediate(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa0)
	cpu.Memory.Write(0x0101, 0xff)

	cpu.Execute()

	if cpu.Registers.Y != 0xff {
		t.Error("Register Y is not 0xff")
	}

	Teardown()
}

func TestLdyZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa4)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.Y != 0xff {
		t.Error("Register Y is not 0xff")
	}

	Teardown()
}

func TestLdyZeroPageX(t *testing.T) {
	Setup()

	cpu.Registers.X = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xb4)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0085, 0xff)

	cpu.Execute()

	if cpu.Registers.Y != 0xff {
		t.Error("Register Y is not 0xff")
	}

	Teardown()
}

func TestLdyAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xac)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0084, 0xff)

	cpu.Execute()

	if cpu.Registers.Y != 0xff {
		t.Error("Register Y is not 0xff")
	}

	Teardown()
}

func TestLdyAbsoluteX(t *testing.T) {
	Setup()

	cpu.Registers.X = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xbc)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0085, 0xff)

	cpu.Execute()

	if cpu.Registers.Y != 0xff {
		t.Error("Register Y is not 0xff")
	}

	Teardown()
}

func TestLdyZFlagSet(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa0)
	cpu.Memory.Write(0x0101, 0x00)

	cpu.Execute()

	if cpu.Registers.P&Z == 0 {
		t.Error("Z flag is not set")
	}

	Teardown()
}

func TestLdyZFlagUnset(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa0)
	cpu.Memory.Write(0x0101, 0x01)

	cpu.Execute()

	if cpu.Registers.P&Z != 0 {
		t.Error("Z flag is set")
	}

	Teardown()
}

func TestLdyNFlagSet(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa0)
	cpu.Memory.Write(0x0101, 0x81)

	cpu.Execute()

	if cpu.Registers.P&N == 0 {
		t.Error("N flag is not set")
	}

	Teardown()
}

func TestLdyNFlagUnset(t *testing.T) {
	Setup()

	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0xa0)
	cpu.Memory.Write(0x0101, 0x01)

	cpu.Execute()

	if cpu.Registers.P&N != 0 {
		t.Error("N flag is set")
	}

	Teardown()
}

func TestStaZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x85)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaZeroPageX(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.X = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x95)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0085) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x8d)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaAbsoluteX(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.X = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x9d)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0085) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaAbsoluteY(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.Y = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x99)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0085) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaIndirectX(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.X = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x81)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0085, 0x87)
	cpu.Memory.Write(0x0086, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0087) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStaIndirectY(t *testing.T) {
	Setup()

	cpu.Registers.A = 0xff
	cpu.Registers.Y = 1
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x91)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0084, 0x86)
	cpu.Memory.Write(0x0085, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0087) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStxZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.X = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x86)
	cpu.Memory.Write(0x0100, 0x86)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStxZeroPageY(t *testing.T) {
	Setup()

	cpu.Registers.X = 0xff
	cpu.Registers.Y = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x96)
	cpu.Memory.Write(0x0100, 0x96)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0085) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStxAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.X = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x8e)
	cpu.Memory.Write(0x0100, 0x8e)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0102, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStyZeroPage(t *testing.T) {
	Setup()

	cpu.Registers.Y = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x84)
	cpu.Memory.Write(0x0100, 0x84)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStyZeroPageY(t *testing.T) {
	Setup()

	cpu.Registers.Y = 0xff
	cpu.Registers.X = 0x01
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x94)
	cpu.Memory.Write(0x0100, 0x94)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)

	cpu.Execute()

	if cpu.Memory.Read(0x0085) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}

func TestStyAbsolute(t *testing.T) {
	Setup()

	cpu.Registers.Y = 0xff
	cpu.Registers.PC = 0x0100

	cpu.Memory.Write(0x0100, 0x8c)
	cpu.Memory.Write(0x0100, 0x8c)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0101, 0x84)
	cpu.Memory.Write(0x0102, 0x00)
	cpu.Memory.Write(0x0102, 0x00)

	cpu.Execute()

	if cpu.Memory.Read(0x0084) != 0xff {
		t.Error("Memory is not 0xff")
	}

	Teardown()
}
