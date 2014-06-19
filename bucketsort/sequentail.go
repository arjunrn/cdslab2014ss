package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"container/list"
	"sort"
)

const(
	BucketNum = (0x7E - 0x21 + 1)
)

type Bucket struct{
	keys list.List
	bucketArr []string
}

func (b *Bucket) getLength() int {
	return b.keys.Len()
}

func (b *Bucket) sort()  {
	b.bucketArr = make([]string, b.keys.Len())
	i := 0
	for e := b.keys.Front() ; e != nil ; e = e.Next() {
		b.bucketArr[i] = e.Value.(string)
		i++
	}
	sort.Strings(b.bucketArr)
}

func bucketsort(inputString []string, writer *bufio.Writer) {
	buckets := make([]Bucket, BucketNum)

	for i := range inputString {
		charString := inputString[i]
		bucketIndex := uint8(charString[0]) - 0x21
		if bucketIndex >= BucketNum {
			panic(fmt.Sprintf("Index out of bucket slice range: %d for index: %d", bucketIndex, i))
		}
		buckets[bucketIndex].keys.PushFront(charString);

	}

	for i := range buckets {
		fmt.Println("Bucket %d size: %d", i, buckets[i].getLength())
		buckets[i].sort()
	}

	for i := range buckets {
		for j := range buckets[i].bucketArr{
			writer.WriteString(buckets[i].bucketArr[j])
		}
	}
}

func main(){
	fi, err := os.Open("bucketsort.in")
	if err != nil { panic(err) }

	defer func(){
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
    writer := bufio.NewWriter(fo)

	reader := bufio.NewReader(fi)
	list := list.New()
	counter := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF { panic(err) }
		if line == "" { break }
		counter++
		list.PushFront(line)
	}

	fmt.Println("Number of Lines: ", counter)

	number_slice := make([]string, list.Len())
	fmt.Println("Length of Number Slice: ", len(number_slice))

	i := 0
	for e := list.Front() ; e != nil ;  e = e.Next() {
		//fmt.Println(i," : ", e.Value.(string))
		number_slice[i] = e.Value.(string)
		i++
	}
	bucketsort(number_slice, writer)
}