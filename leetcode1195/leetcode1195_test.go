package main

import (
	"bytes"
	"testing"
)

// The sequential version of fizzbuzz
func fizzBuzzSeq(n int) {
	for i := 1; i <= n; i++ {
		if i%3 == 0 {
			if i%5 == 0 { // i is divisible by 3 and 5.
				printFizzBuzz()
			} else { // i is divisible by 3 but not by 5.
				printFizz()
			}
		} else { // i is not divisible by 3.
			if i%5 == 0 { // i is divisible by 5 but not by 3.
				printBuzz()
			} else { // i is not divisible by 3 or 5.
				printNumber(i)
			}
		}
	}
}

// testCase represents a test case to test the fizzbuzz program
type testCase struct {
	input int
	want  string
}

// newTestCase generates a new test case
func newTestCase(input int) *testCase {
	output := testCase{input: input}

	// Modify the output from stdout to a new buffer
	out = new(bytes.Buffer)
	// Get the result from the sequential version of fizzbuzz
	fizzBuzzSeq(input)
	output.want = out.(*bytes.Buffer).String()

	return &output
}

func TestRun(t *testing.T) {
	testCases := []*testCase{
		newTestCase(5),
		newTestCase(10),
		newTestCase(15),
		newTestCase(20),
		newTestCase(32),
		newTestCase(64),
	}

	for _, c := range testCases {
		// Modify the output from stdout to a new buffer
		out = new(bytes.Buffer)

		Run(c.input)
		got := out.(*bytes.Buffer).String()

		if got != c.want {
			t.Errorf("Output of %d: %q. Expecting: %q", c.input, got, c.want)
		}
	}
}
