package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type H2O struct {
	wg              *sync.WaitGroup // Wait for all the go routine to finish
	hydrogenVacancy chan struct{}   // 1 hydrogen need 1 vacancy to release, and create 1 vacancy for oxygen
	oxygenVacancy   chan struct{}   // 1 oxygen need 2 vacancies to release, and create 2 vacancies for hydrogen
}

// Specify the output of printing.
// Will be modified during testing.
var out io.Writer = os.Stdout

// Use a mutex to guard the printing.
// May be useless if printing is thread-safe.
var mu sync.Mutex

// The functions to handle the printing defined in the problem
// For simplicity, error is not handled here
func releaseHydrogen() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprint(out, "H")
}

func releaseOxygen() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprint(out, "O")
}

func (h *H2O) hydrogen(releaseHydrogen func()) {
	defer h.wg.Done()

	// 1 hydrogen need 1 vacancy to release,
	// otherwise it will block until there is vacancy.
	h.hydrogenVacancy <- struct{}{}

	releaseHydrogen()

	// 1 hydrogen can create 1 vacancy for oxygen,
	// so that the oxygen can release.
	<-h.oxygenVacancy
}

func (h *H2O) oxygen(releaseOxygen func()) {
	defer h.wg.Done()

	// 1 oxygen need 2 vacancies to release,
	// which corresponds to the release of 2 hydrogen atoms.
	h.oxygenVacancy <- struct{}{}
	h.oxygenVacancy <- struct{}{}

	releaseOxygen()

	// 1 oxygen can create 2 vacancies for hydrogen.
	<-h.hydrogenVacancy
	<-h.hydrogenVacancy
}

func run(water string) {
	obj := H2O{
		wg:              new(sync.WaitGroup),
		hydrogenVacancy: make(chan struct{}, 2),
		oxygenVacancy:   make(chan struct{}, 2),
	}

	// Create a goroutine for each atom in the input.
	for _, c := range water {
		switch c {
		case 'H':
			obj.wg.Add(1)
			go obj.hydrogen(releaseHydrogen)
		case 'O':
			obj.wg.Add(1)
			go obj.oxygen(releaseOxygen)
		default:
			panic("Input should only contain 'H' and 'O'")
		}
	}

	// Wait for all the goroutines to finish.
	obj.wg.Wait()
}

func main() {
	run("OOHHHH")
}
