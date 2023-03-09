package main

import (
	"sync"
)

// The maximum number of philosopher allowed to eat simultaneously
const maxEatingPhilosopher = 3

// The second implementation of dining philosopher
type diningPhilosophersImpl2 struct {
	forksMu [numPhilosopher]sync.Mutex
	eatSema chan struct{} // Semaphore of philosophers that allowed to eat
}

// This implementation restrict the number of philosopher
// allowed to visit wantsToEat simultaneously
func (p *diningPhilosophersImpl2) wantsToEat(
	philosopher int,
	pickLeftFork func(),
	pickRightFork func(),
	eat func(),
	putLeftFork func(),
	putRightFork func(),
) {
	// Acquire one semaphore of eating vacancy
	p.eatSema <- struct{}{}

	leftForkId := philosopher
	rightForkId := (philosopher + 1) % numPhilosopher

	// Pick the left fork.
	p.forksMu[leftForkId].Lock()
	defer p.forksMu[leftForkId].Unlock()
	pickLeftFork()

	// Pick the right fork.
	p.forksMu[rightForkId].Lock()
	defer p.forksMu[rightForkId].Unlock()
	pickRightFork()

	// Have both forks, start to eat.
	eat()

	// Put both forks.
	// The mutexes will be released by defer statements.
	putLeftFork()
	putRightFork()

	// Release the semaphore
	<-p.eatSema
}

// RunImpl2 starts the whole program of the second implementation.
// n is the number of times each philosopher need to eat.
func RunImpl2(n int) {
	obj2 := diningPhilosophersImpl2{
		// Use a channel to simulate a semaphore with size maxEatingPhilosopher
		eatSema: make(chan struct{}, maxEatingPhilosopher),
	}

	var wg sync.WaitGroup

	// Start the goroutine of each philosopher.
	for i := 0; i < numPhilosopher; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runPhilosopher(i, n, obj2.wantsToEat, think)
		}(i)
	}

	// Wait for all the goroutines to finish.
	wg.Wait()
}
