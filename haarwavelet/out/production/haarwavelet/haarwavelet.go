package main

import (
	"os"
	"bufio"
)



func main() {
	inputFileDes, inputErr := os.Open("bucketsort.in")

	if inputErr != nil {
		panic(inputErr)
	}

	defer func() {
		
		if , inputCloseErr != nil {
			panic(inputCloseErr)
		}
	}();

	inputReader := bufio.NewReader(inputFileDes)

	sizeByte, sizeErr := inputReader.ReadByte()
	if sizeErr != nil {
		panic(sizeErr)
	}

}
