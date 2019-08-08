package main

import (
	"fmt"
	"testing"
)

func TestIsLowMemory(t *testing.T) {
	tests := []struct {
		capacity          string
		capacityAvailable string
		expected          bool
		desc              string
	}{
		{"4.9931317248e+10", "5.1054829568e+10", false, "capacity/capacity_available > 0.15"},
		{"0.9931317248e+10", "6.9054829568e+10", true, "capacity/capacity_available < 0.15"},
		{"0", "6.1054829568e+10", true, "capacity = 0"},
		{"0.9931317248e+10", "0", true, "capacity_available = 0"},
	}

	for _, test := range tests {
		content := fmt.Sprintf("capacity{store=\"1\"} %s\ncapacity_available{store=\"1\"} %s", test.capacity, test.capacityAvailable)
		got := isLowMemory(content)
		if got != test.expected {
			t.Errorf("failed testing %s", test.desc)
		}

	}
}
