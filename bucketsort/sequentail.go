package main

import (
	"fmt"
	"os"
	"bufio"
	"container/list"
	"sort"
	"strings"
	"time"
	"bytes"
	"runtime"
	"sync"
)

const (
	BucketNum      = (0x7E - 0x21 + 1)
	ReadBufferSize = 1024 * 1024 * 1204 * 8
)

type Bucket struct{
	keys        list.List
	bucketArr   []string
}

func (b *Bucket) getLength() int {
	return b.keys.Len()
}

func (b *Bucket) sort() {
	b.bucketArr = make([]string, b.keys.Len())
	i := 0
	for e := b.keys.Front() ; e != nil ; e = e.Next() {
		b.bucketArr[i] = e.Value.(string)
		i++
	}
	sort.Strings(b.bucketArr)
	b.doneChannel <- 1
}

func bucketSort(buckets []Bucket) {
	var wg sync.WaitGroup
	wg.Add(len(buckets))
	for i := range buckets {
		buckets[i].doneChannel = make(chan int)
		go func(i int, wg sync.WaitGroup) {
			defer wg.Done()
			buckets[i].sort()
		}(i, wg)
	}
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(8)
	fi, err := os.Open("bucketsort.in")
	if err != nil { panic(err) }

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	fo, err := os.Create("bucketsort.out")
	if err != nil { panic(err) }
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer


	buckets := make([]Bucket, BucketNum)
	rbuff := make([]byte, ReadBufferSize)

	start := time.Now()

	for n, e := fi.Read(rbuff) ; e == nil ; n, e = fi.Read(rbuff) {
		readString := strings.TrimSpace(string(rbuff[:n]))
		inputValues := strings.Split(readString, "\n")
		for i := range inputValues {
			charString := inputValues[i]
			bucketIndex := uint8(charString[0]) - 0x21
			if bucketIndex >= BucketNum {
				panic(fmt.Sprintf("Index out of bucket slice range: %d for index: %d", bucketIndex, i))
			}
			buckets[bucketIndex].keys.PushFront(charString);
		}
	}

	elapsedRead := time.Since(start)
	fmt.Printf("Reading and file and bucketing took %s\n", elapsedRead)

	start = time.Now()
	bucketSort(buckets)
	elapsedRead = time.Since(start)
	fmt.Printf("Bucket Sort took %s\n", elapsedRead)

	start = time.Now()
	writer := bufio.NewWriter(fo)
	for i := range buckets {
		stringBuilder := bytes.NewBufferString("")
		for j := range buckets[i].bucketArr {
			stringBuilder.WriteString(fmt.Sprintf("%s\n", buckets[i].bucketArr[j]))
		}
		writer.WriteString(stringBuilder.String())
		writer.Flush()
	}
	elapsedRead = time.Since(start)
	fmt.Printf("Writing to file took %s\n", elapsedRead)
}
