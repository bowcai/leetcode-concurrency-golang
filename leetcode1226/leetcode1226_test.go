package main

import (
	"testing"
)

var test = []int{
	1,
	2,
	3,
	4,
	5,
	10,
	15,
	20,
}

// TestRun only test whether the program has deadlock,for simplicity
func TestRun(t *testing.T) {
	for _, n := range test {
		Run(n)
	}
}

// Test for the second implementation
func TestRunImpl2(t *testing.T) {
	for _, n := range test {
		RunImpl2(n)
	}
}

// Test for the third implementation
func TestRunImpl3(t *testing.T) {
	for _, n := range test {
		RunImpl3(n)
	}
}

// Test for the fourth implementation
func TestRunImpl4(t *testing.T) {
	for _, n := range test {
		RunImpl4(n)
	}
}
