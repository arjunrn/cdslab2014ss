package main

import "fmt"

func gcd(u , v uint64) uint64 {
	if v == 0 {
		return u
	}
	return gcd(v, u%v)
}


func friendlyNumbers(start, end uint64) {
	last := end - start + 1;

	theNum := make([]uint64, last)
	num := make([]uint64, last)
	den := make([]uint64, last)

	var ii, sum, done, n, factor uint64

	for i := start ; i <= end ; i++ {
		ii = i-start
		sum = 1+i
		theNum[ii] = i
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
		num[ii] = sum
		den[ii] = i
		n = gcd(num[ii], den[ii])
		num[ii]/=n
		den[ii]/=n
	}

	var i, j uint64
	for i = 0 ; i < last ; i++ {
		for j = i+1 ; j < last ; j++ {
			if (num[i] == num[j]) && (den[i] == den[j]) {
				fmt.Printf("%d and %d are FRIENDLY\n", int(theNum[i]), int(theNum[j]))
			}
		}
	}

}

func main() {
	var start, end uint64;

	for {
		fmt.Scanf("%d %d", &start, &end)
		if start == 0 && end == 0 {
			break
		}
		fmt.Printf("Numbers %d to %d\n", start, end)
		friendlyNumbers(start, end)
	}

}
