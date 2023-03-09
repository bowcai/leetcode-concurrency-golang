package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type h2O struct {
	wg              *sync.WaitGroup // Wait for all the go routine to finish
	hydrogenVacancy chan struct{}   // 1 hydrogen need 1 vacancy to release, and create 1 vacancy for oxygen
	oxygenVacancy   chan struct{}   // 1 oxygen need 2 vacancies to release, and create 2 vacancies for hydrogen
	oxygenMu        sync.Mutex      // oxygenMu guard the acquiring of 2 oxygenVacancy
}

// Specify the output of printing.
// Will be modified during testing.
var out io.Writer = os.Stdout

// Use a mutex to guard the printing.
// May be useless if printing is thread-safe.
var printMu sync.Mutex

// The functions to handle the printing defined in the problem
func releaseHydrogen() {
	printMu.Lock()
	defer printMu.Unlock()

	if _, err := fmt.Fprint(out, "H"); err != nil {
		panic(err)
	}
}

func releaseOxygen() {
	printMu.Lock()
	defer printMu.Unlock()

	if _, err := fmt.Fprint(out, "O"); err != nil {
		panic(err)
	}
}

func (h *h2O) hydrogen(releaseHydrogen func()) {
	defer h.wg.Done()

	// 1 hydrogen need 1 vacancy to release,
	// otherwise it will block until there is vacancy.
	h.hydrogenVacancy <- struct{}{}

	releaseHydrogen()

	// 1 hydrogen can create 1 vacancy for oxygen,
	// so that the oxygen can release.
	<-h.oxygenVacancy
}

func (h *h2O) oxygen(releaseOxygen func()) {
	defer h.wg.Done()

	// 1 oxygen need 2 vacancies to release,
	// which corresponds to the release of 2 hydrogen atoms.
	// Use a mutex to guard the simultaneously acquiring of 2 vacancies.
	h.oxygenMu.Lock()
	h.oxygenVacancy <- struct{}{}
	h.oxygenVacancy <- struct{}{}
	h.oxygenMu.Unlock()

	releaseOxygen()

	// 1 oxygen can create 2 vacancies for hydrogen.
	<-h.hydrogenVacancy
	<-h.hydrogenVacancy
}

func Run(water string) {
	obj := h2O{
		wg:              new(sync.WaitGroup),
		hydrogenVacancy: make(chan struct{}, 2),
		oxygenVacancy:   make(chan struct{}, 2),
		oxygenMu:        sync.Mutex{},
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
	Run("OOHHHH")
}
