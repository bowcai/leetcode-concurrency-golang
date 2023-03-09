package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type foo struct {
	firstDone  chan struct{} // Notify the function "first" is done
	secondDone chan struct{} // Notify the function "second" is done
}

// Specify the output of printing
// Will be modified during testing
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem
func printFirst() {
	if _, err := fmt.Fprint(out, "first"); err != nil {
		panic(err)
	}
}
func printSecond() {
	if _, err := fmt.Fprint(out, "second"); err != nil {
		panic(err)
	}
}
func printThird() {
	if _, err := fmt.Fprint(out, "third"); err != nil {
		panic(err)
	}
}

func (f *foo) first(printFirst func()) {
	printFirst()

	// Notify that the function "first" is done.
	f.firstDone <- struct{}{}
}

func (f *foo) second(printSecond func()) {
	// Wait for the function "first" to finish.
	<-f.firstDone

	printSecond()

	// Notify that the function "second" is done.
	f.secondDone <- struct{}{}
}

func (f *foo) third(printThird func()) {
	// Wait for the function "second" to finish.
	<-f.secondDone

	printThird()
}

// Run starts the three goroutines, in a specified order.
// The order must be a permutation of [1,2,3].
func Run(order [3]int) {
	f := foo{
		firstDone:  make(chan struct{}),
		secondDone: make(chan struct{}),
	}

	var wg sync.WaitGroup

	for _, idx := range order {
		wg.Add(1)

		// Start the goroutines with specified order.
		switch idx {
		case 1:
			go func() {
				defer wg.Done()
				f.first(printFirst)
			}()
		case 2:
			go func() {
				defer wg.Done()
				f.second(printSecond)
			}()
		case 3:
			go func() {
				defer wg.Done()
				f.third(printThird)
			}()
		}

		// Add a short delay after the start of each goroutine
		// to simulate the start order.
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}

func main() {
	Run([3]int{1, 2, 3})
}
