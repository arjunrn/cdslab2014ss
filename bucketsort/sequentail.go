package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"container/list"
)

const(
	BucketNum = 100
)

type Bucket struct{
	keys []string;
}




func bucketsort(inputString []string) {
	for i := range inputString {
		charString := inputString[i]
		var stringValue uint64 = 0
		for j := 0 ; j < len(charString) ; j++ {
			line := uint64(charString[j])
			stringValue += (line << ( 8 * uint(j) ))
		}
	}
	//buckets []Bucket := make([]Bucket, BucketNum)
	for i := 0 ; i < BucketNum ; i++{
		
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
	for e := list.Front() ; i < len(number_slice) ;  e.Next() {
		number_slice[i] = e.Value.(string)
		i++
	}
	bucketsort(number_slice)
}