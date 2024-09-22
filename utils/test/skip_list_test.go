package test

import (
	"fmt"
	"search_engine/utils"
	"testing"

	"github.com/huandu/skiplist"
)

func TestSkipList(t *testing.T) {
	list := skiplist.New(skiplist.Int32)
	list.Set(24, 18)
	list.Set(12, 8)
	list.Set(3, 4)
	list.Remove(12)
	if v, ok := list.GetValue(24); ok {
		fmt.Println(v)
	}

	node := list.Front() //找到链表首元素
	for node != nil {
		fmt.Println(node.Key(), node.Value)
		node = node.Next()
	}
}

func TestIntersectionOfSkipList(t *testing.T) {
	l1 := skiplist.New(skiplist.Uint64)
	l1.Set(uint64(5), 0)
	l1.Set(uint64(1), 0)
	l1.Set(uint64(4), 0)
	l1.Set(uint64(9), 0)
	l1.Set(uint64(11), 0)
	l1.Set(uint64(7), 0)

	l2 := skiplist.New(skiplist.Uint64)
	l2.Set(uint64(4), 0)
	l2.Set(uint64(5), 0)
	l2.Set(uint64(9), 0)
	l2.Set(uint64(8), 0)
	l2.Set(uint64(2), 0)

	l3 := skiplist.New(skiplist.Uint64)
	l3.Set(uint64(3), 0)
	l3.Set(uint64(5), 0)
	l3.Set(uint64(7), 0)
	l3.Set(uint64(9), 0)

	insert := utils.IntersectionOfSkipList(l1, l2, l3)
	if insert != nil {
		node := insert.Front()
		for node != nil {
			fmt.Println(node.Key().(uint64))
			node = node.Next()
		}
	}

	fmt.Println("------------------------------------------")
	union := utils.UnionsetOfSkipList(l1, l2, l3)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Println(node.Key().(uint64))
			node = node.Next()
		}
	}
}
