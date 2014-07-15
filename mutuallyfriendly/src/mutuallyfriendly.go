package main

import (
	"fmt"
	"time"
	"runtime"
	"sync"
	"os"
)

var workerCount int

func gcd(u , v uint64) uint64 {
	if v == 0 {
		return u
	}
	return gcd(v, u%v)
}

type FriendlyResult struct{
	number, numerator, denominator uint64
}

func friendlyNumbers(start, end uint64) {
	var wg sync.WaitGroup
	last := end - start + 1;

	theNum := make([]uint64, last)
	num := make([]uint64, last)
	den := make([]uint64, last)

	coreAssignmentCount := uint64(last / uint64(workerCount))
	if coreAssignmentCount == 0 {
		coreAssignmentCount = 1
	}
	fmt.Printf("Core Assignment: %d\n", coreAssignmentCount)

	for n := start ; n <= end ; n+=coreAssignmentCount {
		wg.Add(1)
		go func(startInt uint64) {
			var sum, done, factor uint64
			var endInt uint64
			if startInt+coreAssignmentCount > end {
				endInt = end
			} else {
				endInt = startInt+coreAssignmentCount
			}
			for i := startInt; i < endInt; i++ {
				sum = 1+i
				done = i
				factor = 2

				for factor < done {
					if (i%factor) == 0 {
						sum += (factor+(i/factor))
						done = i/factor
						if done == factor {
							sum -= factor
						}
					}
					factor++
				}

				numerator := sum
				denominator := i
				n := gcd(numerator, denominator)
				numerator/=n
				denominator/=n
				solIndex := i - start
				num[solIndex] = numerator
				den[solIndex] = denominator
				theNum[solIndex] = i
			}
			wg.Done()
		}(n)
	}

	wg.Wait()

	var printWg sync.WaitGroup
	var i uint64
	printWg.Add(int(last))
	for i = 0 ; i < last ; i++ {
		go func(row uint64) {
			defer printWg.Done()
			for j := row+1 ; j < last ; j++ {
				if (num[row] == num[j]) && (den[row] == den[j]) {
					fmt.Printf("%d and %d are FRIENDLY\n", theNum[row], theNum[j])
				}
			}
		}(i)
	}
	printWg.Wait()
}

func main() {
	logHandler , logErr := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logErr != nil {
		panic(logErr)
	}
	defer logHandler.Close()

	workerCount = runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount)
	var start, end uint64;

	for {
		fmt.Scanf("%d %d", &start, &end)
		if start == 0 && end == 0 {
			break
		}
		fmt.Printf("Numbers %d to %d\n", start, end)

		startTime := time.Now()
		friendlyNumbers(start, end)
		fmt.Fprintf(logHandler, "Time to compute: %s\n", time.Since(startTime))
	}

}
