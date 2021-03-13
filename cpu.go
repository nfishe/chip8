package chip8

import (
	"fmt"
	"log"

	"github.com/nfishe/chip8/internal/atomic"
)

type Register struct {
	x uint8
}

func (r *Register) Equal(s Register) bool {
	return r.x == s.x
}

func (r *Register) Set(x uint8) {
	atomic.Store8(&r.x, x)
}

func (r *Register) Value() uint8 { return r.x }

type CPU struct {
	Display [64][32]uint8

	Memory [4096]uint8

	PC, I uint16

	Stack [16]uint16
	SP    uint8

	// Timers
	Timers struct {
		// Delay timer
		D uint8

		// Sound
		S uint8
	}

	V [16]Register
}

// New initiates creates a new CPU
func New() *CPU {
	return &CPU{
		PC: 0x200,
	}
}

func (cpu *CPU) Cycle() error {
	opcode := cpu.decodeOp()

	log.Println(fmt.Printf("%x", opcode))

	if err := cpu.executeOp(opcode); err != nil {
		return err
	}
	return nil
}

func (cpu *CPU) decodeOp() uint16 {
	return (uint16(cpu.Memory[cpu.PC])<<8 | uint16(cpu.Memory[cpu.PC+1]))
}
