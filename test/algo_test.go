package test

import (
	"fmt"
	"testing"
)

func TestAlgo(t *testing.T) {
	a := []int{}
	a = append(a, 6)

	b := make([]int, 0)
	b = append(b, 6)

	fmt.Println(a)
	fmt.Println(b)

}
