package chip8

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

type OpcodeType string

const (
	OpTypeCall OpcodeType = "call"
)

type Op interface {
}

type OpType struct{}
