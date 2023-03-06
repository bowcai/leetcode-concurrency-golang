package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	var test = []int{
		5,
		10,
		15,
		20,
		30,
		64,
		100,
		200,
		1000,
	}

	for _, input := range test {
		// Modify the output from stdout to a new buffer
		out = new(bytes.Buffer)

		run(input)

		got := out.(*bytes.Buffer).String()
		want := strings.Repeat("foobar", input)

		if got != want {
			t.Errorf("Output of %d: %q. Expecting: %q", input, got, want)
		}
	}
}
