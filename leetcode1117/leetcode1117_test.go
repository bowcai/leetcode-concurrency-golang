package main

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

// generateInput generates an input contains n H2O molecules randomly.
func generateInput(n int, rng *rand.Rand) string {
	inputLen := 3 * n
	inputBytes := make([]byte, inputLen)

	for i := 0; i < n; i++ {
		inputBytes[i] = 'O'
	}

	for i := n; i < inputLen; i++ {
		inputBytes[i] = 'H'
	}

	// Use Random number generator to shuffle the input
	rng.Shuffle(inputLen, func(i, j int) {
		inputBytes[i], inputBytes[j] = inputBytes[j], inputBytes[i]
	})

	return string(inputBytes)
}

// isValidResult tests whether a result is valid
func isValidResult(result string, n int) bool {
	if len(result) != 3*n {
		return false
	}

	hydrogenCnt := 0
	oxygenCnt := 0
	for i, c := range result {
		switch c {
		case 'H':
			hydrogenCnt++
		case 'O':
			oxygenCnt++
		default:
			return false
		}

		// Check if the result is valid every 3 characters
		if i%3 == 2 {
			if hydrogenCnt != 2 || oxygenCnt != 1 {
				return false
			}

			hydrogenCnt = 0
			oxygenCnt = 0
		}
	}

	// Passed through every test, is a valid result
	return true
}

func TestRun(t *testing.T) {
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, n := range test {
		// Modify the output from stdout to a new buffer
		out = new(bytes.Buffer)

		input := generateInput(n, rng)
		run(input)

		got := out.(*bytes.Buffer).String()
		valid := isValidResult(got, n)

		if !valid {
			t.Errorf("Output of %q: %q. Incorrect output.", input, got)
		}
	}
}
