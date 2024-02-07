package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	lenDoubles = 21
	lenDarts   = 43
	DOUBLES    = []int{50, 40, 38, 36, 34, 32, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4, 2}
	DARTS      = []int{60, 57, 54, 51, 50, 48, 45, 42, 40, 39, 38, 36, 34, 33, 32, 30, 28, 27, 26, 25, 24, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
)

var target, totalDarts, totalSolutions int

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

func reducePossibleDarts() {
	// If the total score required is high we can skip
	// low scoring darts that couldn't feasibly be used
	// to reach the total.
	approxLowestNormalDart := (target - 50) - ((totalDarts - 2) * 60)
	approxLowestFinishDart := target - ((totalDarts - 1) * 60)
	if approxLowestNormalDart > 0 {
		for i := 0; i < lenDarts; i++ {
			if DARTS[i] < approxLowestNormalDart {
				fmt.Printf("Capping possible darts at %d\n", approxLowestNormalDart)
				lenDarts = i
			}
		}
	}

	if approxLowestFinishDart > 0 {
		for i := 0; i < lenDoubles; i++ {
			if DOUBLES[i] < approxLowestFinishDart {
				fmt.Printf("Capping possible finish dart at %d\n", approxLowestFinishDart)
				lenDoubles = i
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

func normalDarts(points int, darts []int, dartCount int) {
	//fmt.Printf("normalDarts: %d %d %v\n", points, dartCount, darts)
	// normal darts are any dart which can be used when not
	// checking out.
	if points == 0 {
		// printCheckout(darts, dartCount)
		totalSolutions++
		return
	}

	if points < 0 || dartCount == totalDarts {
		// Invalid state, return
		return
	}

	for i := 0; i < lenDarts; i++ {
		darts[dartCount] = DARTS[i]
		normalDarts(points-DARTS[i], darts, dartCount+1)
	}
}

func findDartCheckouts() int {
	// Loop over the possible checkout darts, then recursively
	// search all of the other darts to find possible combinations
	// that reach the total

	darts := make([]int, totalDarts)
	totalSolutions = 0

	reducePossibleDarts()

	for i := 0; i < lenDoubles; i++ {
		fmt.Printf("=== Calculating for finishing dart %d ===\n", DOUBLES[i])
		darts[0] = DOUBLES[i]
		normalDarts(target-DOUBLES[i], darts, 1)
	}

	return totalSolutions
}

func main() {
	// takes args from stdin of target and then total_darts, otherwise sets a default value
	if len(os.Args) == 3 {
		var err error

		target, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}

		totalDarts, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	} else {
		target = 125
		totalDarts = 3
	}

	printDart()
	fmt.Printf("Finding %d dart finishes for %d\n", totalDarts, target)

	ts := time.Now()

	findDartCheckouts()

	duration := time.Since(ts)

	fmt.Printf("Total solutions: %d\n", totalSolutions)
	fmt.Printf("%.2fs\n", duration.Seconds())
}
