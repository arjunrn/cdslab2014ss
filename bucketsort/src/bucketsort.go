package main

import (
	"fmt"
	"os"
	"bufio"
	"container/list"
	"sort"
	"strings"
	"bytes"
	"runtime"
	"sync"
	"time"
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
}

func bucketSort(buckets []Bucket) {
	var wg sync.WaitGroup
	wg.Add(len(buckets))
	for i := range buckets {
		go func(i int) {
			buckets[i].sort()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logHandler , logErr := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logErr != nil {
		panic(logErr)
	}
	defer logHandler.Close()


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

	startSortTime := time.Now()
	bucketSort(buckets)
	_, writeErr := fmt.Fprintf(logHandler, "Time for bucket sort: %s\n", time.Since(startSortTime))
	if writeErr != nil {
		panic(writeErr)
	}


	writer := bufio.NewWriter(fo)
	for i := range buckets {
		stringBuilder := bytes.NewBufferString("")
		for j := range buckets[i].bucketArr {
			stringBuilder.WriteString(fmt.Sprintf("%s\n", buckets[i].bucketArr[j]))
		}
		writer.WriteString(stringBuilder.String())
		writer.Flush()
	}
}
