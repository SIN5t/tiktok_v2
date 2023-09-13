package test

import (
	"fmt"
	"testing"
)

func TestAlgo(t *testing.T) {
	res := isAnagram("cx", "lt")
	fmt.Println(res)
}
func isAnagram(s string, t string) bool {
	table := make([]int, 26)
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		index1 := s[i] - 'a'
		index2 := t[i] - 'a'

		table[index1]++
		table[index2]--
	}
	for i := 0; i < len(s); i++ {
		if table[i] != 0 {
			return false
		}
	}
	return true
}
