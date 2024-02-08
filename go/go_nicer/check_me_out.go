package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func printDart() {
	fmt.Println("    __________                                    ")
	fmt.Println("   /M\\\\\\M|||M//.                                  ")
	fmt.Println("  /MMM\\\\\\|||///M:.                                ")
	fmt.Println(" /MMMMM\\\\\\ | //MMMM:.            ______________________ ")
	fmt.Println("(=========+======<]]]]::::::::::<|||_|||_|||_|||_|||_|||>=========-")
	fmt.Println(" \\#MMMM// | \\\\MMMM:'                              ")
	fmt.Println("  \\#MM///|||\\\\\\M:'                                 ")
	fmt.Println("   \\M///M|||M\\\\'                                  ")
	fmt.Println("")
}

type dartsSolver struct {
	target     int
	totalDarts int

	lenDoubles int
	lenDarts   int
	allDoubles []int
	allDarts   []int

	totalSolutions atomic.Int32
}

func newDartsSolver(target, totalDarts int) *dartsSolver {
	s := &dartsSolver{
		target:     target,
		totalDarts: totalDarts,
		allDoubles: []int{50, 40, 38, 36, 34, 32, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4, 2},
		allDarts:   []int{60, 57, 54, 51, 50, 48, 45, 42, 40, 39, 38, 36, 34, 33, 32, 30, 28, 27, 26, 25, 24, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
	}

	s.lenDoubles = len(s.allDoubles)
	s.lenDarts = len(s.allDarts)

	return s
}

func (s *dartsSolver) reducePossibleDarts() {
	// If the total score required is high we can skip
	// low scoring darts that couldn't feasibly be used
	// to reach the total.
	approxLowestNormalDart := (s.target - 50) - ((s.totalDarts - 2) * 60)
	approxLowestFinishDart := s.target - ((s.totalDarts - 1) * 60)
	if approxLowestNormalDart > 0 {
		for i := 0; i < s.lenDarts; i++ {
			if s.allDarts[i] < approxLowestNormalDart {
				fmt.Printf("Capping possible darts at %d\n", approxLowestNormalDart)
				s.lenDarts = i
			}
		}
	}

	if approxLowestFinishDart > 0 {
		for i := 0; i < s.lenDoubles; i++ {
			if s.allDoubles[i] < approxLowestFinishDart {
				fmt.Printf("Capping possible finish dart at %d\n", approxLowestFinishDart)
				s.lenDoubles = i
			}
		}
	}
}

func printCheckout(darts []int, dartCount int) {
	fmt.Print("Checkout found: ")
	for i := dartCount; i > 0; i-- {
		fmt.Printf("%02d ", darts[i-1])
	}
	fmt.Println("")
}

func (s *dartsSolver) normalDarts(points int, darts []int, dartCount int) {
	//fmt.Printf("normalDarts: %d %d %v\n", points, dartCount, darts)
	// normal darts are any dart which can be used when not
	// checking out.
	if points == 0 {
		// printCheckout(darts, dartCount)
		s.totalSolutions.Add(1)
		return
	}

	if points < 0 || dartCount == s.totalDarts {
		// Invalid state, return
		return
	}

	for i := 0; i < s.lenDarts; i++ {
		darts[dartCount] = s.allDarts[i]
		s.normalDarts(points-s.allDarts[i], darts, dartCount+1)
	}
}

func (s *dartsSolver) findDartCheckouts() {
	// Concurrently look at all checkout darts and recursively
	// search all of the other darts to find possible combinations
	// that reach the total

	s.totalSolutions.Store(0)

	s.reducePossibleDarts()

	var wg sync.WaitGroup

	for i := 0; i < s.lenDoubles; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			darts := make([]int, s.totalDarts)
			fmt.Printf("=== Calculating for finishing dart %d ===\n", s.allDoubles[i])
			darts[0] = s.allDoubles[i]
			s.normalDarts(s.target-s.allDoubles[i], darts, 1)
		}(i)
	}

	wg.Wait()
}

func main() {
	totalDarts := 3
	target := 125

	// takes args from stdin of target and then total_darts, otherwise sets a default value
	if len(os.Args) == 3 {
		var err error

		target, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Error: invalid target %s: %v\n", os.Args[1], err)
			os.Exit(1)
		}

		totalDarts, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Error: invalid totalDarts %s: %v\n", os.Args[2], err)
			os.Exit(1)
		}
	}

	printDart()

	fmt.Printf("Finding %d dart finishes for %d\n", totalDarts, target)

	ts := time.Now()

	s := newDartsSolver(target, totalDarts)

	s.findDartCheckouts()

	duration := time.Since(ts)

	fmt.Printf("Total solutions: %d\n", s.totalSolutions.Load())
	fmt.Printf("%.2fs\n", duration.Seconds())
}
