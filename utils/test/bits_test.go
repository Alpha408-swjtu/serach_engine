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

var arr1 = []int{15, 18, 29, 36, 47, 60}
var arr2 = []int{18, 28, 36, 43, 47, 61}

func TestIntersection(t *testing.T) {
	min := 10

	bm1 := utils.CreateBitMap(min, arr1)
	bm2 := utils.CreateBitMap(min, arr2)
	arr := utils.IntersectionOfBitMap(bm1, bm2, min)
	for _, v := range arr {
		fmt.Println(v)
	}

	c := utils.IntersectionOfOrderedList(arr1, arr2)
	fmt.Println(c)
}
