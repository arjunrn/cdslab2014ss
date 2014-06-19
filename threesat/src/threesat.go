package main

import "fmt"

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
		iVar[i] = int64(2 ^ i)

	}

	var maxNumber int64 = int64(2 ^ nVar)
	var number int64
	var variable int16
	var c int

	for number = 0; number < maxNumber; number++ {
		for c = 0 ; c < nClauses ; c++ {

			variable = clauses[0][c]
			if variable > 0 && (number&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (number&iVar[-variable-1] == 0) {
				continue
			}

			variable = clauses[1][c]
			if variable > 0 && (number&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (number&iVar[-variable-1] == 0) {
				continue
			}

			variable = clauses[2][c]
			if variable > 0 && (number&iVar[variable-1]) > 0 {
				continue
			} else if variable < 0 && (number&iVar[-variable-1] == 0) {
				continue
			}
			break
		}

		if c == nClauses {
			return number
		}

	}
	return -1
}

func main() {
	var nClauses, nVar int
	fmt.Scanf("%d %d", &nClauses, &nVar)

	clauses := readClauses(nClauses)

	solution := solveClauses(clauses, nClauses, nVar)

	if solution > 0 {
		fmt.Printf("Solution found [%d]: ", solution)
		for i := 0; i < nVar; i++ {
			fmt.Printf("%d ", int((solution & int64(2 ^ i)) / int64(2 ^ i)))
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Solution not found. \n")
	}
}
