package utils

func IsBit1(n uint64, i int) bool {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	if n&c == c {
		return true
	} else {
		return false
	}
}

func SetBit1(n uint64, i int) uint64 {
	if i > 64 {
		panic(i)
	}
	c := uint64(1 << (i - 1))
	return n | c
}

func CountBit1(n uint64) int {
	c := uint64(1)
	sum := 0
	for i := 0; i < 64; i++ {
		if c&n == c {
			sum++
		}
		c = c << 1
	}
	return sum
}
