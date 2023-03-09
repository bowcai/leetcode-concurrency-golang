package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type fooBar struct {
	n       int
	fooDone chan struct{} // Notify the function foo is done
	barDone chan struct{} // Notify the function bar is done
}

// Specify the output of printing
// Will be modified during testing
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem
func printFoo() {
	if _, err := fmt.Fprint(out, "foo"); err != nil {
		panic(err)
	}
}
func printBar() {
	if _, err := fmt.Fprint(out, "bar"); err != nil {
		panic(err)
	}
}

func (f *fooBar) foo(printFoo func()) {
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

func (f *fooBar) bar(printBar func()) {
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

// Run starts the goroutines.
func Run(n int) {
	f := fooBar{
		n:       n,
		fooDone: make(chan struct{}),
		barDone: make(chan struct{}),
	}

	var wg sync.WaitGroup

	// Totally 2 goroutines are triggered
	wg.Add(2)

	go func() {
		defer wg.Done()
		f.foo(printFoo)
	}()
	go func() {
		defer wg.Done()
		f.bar(printBar)
	}()

	// Wait for all goroutines to finish.
	wg.Wait()
}

func main() {
	Run(10)
}
