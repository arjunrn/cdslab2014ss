package main

import (
	"fmt"
	"math"
	"time"
	"runtime"
	"sync"
	"os"
)

var workerCount int64
var nClauses, nVar int
var clauses [][]int16
var iVar []int64
var maxNumber int64
var solution int64 = -1

func clausesSolver(ith int64, wg *sync.WaitGroup) {
	defer wg.Done()

	var variable int16
	var solNum, i int64
	for solNum = ith; solNum < maxNumber && solution < 0; solNum+=workerCount {
		for i = 0; i < int64(nClauses); i++ {
			variable = clauses[0][i]
			if variable > 0 && (solNum&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
				continue
			}

			variable = clauses[1][i]
			if variable > 0 && (solNum&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
				continue
			}

			variable = clauses[2][i]
			if variable > 0 && (solNum&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
				continue
			}
			break
		}

		if i == int64(nClauses) && solution < 0 {
			solution = solNum
		}

		if solution >= 0 {
			break
		}
	}
}

func solveClauses() int64 {
	iVar = make([]int64, nVar)
	for i := 0; i < nVar; i++ {
		iVar[i] = int64(math.Pow(2, float64(i)))
	}

	maxNumber = int64(math.Pow(2, float64(nVar)))
	var wg sync.WaitGroup
	var i int64
	for i = 0 ; i < workerCount; i++ {
		wg.Add(1)
		go clausesSolver(int64(i), &wg)
	}

	wg.Wait()
	return solution

}

func readClauses() {
	clauses = make([][]int16, 3)
	clauses[0] = make([]int16, nClauses)
	clauses[1] = make([]int16, nClauses)
	clauses[2] = make([]int16, nClauses)

	for i := 0; i < nClauses; i++ {
		fmt.Scanf("%d %d %d", &clauses[0][i], &clauses[1][i], &clauses[2][i])
	}
}

func main() {
	workerCount = int64(runtime.NumCPU())
	runtime.GOMAXPROCS(int(workerCount))

	logHandler , logErr := os.Create("logfile.log")
	if logErr != nil {
		panic(logErr)
	}
	defer logHandler.Close()

	fmt.Scanf("%d %d", &nClauses, &nVar)

	readClauses()

	startTime := time.Now()
	solution := solveClauses()
	fmt.Fprintf(logHandler, "Time to solve: %s\n", time.Since(startTime))

	if solution > 0 {
		fmt.Printf("Solution found [%d]: ", solution)
		for i := 0; i < nVar; i++ {
			fmt.Printf("%d ", int((solution & int64(math.Pow(2, float64(i)))) / int64(math.Pow(2, float64(i)))))
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Solution not found. \n")
	}
}
