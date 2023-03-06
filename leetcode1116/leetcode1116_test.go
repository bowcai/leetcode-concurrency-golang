package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

// Sequential version of zeroEvenOdd
func zeroEvenOddSeq(n int) string {
	sb := strings.Builder{}

	for i := 1; i <= n; i++ {
		sb.WriteByte('0')
		sb.WriteString(strconv.Itoa(i))
	}

	return sb.String()
}

func TestRun(t *testing.T) {
	var test = []int{
		1,
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
		want := zeroEvenOddSeq(input)

		if got != want {
			t.Errorf("Output of %d: %q. Expecting: %q", input, got, want)
		}
	}
}
