package test

import (
	"fmt"
	"search_engine/types"
	"testing"
)

func TestS(t *testing.T) {
	a := types.NewTermQuery("", "A")
	b := types.NewTermQuery("", "B")
	c := types.NewTermQuery("", "C")

	s := a.Or(b).And(c)
	fmt.Println(s.ToString())
}
