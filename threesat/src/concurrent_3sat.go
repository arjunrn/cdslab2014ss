package main

import (
	"fmt"
	"math"
	"time"
	"runtime"
	"sync"
)

const CORE_COUNT = 4

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
					if !allCondTrue{
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
	solution := solveClauses(clauses, nClauses, nVar)
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
