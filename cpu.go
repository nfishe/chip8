package chip8

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

// Memory
type Memory [4096]uint

type CPU struct {
	memory Memory

	mu sync.Mutex
	pc uint64
}

// New initiates creates a new CPU
func New() *CPU {
	var m Memory
	return &CPU{
		memory: m,
		pc:     0x200,
	}
}

func (cpu *CPU) Cycle() error {
	opcode := (cpu.memory[cpu.pc]<<8 | cpu.memory[cpu.pc+1])
	err := func() error {
		cpu.executeOp(opcode)
		return nil
	}()

	return err
}

// executeOp
func (cpu *CPU) executeOp(opcode uint) {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	atomic.AddUint64(&cpu.pc, 2)

	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	log.Println(fmt.Printf("%v %v", x, y))
}
