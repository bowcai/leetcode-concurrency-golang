package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type FooBar struct {
	n       int
	wg      *sync.WaitGroup // Wait for all the go routine to finish
	fooDone chan struct{}   // Notify the function foo is done
	barDone chan struct{}   // Notify the function bar is done
}

// Specify the output of printing
// Will be modified during testing
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem
// For simplicity, error is not handled here
func printFoo() {
	fmt.Fprint(out, "foo")
}
func printBar() {
	fmt.Fprint(out, "bar")
}

func (f *FooBar) foo() {
	defer f.wg.Done()

	for i := 0; i < f.n; i++ {
		// Wait till bar is done.
		<-f.barDone

		printFoo()

		// Notify that foo is done.
		f.fooDone <- struct{}{}
	}

	// Clear the last object inserted into barDone,
	// otherwise bar function cannot exit and result in a deadlock
	<-f.barDone
}

func (f *FooBar) bar() {
	defer f.wg.Done()

	// Let function foo execute first.
	f.barDone <- struct{}{}

	for i := 0; i < f.n; i++ {
		// Wait till foo is done.
		<-f.fooDone

		printBar()

		// Notify that foo is done.
		f.barDone <- struct{}{}
	}
}

// run starts the goroutines.
func run(n int) {
	f := FooBar{
		n:       n,
		wg:      new(sync.WaitGroup),
		fooDone: make(chan struct{}),
		barDone: make(chan struct{}),
	}

	// Totally 2 goroutines are triggered
	f.wg.Add(2)

	go f.foo()
	go f.bar()

	// Wait for all goroutines to finish.
	f.wg.Wait()
}

func main() {
	run(10)
}
