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
	cpu.debug.enabled = true
}

func Teardown() {

}

func TestBedOpcodeError(t *testing.T) {
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
