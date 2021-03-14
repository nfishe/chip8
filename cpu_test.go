package chip8

import (
	"testing"

	"github.com/nfishe/chip8/roms"

	"gotest.tools/assert"
)

func TestRegister(t *testing.T) {
	r := Register{0}
	r.Set(1)

	assert.Assert(t, r.Value() == uint8(1))
}

func TestChip8(t *testing.T) {
	fixture := &fixture{}
	cpu := fixture.getCPU()

	assert.Assert(t, cpu.Memory[512] == roms.Blinky[0])
}

func TestCycle(t *testing.T) {
	fixture := &fixture{}
	cpu := fixture.getCPU()

	program := fixture.loadProgram()
	k := len(program) / 2
	for i := 0; i < k; i++ {
		if err := cpu.Cycle(); err != nil {
			t.Error(err)
		}
	}

	assert.Assert(t, true)
}

type fixture struct {
}

func (f *fixture) getCPU() *CPU {
	cpu := New()
	for i, b := range f.loadProgram() {
		cpu.Memory[i+512] = b
	}

	return cpu
}

func (f *fixture) loadProgram() []byte {
	return roms.Blinky
}
