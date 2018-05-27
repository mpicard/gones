package cpu

import (
	"io"
	"os"
)

const DEFAULT_MEMORY_SIZE uint32 = 65536

type Memory interface {
	Reset()
	Fetch(address uint16) (value uint8)
	Store(address uint16, value uint8) (oldValue uint8)
}

type BasicMemory struct {
	m              []uint8
	disableReads   bool
	disabledWrites bool
}

func NewBasicMemory(size uint32) *BasicMemory {
	return &BasicMemory{
		m: make([]uint8, size),
	}
}

func (mem *BasicMemory) DisableReads() {
	mem.disableReads = true
}

func (mem *BasicMemory) EnableReads() {
	mem.disableReads = false
}

func (mem *BasicMemory) DisabledWrites() {
	mem.disabledWrites = true
}

func (mem *BasicMemory) EnableWrites() {
	mem.disabledWrites = false
}

func (mem *BasicMemory) Reset() {
	for i := range mem.m {
		mem.m[i] = 0x00
	}
}

func (mem *BasicMemory) Fetch(address uint16) (value uint8) {
	if mem.disableReads {
		value = 0xff
	} else {
		value = mem.m[address]
	}
	return
}

func (mem *BasicMemory) Store(address uint16, value uint8) (oldValue uint8) {
	if !mem.disabledWrites {
		oldValue = mem.m[address]
		mem.m[address] = value
	}
	return
}

// SamePage returns true if the two addresses are located on the same
// page in memory. Two addresses are on the same page if their high
//  bytes are both the same 0x0101 and 0x0103 but not 0x0101 and 0x0203
func SamePage(addr1 uint16, addr2 uint16) bool {
	return (addr1^addr2)>>8 == 0
}

func (mem *BasicMemory) load(path string) {
	fi, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	total := 0
	buf := make([]byte, 65536)

	for {
		n, err := fi.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		total++
	}

	j := 0xc000

	for i, b := range buf {
		if i <= 15 {
			continue
		}

		mem.m[j] = b
		j++

		if j == 0xffff {
			break
		}
	}
	return
}
