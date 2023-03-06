package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type FizzBuzz struct {
	wg             *sync.WaitGroup // Wait for all the goroutines to finish.
	n              int             // n represents the length of sequence to be printed.
	fizzRun        chan struct{}   // Notify fizz function to execute.
	buzzRun        chan struct{}   // Notify buzz function to execute.
	fizzbuzzRun    chan struct{}   // Notify fizzbuzz function to execute.
	fizzOrBuzzDone chan struct{}   // Notify number function that fizz/buzz/fizzbuzz is done.
}

// Specify the output of printing.
// Will be modified during testing.
var out io.Writer = os.Stdout

// The functions to handle the printing defined in the problem.
// For simplicity, error is not handled here.
func printFizz() {
	fmt.Fprint(out, "fizz ")
}
func printBuzz() {
	fmt.Fprint(out, "buzz ")
}
func printFizzBuzz() {
	fmt.Fprint(out, "fizzbuzz ")
}
func printNumber(x int) {
	fmt.Fprint(out, x, " ")
}

func (f *FizzBuzz) fizz() {
	defer f.wg.Done()

	// The block is executed when it can get message from fizzRun channel.
	// Loop until fizzRun is closed.
	for range f.fizzRun {
		printFizz()

		// Notify that fizz is done.
		f.fizzOrBuzzDone <- struct{}{}
	}
}

func (f *FizzBuzz) buzz() {
	defer f.wg.Done()

	// The block is executed when it can get message from buzzRun channel.
	// Loop until buzzRun is closed.
	for range f.buzzRun {
		printBuzz()

		// Notify that buzz is done.
		f.fizzOrBuzzDone <- struct{}{}
	}
}

func (f *FizzBuzz) fizzbuzz() {
	defer f.wg.Done()

	// The block is executed when it can get message from fizzbuzzRun channel.
	// Loop until fizzbuzzRun is closed.
	for range f.fizzbuzzRun {
		printFizzBuzz()

		// Notify that fizzbuzz is done.
		f.fizzOrBuzzDone <- struct{}{}
	}
}

// number handles the loop over numbers until n.
// It will notify fizz/buzz/fizzbuzz goroutines
// when they need to print the corresponding strings.
// It will also notify these goroutines that the loop is finished
// by closing the channels.
func (f *FizzBuzz) number() {
	defer f.wg.Done()

	for i := 1; i <= f.n; i++ {
		if i%3 == 0 {
			if i%5 == 0 { // i is divisible by 3 and 5.
				f.fizzbuzzRun <- struct{}{}

				// Wait for fizzbuzz to finish.
				<-f.fizzOrBuzzDone
			} else { // i is divisible by 3 but not by 5.
				f.fizzRun <- struct{}{}

				// Wait for fizz to finish.
				<-f.fizzOrBuzzDone
			}
		} else { // i is not divisible by 3.
			if i%5 == 0 { // i is divisible by 5 but not by 3.
				f.buzzRun <- struct{}{}

				// Wait for buzz to finish.
				<-f.fizzOrBuzzDone
			} else { // i is not divisible by 3 or 5.
				printNumber(i)
			}
		}
	}

	// Close the channels to notify other goroutines to exit.
	close(f.fizzRun)
	close(f.buzzRun)
	close(f.fizzbuzzRun)
}

func run(n int) {
	obj := FizzBuzz{
		wg:             new(sync.WaitGroup),
		n:              n,
		fizzRun:        make(chan struct{}),
		buzzRun:        make(chan struct{}),
		fizzbuzzRun:    make(chan struct{}),
		fizzOrBuzzDone: make(chan struct{}),
	}

	// Totally 4 goroutines are triggered.
	obj.wg.Add(4)

	go obj.fizz()
	go obj.buzz()
	go obj.fizzbuzz()
	go obj.number()

	// Wait for all the goroutines to finish.
	obj.wg.Wait()
}

func main() {
	run(15)
}