package main

import (
	"fmt"
	"sort"
)

type Item struct{
	value   int
	weight  int
	density float64
}

type ItemV []Item

func (slice ItemV) Len() int {
	return len(slice)
}

func (slice ItemV) Less(i, j int) bool {
	return slice[i].density < slice[j].density;
}

func (slice ItemV) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func knapsack_f(n int, M int, items []Item) int {
	var v, w, r int
	v = 0
	w = 0
	r = 0

	if n < 1 {
		return 0
	}

	for M-w >= 0 {
		r = greater_f(r, v+knapsack_f(n-1, M-w, items[1:]))
		v += items[0].value
		w += items[0].weight
	}

	return r
}

func greater_f(x, y int) int {
	if (x > y) {
		return x
	}else {
		return y
	}
}



func main() {
	var n, M int

	fmt.Scanf("%d %d", &n, &M)

	items := make([]Item, n)

	for i := 0 ; i < n; i++ {
		fmt.Scanf("%d %d", &items[i].value, &items[i].weight)
		items[i].density = float64(items[i].value / items[i].weight)
	}

	sort.Sort(ItemV(items))

	fmt.Printf("%d\n", knapsack_f(n, M, items))

}
