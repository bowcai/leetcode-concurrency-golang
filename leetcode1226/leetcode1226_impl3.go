package main

import "sync"

// The third implementation of dining philosopher.
// In this implementation, each philosopher try to lock two forks simultaneously.
type diningPhilosophersImpl3 struct {
	forksMu [numPhilosopher]sync.Mutex
}

// simulLock locks multiple mutexes simultaneously, resembling std::lock in C++.
// If any of the mutexes cannot be locked, the function get blocked.
// The function does not hold any locked mutex until successfully locking all the mutexes.
func simulLock(mutexes ...*sync.Mutex) {
	lenLocks := len(mutexes)

	// Check if the parameter list is empty.
	if lenLocks == 0 {
		return
	}

	// first is the index of the first mutex to be locked.
	// The function will get blocked if this mutex cannot be locked
	// (i.e. using Lock() instead of TryLock()).
	first := 0

	for {
		// Lock the mutex with first index.
		mutexes[first].Lock()

		// First mutex get locked. Try to lock other mutexes with TryLock().
		// If any mutex cannot be locked, release all the mutexes that already get locked,
		// and change the first index to wait for the failed mutex.
		var j int
		for j = 1; j < lenLocks; j++ {
			idx := (first + j) % lenLocks

			// Try to lock the mutex. If fails, release all mutexes.
			if !mutexes[idx].TryLock() {
				// Release all mutexes already locked.
				for k := j; k != 0; k-- {
					mutexes[(first+k-1)%lenLocks].Unlock()
				}

				// Change the index first to the failed mutex.
				// Will blocked to wait for this mutex.
				first = idx

				break
			}
		}

		// All mutexes are locked successfully,
		// break the for loop and return.
		if j == lenLocks {
			break
		}
	}
}

func (p *diningPhilosophersImpl3) wantsToEat(
	philosopher int,
	pickLeftFork func(),
	pickRightFork func(),
	eat func(),
	putLeftFork func(),
	putRightFork func(),
) {
	leftForkId := philosopher
	rightForkId := (philosopher + 1) % numPhilosopher

	// Pick two forks simultaneously.
	simulLock(&p.forksMu[leftForkId], &p.forksMu[rightForkId])
	defer p.forksMu[leftForkId].Unlock()
	defer p.forksMu[rightForkId].Unlock()
	pickLeftFork()
	pickRightFork()

	// Have both forks, start to eat.
	eat()

	// Put both forks.
	// The mutexes will be released by defer statements.
	putLeftFork()
	putRightFork()
}

// RunImpl3 starts the whole program of the third implementation.
// n is the number of times each philosopher need to eat.
func RunImpl3(n int) {
	obj3 := diningPhilosophersImpl3{}

	var wg sync.WaitGroup

	// Start the goroutine of each philosopher.
	for i := 0; i < numPhilosopher; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			runPhilosopher(i, n, obj3.wantsToEat, think)
		}(i)
	}

	// Wait for all the goroutines to finish.
	wg.Wait()
}
