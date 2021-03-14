package atomic

func And8(ptr *uint8, val uint8) {
	*ptr = *ptr & val
}

func Store8(ptr *uint8, val uint8) {
	*ptr = val
}

func Store(ptr *uint16, val uint16) {
	*ptr = val
}

func Add(ptr *uint8, delta uint8) {
	*ptr = *ptr + delta
}
