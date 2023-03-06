package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	const want = "firstsecondthird"

	var test = [][3]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	for _, input := range test {
		// Modify the output from stdout to a new buffer
		out = new(bytes.Buffer)

		// Run the code with different order
		run(input)

		got := out.(*bytes.Buffer).String()
		if got != want {
			t.Errorf("Output of %v: %q. Expecting: %q", input, got, want)
		}
	}
}
