package test

import (
	"fmt"
	"search_engine/types"
	"strings"
	"testing"
)

//输出(A|B|C)&D|E&(F|G)&H

func should(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString("(")
	for _, ele := range s {
		if ele != "" {
			sb.WriteString(ele + "|")
		}
	}
	result := sb.String()
	return result[0:len(result)-1] + ")"
}

func must(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString("(")
	for _, ele := range s {
		if ele != "" {
			sb.WriteString(ele + "&")
		}
	}
	result := sb.String()
	return result[0:len(result)-1] + ")"
}

func TestA(t *testing.T) {
	S := must(must(should(must(should("A", "B", "C"), "D"), "E"), should("F", "G")), "H")
	fmt.Println(S)
}
func TestB(t *testing.T) {
	T := types.TermQueryV0{Keyword: "A"}
	fmt.Println(T.String())
}
