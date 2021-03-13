package chip8

import (
	"testing"

	"github.com/nfishe/chip8/roms"

	"gotest.tools/assert"
)

func TestChip8(t *testing.T) {
	fixture := &fixture{}
	cpu := fixture.getCPU()

	assert.Assert(t, cpu.memory[512] == uint(roms.Blinky[0]))
}

func TestCycle(t *testing.T) {
	fixture := &fixture{}
	cpu := fixture.getCPU()

	if err := cpu.Cycle(); err != nil {
		t.Error(err)
	}
}

type fixture struct {
}

func (f *fixture) getCPU() *CPU {
	cpu := New()
	for i, b := range f.loadProgram() {
		cpu.memory[i+512] = uint(b)
	}

	return cpu
}

func (f *fixture) loadProgram() []byte {
	return roms.Blinky
}
