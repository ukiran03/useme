package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestParseOnlyInts(t *testing.T) {
	testcases := []struct {
		input    string
		expected []int
	}{
		{"1,2 ,3  ,4, 5", []int{1, 2, 3, 4, 5}},
		{"1,2,3", []int{1, 2, 3}},
		{"1", []int{1}},
	}

	for i, tc := range testcases {
		tname := fmt.Sprintf("test %d: ", i)
		t.Run(tname, func(t *testing.T) {
			out := parseOnlyInts(tc.input)
			if !slices.Equal(out, tc.expected) {
				t.Errorf("got %v, expected %v", out, tc.expected)
			}
		})
	}
}
