package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type zeroEvenOdd struct {
	n           int
	evenChan    chan int      // Channel for function even to get the number to print
	oddChan     chan int      // Channel for function odd the get the number to print
	evenOddDone chan struct{} // Notify zero that even/odd is done
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
func (z *zeroEvenOdd) zero(printNumber func(int)) {
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

func (z *zeroEvenOdd) even(printNumber func(int)) {
	// The block is executed when it gets number from evenChan.
	// Loop until evenChan is closed.
	for num := range z.evenChan {
		printNumber(num)

		// Notify that even is done.
		z.evenOddDone <- struct{}{}
	}
}

func (z *zeroEvenOdd) odd(printNumber func(int)) {
	// The block is executed when it gets number from oddChan.
	// Loop until oddChan is closed.
	for num := range z.oddChan {
		printNumber(num)

		// Notify that odd is done.
		z.evenOddDone <- struct{}{}
	}
}

func Run(n int) {
	obj := zeroEvenOdd{
		n:           n,
		evenChan:    make(chan int),
		oddChan:     make(chan int),
		evenOddDone: make(chan struct{}),
	}

	var wg sync.WaitGroup

	// Totally 4 goroutines are triggered.
	wg.Add(3)

	go func() {
		defer wg.Done()
		obj.zero(printNumber)
	}()
	go func() {
		defer wg.Done()
		obj.even(printNumber)
	}()
	go func() {
		defer wg.Done()
		obj.odd(printNumber)
	}()

	// Wait for all the goroutines to finish.
	wg.Wait()
}

func main() {
	Run(9)
}
