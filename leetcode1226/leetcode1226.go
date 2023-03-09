package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	numPhilosopher = 5                      // The total number of philosophers
	thinkTime      = 100 * time.Millisecond // The amount of time a philosopher think
)

// Specify the output of printing.
// Can be modified during testing.
var out io.Writer = os.Stdout

// Use a mutex to guard the printing.
// May be useless if printing is thread-safe.
var printMu sync.Mutex

// The functions to handle the printing defined in the problem.
func eat(philosopher int) {
	printMu.Lock()
	defer printMu.Unlock()

	// The printing defined in the LeetCode problem
	// _, err := fmt.Fprintf(out, "[%d, 0, 3] ", philosopher)

	// Use self-defined printing for a better understanding
	_, err := fmt.Fprintf(out, "%d is eating.\n", philosopher)

	if err != nil {
		panic(err)
	}
}

func pickLeftFork(philosopher int) {
	printMu.Lock()
	defer printMu.Unlock()

	// The printing defined in the LeetCode problem
	// _, err := fmt.Fprintf(out, "[%d, 1, 1] ", philosopher)

	// Use self-defined printing for a better understanding
	_, err := fmt.Fprintf(out, "%d is picking left fork.\n", philosopher)

	if err != nil {
		panic(err)
	}
}

func pickRightFork(philosopher int) {
	printMu.Lock()
	defer printMu.Unlock()

	// The printing defined in the LeetCode problem
	// _, err := fmt.Fprintf(out, "[%d, 2, 1] ", philosopher)

	// Use self-defined printing for a better understanding
	_, err := fmt.Fprintf(out, "%d is picking right fork.\n", philosopher)

	if err != nil {
		panic(err)
	}
}

func putLeftFork(philosopher int) {
	printMu.Lock()
	defer printMu.Unlock()

	// The printing defined in the LeetCode problem
	// _, err := fmt.Fprintf(out, "[%d, 1, 2] ", philosopher)

	// Use self-defined printing for a better understanding
	_, err := fmt.Fprintf(out, "%d is putting left fork.\n", philosopher)

	if err != nil {
		panic(err)
	}
}

func putRightFork(philosopher int) {
	printMu.Lock()
	defer printMu.Unlock()

	// The printing defined in the LeetCode problem
	// _, err := fmt.Fprintf(out, "[%d, 2, 2] ", philosopher)

	// Use self-defined printing for a better understanding
	_, err := fmt.Fprintf(out, "%d is putting right fork.\n", philosopher)

	if err != nil {
		panic(err)
	}
}

// think represents the philosopher is thinking.
// Do nothing or sleep for a short while.
func think(philosopher int) {
	// Sleep for a short while.
	time.Sleep(thinkTime)
}

type diningPhilosophers struct {
	wg      *sync.WaitGroup // Wait for all the goroutines to finish.
	forksMu [numPhilosopher]sync.Mutex
}

// This implementation lets the philosopher pick the fork in a lowest-first order.
// It does not restrict the number of eating philosophers.
func (p *diningPhilosophers) wantsToEat(
	philosopher int,
	pickLeftFork func(),
	pickRightFork func(),
	eat func(),
	putLeftFork func(),
	putRightFork func(),
) {
	leftForkId := philosopher
	rightForkId := (philosopher + 1) % numPhilosopher

	var smallForkId, largeForkId int
	var pickSmallFork, pickLargeFork func()

	if leftForkId <= rightForkId {
		smallForkId, largeForkId = leftForkId, rightForkId
		pickSmallFork, pickLargeFork = pickLeftFork, pickRightFork
	} else {
		smallForkId, largeForkId = rightForkId, leftForkId
		pickSmallFork, pickLargeFork = pickRightFork, pickLeftFork
	}

	// Lock the small fork first.
	p.forksMu[smallForkId].Lock()
	defer p.forksMu[smallForkId].Unlock()
	pickSmallFork()

	// And then lock the large fork.
	p.forksMu[largeForkId].Lock()
	defer p.forksMu[largeForkId].Unlock()
	pickLargeFork()

	// Have both forks, start to eat.
	eat()

	// Put both forks.
	// The mutexes will be released by defer statements.
	putLeftFork()
	putRightFork()
}

// runPhilosopher represents a philosopher that eat and think for n times.
func (p *diningPhilosophers) runPhilosopher(philosopher, n int) {
	defer p.wg.Done()

	pickLeftForkFunc := func() { pickLeftFork(philosopher) }
	pickRightForkFunc := func() { pickRightFork(philosopher) }
	eatFunc := func() { eat(philosopher) }
	putLeftForkFunc := func() { putLeftFork(philosopher) }
	putRightForkFunc := func() { putRightFork(philosopher) }

	for i := 0; i < n; i++ {
		p.wantsToEat(
			philosopher,
			pickLeftForkFunc,
			pickRightForkFunc,
			eatFunc,
			putLeftForkFunc,
			putRightForkFunc,
		)

		think(philosopher)
	}
}

// Run starts the whole program.
// n is the number of times each philosopher need to eat.
func Run(n int) {
	obj := diningPhilosophers{
		wg: new(sync.WaitGroup),
	}

	// Start the goroutine of each philosopher.
	for i := 0; i < numPhilosopher; i++ {
		obj.wg.Add(1)
		go obj.runPhilosopher(i, n)
	}

	// Wait for all the goroutines to finish.
	obj.wg.Wait()
}

func main() {
	Run(2)
}
