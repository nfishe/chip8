package chip8

import (
	"github.com/nfishe/chip8/internal/atomic"
	"github.com/nfishe/chip8/util/rand"
)

type OpRegistryBuilder []func(*OpRegistry) error

func NewOp(ops ...func(*OpRegistry) error) OpRegistryBuilder {
	var or OpRegistryBuilder
	for _, op := range ops {
		or.AddOp(op)
	}
	return or
}

func (orb *OpRegistryBuilder) AddOp(f func(*OpRegistry) error) {
	*orb = append(*orb, f)
}

type OpRegistry struct{}

func (r *OpRegistry) getOp() {}

type OpcodeType string

const (
	OpTypeCall OpcodeType = "call"
)

type Op interface {
}

type OpType struct{}

// executeOp
func (cpu *CPU) executeOp(opcode uint16) error {
	switch opcode & 0xf000 {
	case 0x0000:
		switch opcode {
		case 0x00e0:
			atomic.Store(&cpu.PC, 2)
		case 0x00ee:
			atomic.Store(&cpu.PC, cpu.Stack[cpu.SP])
			atomic.Add(&cpu.SP, ^uint8(0))

			atomic.Store(&cpu.PC, 2)
		default:
			//panic("")
		}

	// JP
	case 0x1000:
		atomic.Store(&cpu.PC, opcode&0x0fff)
	// CALL
	case 0x2000:
		// Increment the stack pointer
		atomic.Add(&cpu.SP, 1)

		cpu.Stack[cpu.SP] = cpu.PC
		cpu.PC = opcode & 0xffff
	// SE Vx, byte
	case 0x3000:
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode)

		atomic.Store(&cpu.PC, 2)

		vx := cpu.V[x]
		if vx.Equal(Register{kk}) {
			atomic.Store(&cpu.PC, 2)
		}
	case 0x4000:
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode)

		atomic.Store(&cpu.PC, 2)

		vx := cpu.V[x]
		if !vx.Equal(Register{kk}) {
			atomic.Store(&cpu.PC, 2)
		}
	case 0x5000:
		switch opcode & 0xf00f {
		case 0x5000:
			vx, vy := cpu.V[(opcode&0x0F00)>>8], cpu.V[(opcode&0x00F0)>>4]

			atomic.Store(&cpu.PC, 2)

			if vx.Equal(vy) {
				atomic.Store(&cpu.PC, 2)
			}
		}
	// 6xkk - LD Vx, byte
	case 0x6000:
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode)

		vx := cpu.V[x]
		vx.Set(kk)

		atomic.Store(&cpu.PC, 2)

		// 7xkk - ADD Vx, byte
	case 0x7000:
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode)

		vx := cpu.V[x]
		vx.Set(vx.Value() + kk)

		atomic.Store(&cpu.PC, 2)
	case 0x8000:
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4

		switch opcode & 0x000f {
		case 0x0000:
			vx, vy := cpu.V[x], cpu.V[y]
			vx.Set(vy.Value())

			atomic.Store(&cpu.PC, 2)
		case 0x0001:
			vx, vy := cpu.V[x], cpu.V[y]
			vx.Set(vy.Value() | vx.Value())

			atomic.Store(&cpu.PC, 2)
		case 0x0002:
			vx, vy := cpu.V[x], cpu.V[y]
			vx.Set(vy.Value() & vx.Value())

			atomic.Store(&cpu.PC, 2)
		case 0x0003:
			vx, vy := cpu.V[x], cpu.V[y]
			vx.Set(vy.Value() ^ vx.Value())

			atomic.Store(&cpu.PC, 2)
		case 0x0004:
			vx, vy := cpu.V[x], cpu.V[y]
			r := uint16(vx.Value()) + uint16(vy.Value())

			var cf uint8
			if r > 0xFF {
				cf = 1
			}
			v := cpu.V[0xF]
			v.Set(cf)

			b := uint8(r)
			vx = cpu.V[x]
			vx.Set(b)

			atomic.Store(&cpu.PC, 2)
		case 0x0005:
			vx, vy := cpu.V[x], cpu.V[y]
			var cf uint8

			if vx.Value() > vy.Value() {
				cf = 1
			}
			r := cpu.V[0xf]
			r.Set(cf)

			vx.Set(vx.Value() - vy.Value())

			atomic.Store(&cpu.PC, 2)
		case 0x0006:
			var cf uint8
			vx := cpu.V[x]
			if (vx.Value() & 0x01) == 0x01 {
				cf = 1
			}

			v := cpu.V[0xf]
			v.Set(cf)

			vx.Set(vx.Value() / 2)

			atomic.Store(&cpu.PC, 2)
		// 8xy7 - SUBN Vx, Vy
		case 0x0007:
			vx, vy := cpu.V[x], cpu.V[y]
			var cf uint8
			if vy.Value() > vx.Value() {
				cf = 1
			}

			v := cpu.V[0xf]
			v.Set(cf)

			vx.Set(vy.Value() - vx.Value())

			atomic.Store(&cpu.PC, 2)
		// 8xyE - SHL Vx {, Vy}
		case 0x000E:
			vx := cpu.V[x]
			var cf uint8
			if (vx.Value() & 0x80) == 0x80 {
				cf = 1
			}

			v := cpu.V[0xf]
			v.Set(cf)

			vx.Set(vx.Value() * 2)

			atomic.Store(&cpu.PC, 2)
		// 0x9XY0
		case 0x9000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			switch opcode & 0x000F {
			// 9xy0 - SNE Vx, Vy
			case 0x0000:
				atomic.Store(&cpu.PC, 2)

				vx, vy := cpu.V[x], cpu.V[y]
				if vx != vy {
					atomic.Store(&cpu.PC, 2)
				}
			default:
				//
			}
		case 0xA000: // Annn - LD I, addr
			atomic.Store(&cpu.I, opcode&0x0fff)
			atomic.Store(&cpu.PC, 2)
		case 0xB000: // Bnnn - JP V0, addr
			v0 := cpu.V[0]
			atomic.Store(&cpu.PC, opcode&0x0fff+uint16(v0.Value()))
		case 0xC000: // Cxkk - RND Vx, byte
			x := (opcode & 0x0F00) >> 8
			kk := uint8(opcode)

			vx := cpu.V[x]
			vx.Set(kk + uint8(rand.Intn(255)))
		case 0xD000: // Dxyn - DRW Vx, Vy, nibble
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			vx, vy := cpu.V[x], cpu.V[y]
			n := opcode & 0x000F

			var collision bool
			sprite := cpu.Memory[cpu.I : cpu.I+n]
			for i := 0; i < int(n); i++ {
				r := sprite[i]

				for j := 0; j < 8; j++ {
					ii := 0x80 >> uint8(i)
					on := (r & uint8(ii)) == uint8(ii)

					ip := uint16(vx.Value()) + uint16(i)
					jp := uint16(vy.Value()) + uint16(j)

					collision = func(i, j uint16, on bool) (ok bool) {
						if cpu.Display[ip][jp] == 0x01 {
							ok = true
						}

						var v uint8
						if on {
							v = 0x01
						}

						cpu.Display[i][j] = cpu.Display[i][j] ^ v
						return
					}(ip, jp, on)
				}
			}

			var cf uint8
			if collision {
				cf = 0x01
			}

			v := cpu.V[0xf]
			v.Set(cf)

			atomic.Store(&cpu.PC, 2)

			// draw
		case 0xE000:
			x := (opcode & 0x0F00) >> 8
			switch opcode & 0x00FF {
			case 0x9E: // Ex9E - SKP Vx
				atomic.Store(&cpu.PC, 2)

				vx := cpu.V[x]
				if vx.Value() < uint8(0) {
					//
				}
			case 0xA1: // ExA1 - SKNP Vx
				atomic.Store(&cpu.PC, 2)

				vx := cpu.V[x]
				if vx.Value() < uint8(0) {
					//
				}
			default:
				//
			}
		case 0xF000:
			x := (opcode & 0x0F00) >> 8
			switch opcode & 0x00FF {
			case 0x07: // Fx07 - LD Vx, DT
				vx := cpu.V[x]
				vx.Set(cpu.Timers.D)

				atomic.Store(&cpu.PC, 2)
			case 0x0A: // Fx0A - LD Vx, K
				atomic.Store(&cpu.PC, 2)
			case 0x15: // Fx15 - LD DT, Vx
				vx := cpu.V[x]

				atomic.Store8(&cpu.Timers.D, vx.Value())
				atomic.Store(&cpu.PC, 2)
			case 0x18: // Fx18 - LD ST, Vx
				vx := cpu.V[x]

				atomic.Store8(&cpu.Timers.S, vx.Value())
				atomic.Store(&cpu.PC, 2)
			case 0x1E: // Fx1E - ADD I, Vx
				vx := cpu.V[x]

				atomic.Store(&cpu.I, cpu.I+uint16(vx.Value()))
				atomic.Store(&cpu.PC, 2)
			case 0x29: // Fx29 - LD F, Vx
				vx := cpu.V[x]

				atomic.Store(&cpu.I, uint16(vx.Value())*uint16(0x05))
				atomic.Store(&cpu.PC, 2)
			case 0x33: // Fx33 - LD B, Vx
				vx := cpu.V[x]

				atomic.Store8(&cpu.Memory[cpu.I], vx.Value()/100)
				atomic.Store8(&cpu.Memory[cpu.I+1], (vx.Value()/10)%10)
				atomic.Store8(&cpu.Memory[cpu.I+2], (vx.Value()%100)%10)
				atomic.Store(&cpu.PC, 2)
			case 0x55: // Fx55 - LD [I], Vx
				for i := 0; uint16(i) <= x; i++ {
					vi := cpu.V[i]
					atomic.Store8(&cpu.Memory[cpu.I+uint16(i)], vi.Value())
				}

				atomic.Store(&cpu.PC, 2)
			case 0x65: // Fx65 - LD Vx, [I]
				for i := 0; uint8(i) <= uint8(x); i++ {
					vi := cpu.V[i]
					vi.Set(cpu.Memory[cpu.I+uint16(i)])
				}

				atomic.Store(&cpu.PC, 2)
			default:
			}
		}
	default:
	}
	return nil
}
