package main

import (
	"fmt"
	"time"
	"runtime"
	"sync"
)

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


	doneChannel :=make(chan FriendlyResult)

	for i := start ; i <= end ; i++ {
		wg.Add(1)
		go func(i uint64,wg *sync.WaitGroup, doneChan chan FriendlyResult) {
			var sum, done, factor uint64

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
			res := FriendlyResult{}
			res.number = i; res.numerator = numerator; res.denominator=denominator
			doneChan <- res
			wg.Done()

		}(i,&wg,doneChannel)
	}

	go func(){
		wg.Wait()
		close(doneChannel)
//		fmt.Printf("Closed channel\n")
	}()

	recCount := 0
	for res := range doneChannel{
		i:=res.number-start
		theNum[i]=res.number
		num[i]=res.numerator
		den[i]=res.denominator
		recCount++
	}
	fmt.Printf("Received Count: %d\n",recCount)
//	fmt.Printf("Finished waiting for go routines.\n")

	var i, j uint64

	for i = 0 ; i < last ; i++ {
		for j = i+1 ; j < last ; j++ {
			if (num[i] == num[j]) && (den[i] == den[j]) {
				fmt.Printf("%d and %d are FRIENDLY\n", theNum[i], theNum[j])
//				fmt.Printf("%d/%d and %d/%d\n",num[i],den[i],num[j],den[j])
			}
		}
	}

}

func main() {
	runtime.GOMAXPROCS(1)
	var start, end uint64;

	for {
		fmt.Scanf("%d %d", &start, &end)
		if start == 0 && end == 0 {
			break
		}
		fmt.Printf("Numbers %d to %d\n", start, end)
		startTime := time.Now()
		friendlyNumbers(start, end)
		fmt.Printf("Friendly Numbers Compution took: %s\n", time.Since(startTime))
	}

}
