package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Foo struct {
	wg         *sync.WaitGroup // Wait for all the go routine to finish
	firstDone  chan struct{}   // Notify the function "first" is done
	secondDone chan struct{}   // Notify the function "second" is done
}

// Specify the output of printing
// Will be modified during testing
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem
// For simplicity, error is not handled here
func (f *Foo) printFirst() {
	fmt.Fprint(out, "first")
}
func (f *Foo) printSecond() {
	fmt.Fprint(out, "second")
}
func (f *Foo) printThird() {
	fmt.Fprint(out, "third")
}

func (f *Foo) first() {
	defer f.wg.Done()

	f.printFirst()

	// Notify that the function "first" is done.
	f.firstDone <- struct{}{}
}

func (f *Foo) second() {
	defer f.wg.Done()

	// Wait for the function "first" to finish.
	<-f.firstDone

	f.printSecond()

	// Notify that the function "second" is done.
	f.secondDone <- struct{}{}
}

func (f *Foo) third() {
	defer f.wg.Done()

	// Wait for the function "second" to finish.
	<-f.secondDone

	f.printThird()
}

// run starts the three goroutines, in a specified order.
// The order must be a permutation of [1,2,3].
func run(order [3]int) {
	f := Foo{
		wg:         new(sync.WaitGroup),
		firstDone:  make(chan struct{}),
		secondDone: make(chan struct{}),
	}

	for _, idx := range order {
		f.wg.Add(1)

		// Start the goroutines with specified order.
		switch idx {
		case 1:
			go f.first()
		case 2:
			go f.second()
		case 3:
			go f.third()
		}

		// Add a short delay after the start of each goroutine
		// to simulate the start order.
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all goroutines to finish.
	f.wg.Wait()
}

func main() {
	run([3]int{1, 2, 3})
}
