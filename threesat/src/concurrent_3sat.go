package main

import (
	"fmt"
	"math"
	"sync"
	"time"
	"runtime"
)

func readClauses(nClauses int) [][]int16 {
	clauses := make([][]int16, 3)
	clauses[0] = make([]int16, nClauses)
	clauses[1] = make([]int16, nClauses)
	clauses[2] = make([]int16, nClauses)

	for i := 0; i < nClauses; i++ {
		fmt.Scanf("%d %d %d", &clauses[0][i], &clauses[1][i], &clauses[2][i])
	}

	return clauses
}

func solveClauses(clauses [][]int16, nClauses, nVar int) int64 {
	iVar := make([]int64, nVar)
	for i := 0; i < nVar; i++ {
		iVar[i] = int64(math.Pow(2, float64(i)))
	}

	var maxNumber int64 = int64(math.Pow(2, float64(nVar)))
	var wg sync.WaitGroup
	wg.Add(maxNumber)

	resultChan := make(chan int64, maxNumber)

	var i int64
	for i = 0; i < maxNumber; i++ {
		go func(clauses [][]int16, solNum int64, resChan chan int64) {
			defer wg.Done()
			var c int
			var variable int16
			for c = 0 ; c < nClauses ; c++ {

				variable = clauses[0][c]
				if variable > 0 && (solNum&iVar[variable-1]) > 0 {
					continue
				} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
					continue
				}

				variable = clauses[1][c]
				if variable > 0 && (solNum&iVar[variable-1]) > 0 {
					continue
				} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
					continue
				}

				variable = clauses[2][c]
				if variable > 0 && (solNum&iVar[variable-1]) > 0 {
					continue
				} else if variable < 0 && (solNum&iVar[-variable-1] == 0) {
					continue
				}
				break
			}

			if c == nClauses {
				resChan <- solNum
			}else {
				resChan <- -1
			}


		}(clauses, i, resultChan)
	}

	go func(wg sync.WaitGroup, resultChan chan int64) {
		wg.Wait()
		close(resultChan)
	}(wg, resultChan)

	for result := range resultChan {
		if result > 0 {
			return result
		}
	}
	return -1
}

func main() {
	runtime.GOMAXPROCS(1)
	var nClauses, nVar int
	fmt.Scanf("%d %d", &nClauses, &nVar)

	clauses := readClauses(nClauses)

	startTime := time.Now()
	solution := solveClauses(clauses,nClauses, nVar)
	fmt.Printf("Time to solve: %s\n",time.Since(startTime))

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
