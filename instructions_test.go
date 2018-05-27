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
