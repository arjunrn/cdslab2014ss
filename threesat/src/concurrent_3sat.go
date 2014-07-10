package main

import (
	"fmt"
	"math"
	"time"
	"runtime"
	"sync"
)

const CORE_COUNT = 1

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

func testClause(clauses [][]int16, iVar []int64, nClauses int64, solChan chan int64, numberChan chan int64) {
	var solNum, i int64
	var variable int16
	for {
		solNum = <-numberChan
		if solNum == -1 {
			return
		}
		for i = 0; i < nClauses; i++ {
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

		if i == nClauses {
			solChan <- solNum
		} else {
			solChan <- -1
		}
	}

}

func solveClauses(nClauses int, clauses [][]int16, nVar int) int64 {
	iVar := make([]int64, nVar)
	for i := 0; i < nVar; i++ {
		iVar[i] = int64(math.Pow(2, float64(i)))
	}

	var maxNumber int64 = int64(math.Pow(2, float64(nVar)))

	workerChans := make([]chan int64, CORE_COUNT)

	for i := range workerChans {
		workerChans[i] = make(chan int64, 1)
	}

	solChan := make(chan int64)
	for i := range workerChans {
		go testClause(clauses, iVar, int64(nClauses), solChan, workerChans[i])
	}

	var counter int64 = 0
	var i int64
	var result int64 = -1
	var resultFinal int64 = -1
	for counter < maxNumber {

		for i = 0; i < CORE_COUNT && counter+i < maxNumber; i++ {
			workerChans[i] <- counter
			counter++
		}

		for i = 0; i < CORE_COUNT && counter+i < maxNumber; i++ {
			result = <-solChan
			if result >= 0 {
				resultFinal = result
			}
		}

		if resultFinal >= 0 {
			for i = 0; i < CORE_COUNT; i++ {
				workerChans[i] <- -1
			}
			return resultFinal
		}
	}

	return -1
}

func solveClausesOld(nClauses int, clauses [][]int16, nVar int) int64 {
	iVar := make([]int64, nVar)
	for i := 0; i < nVar; i++ {
		iVar[i] = int64(math.Pow(2, float64(i)))
	}

	var maxNumber int64 = int64(math.Pow(2, float64(nVar)))
	var solNum int64

	coreAssignmentCount := int(nClauses / CORE_COUNT)
	if coreAssignmentCount == 0 {
		coreAssignmentCount = 1
	}
	fmt.Printf("Core Assignment: %d\n", coreAssignmentCount)

	for solNum = 0; solNum < maxNumber; solNum++ {
		var c int

		var wg sync.WaitGroup
		allCondTrue := true

		for c = 0 ; c < nClauses ; c+=coreAssignmentCount {

			wg.Add(1)

			go func(startClause int) {
				defer wg.Done()
				var endClause int
				if startClause+coreAssignmentCount > nClauses {
					endClause = nClauses
				} else {
					endClause = startClause+coreAssignmentCount
				}
				var variable int16
				var i int

				for i = startClause; i < endClause; i++ {
					if !allCondTrue {
						return
					}

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

				if i != endClause {
					allCondTrue = false
				}

			}(c)
		}
		wg.Wait()
		if allCondTrue {
			return solNum
		}

	}
	return -1
}

func main() {
	runtime.GOMAXPROCS(CORE_COUNT)
	var nClauses, nVar int
	fmt.Scanf("%d %d", &nClauses, &nVar)

	clauses := readClauses(nClauses)

	startTime := time.Now()
	solution := solveClauses(nClauses, clauses, nVar)
	fmt.Printf("Time to solve: %s\n", time.Since(startTime))

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
