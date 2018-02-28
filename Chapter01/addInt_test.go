package main

import (
	"testing"
)

func TestAddInt(t *testing.T) {
	testCases := []struct {
		Name     string
		Values   []int
		Expected int
	}{
		{"addInt() -> 0", []int{}, 0},
		{"addInt([]int{10, 20, 100}) -> 130", []int{10, 20, 100}, 130},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sum := addInt(tc.Values...)
			if sum != tc.Expected {
				t.Errorf("%d != %d", sum, tc.Expected)
			} else {
				t.Logf("%d == %d", sum, tc.Expected)
			}
		})
	}
}
