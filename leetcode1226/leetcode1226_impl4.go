package main

import "sync"

// The fourth implementation of dining philosopher.
// In this implementation, use boolean array to check the status of each fork.
// Use a lock to guard the fork array,
// and use condition variable to notify the blocked philosophers to check the fork status again.
type diningPhilosophersImpl4 struct {
	forksTaken    [numPhilosopher]bool
	forkArrayMu   sync.Mutex
	forkArrayCond *sync.Cond
}

func (p *diningPhilosophersImpl4) wantsToEat(
	philosopher int,
	pickLeftFork func(),
	pickRightFork func(),
	eat func(),
	putLeftFork func(),
	putRightFork func(),
) {
	leftForkId := philosopher
	rightForkId := (philosopher + 1) % numPhilosopher

	// Check if the two forks are both available.
	// If not, blocked and wait for any other philosophers to put the forks.
	p.forkArrayMu.Lock()
	for p.forksTaken[leftForkId] || p.forksTaken[rightForkId] {
		p.forkArrayCond.Wait()
	}

	// Release the lock after setting status of both forks,
	// so that other philosopher can check the status.
	func() {
		defer p.forkArrayMu.Unlock()

		// Pick the forks and set the status.
		p.forksTaken[leftForkId] = true
		pickLeftFork()
		p.forksTaken[rightForkId] = true
		pickRightFork()
	}()

	// Have both forks, start to eat.
	eat()

	// Put both forks by setting the status.
	p.forksTaken[leftForkId] = false
	putLeftFork()
	p.forksTaken[rightForkId] = false
	putRightFork()

	// Notify other blocked philosophers to check the fork status.
	p.forkArrayCond.Broadcast()
}

// RunImpl4 starts the whole program of the fourth implementation.
// n is the number of times each philosopher need to eat.
func RunImpl4(n int) {
	obj4 := diningPhilosophersImpl4{}
	obj4.forkArrayCond = sync.NewCond(&obj4.forkArrayMu)

	var wg sync.WaitGroup

	// Start the goroutine of each philosopher.
	for i := 0; i < numPhilosopher; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runPhilosopher(i, n, obj4.wantsToEat, think)
		}(i)
	}

	// Wait for all the goroutines to finish.
	wg.Wait()
}
