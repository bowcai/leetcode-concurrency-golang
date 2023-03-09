package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type ZeroEvenOdd struct {
	n           int
	wg          *sync.WaitGroup // Wait for all the go routine to finish
	evenChan    chan int        // Channel for function even to get the number to print
	oddChan     chan int        // Channel for function odd the get the number to print
	evenOddDone chan struct{}   // Notify zero that even/odd is done
}

// Specify the output of printing
// Will be modified during testing
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem
func printNumber(x int) {
	if _, err := fmt.Fprint(out, x); err != nil {
		panic(err)
	}
}

// zero loops over 1-n and print 0,
// and notifies the other two goroutines to print numbers
func (z *ZeroEvenOdd) zero(printNumber func(int)) {
	defer z.wg.Done()

	// Loop over 1-n
	for i := 1; i <= z.n; i++ {
		// Print 0
		printNumber(0)

		// Insert the number to the corresponding channel
		// and let the function print the number
		if i%2 == 0 {
			z.evenChan <- i
		} else {
			z.oddChan <- i
		}

		// Wait for even/odd function done.
		<-z.evenOddDone
	}

	// Close the channels to notify other goroutines to exit.
	close(z.evenChan)
	close(z.oddChan)
}

func (z *ZeroEvenOdd) even(printNumber func(int)) {
	defer z.wg.Done()

	// The block is executed when it gets number from evenChan.
	// Loop until evenChan is closed.
	for num := range z.evenChan {
		printNumber(num)

		// Notify that even is done.
		z.evenOddDone <- struct{}{}
	}
}

func (z *ZeroEvenOdd) odd(printNumber func(int)) {
	defer z.wg.Done()

	// The block is executed when it gets number from oddChan.
	// Loop until oddChan is closed.
	for num := range z.oddChan {
		printNumber(num)

		// Notify that odd is done.
		z.evenOddDone <- struct{}{}
	}
}

func run(n int) {
	obj := ZeroEvenOdd{
		wg:          new(sync.WaitGroup),
		n:           n,
		evenChan:    make(chan int),
		oddChan:     make(chan int),
		evenOddDone: make(chan struct{}),
	}

	// Totally 4 goroutines are triggered.
	obj.wg.Add(3)

	go obj.zero(printNumber)
	go obj.even(printNumber)
	go obj.odd(printNumber)

	// Wait for all the goroutines to finish.
	obj.wg.Wait()
}

func main() {
	run(9)
}
