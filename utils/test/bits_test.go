package test

import (
	"fmt"
	"search_engine/utils"
	"testing"
)

func TestBits(t *testing.T) {
	var n uint64
	n = utils.SetBit1(n, 12)
	n = utils.SetBit1(n, 29)

	fmt.Println(utils.IsBit1(n, 12))
	fmt.Println(utils.IsBit1(n, 28))
	fmt.Println(utils.IsBit1(n, 29))

	fmt.Println(utils.CountBit1(n))
	fmt.Printf("%64b\n", n)
}
