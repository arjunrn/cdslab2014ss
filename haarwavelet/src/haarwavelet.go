package main

import (
	"os"
	"math"
	"encoding/binary"
	"fmt"
	"bytes"
	"sync"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(8)
	const SIZE_OF_INT32 = 4
	SquareRoot2 := math.Sqrt(2)
	fmt.Printf("Square Root: %f\n", SquareRoot2)
	inputFileDes, inputErr := os.Open("image.in")

	if inputErr != nil {
		panic(inputErr)
	}

	defer func() {
		inputCloseErr := inputFileDes.Close()
		if inputCloseErr != nil {
			panic(inputCloseErr)
		}
	}();

	outputFile, outErr := os.Create("image.out")
	if outErr != nil {
		panic(outErr)
	}
	defer func() {
		outputCloseErr := outputFile.Close();
		if outputCloseErr != nil {
			panic(outputCloseErr)
		}
	}();

	sizeByte := make([]byte, 8)

	readSize, sizeReadErr := inputFileDes.Read(sizeByte)
	fmt.Printf("Meta Read Size: %d\n", readSize)
	if sizeReadErr != nil {
		panic(sizeReadErr)
	}

	size := int(binary.LittleEndian.Uint64(sizeByte))
	fmt.Printf("Size of Image File: %d\n", size)

	outputFile.Write(sizeByte)

	byteBufferSize := (size * size) * SIZE_OF_INT32
	bytePixels := make([]byte, byteBufferSize)

	readBuffer := make([]byte, 1024 * 1024 * 1204)

	readCounter := 0
	inputFileDes.Seek(8, 0)
	fileReadCounter := time.Now()
	for n, e := inputFileDes.Read(readBuffer); e == nil; n, e = inputFileDes.Read(readBuffer) {
		copy(bytePixels[readCounter:readCounter+n], readBuffer[:n])
		readCounter+=n
	}
	if readCounter != byteBufferSize {
		fmt.Printf("Bytes Read: %d Byte Buffer Size: %d\n", readCounter, byteBufferSize)
		panic("Could not read the size")
	}


	pixels := make([]int32, (size*size))
	readBytesBuffer := bytes.NewBuffer(bytePixels)
	binary.Read(readBytesBuffer, binary.LittleEndian, pixels)
	fmt.Printf("Time to read file into buffer: %s\n", time.Since(fileReadCounter))
	transformCounter := time.Now()
	for s := size ; s > 1; s/=2 {
		mid := s / 2

		var xWaitGroup sync.WaitGroup
		xWaitGroup.Add(mid)
		for y := 0; y < mid; y++ {
			go func(y int) {
				defer xWaitGroup.Done()
				for x := 0; x < mid; x++ {
					a := pixels[(y * size) + x]
					d := a
					a = int32(float64((a + pixels[(y * size) + (mid + x)])) / SquareRoot2)
					d = int32(float64((d - pixels[(y * size) + (mid + x)])) / SquareRoot2)
					pixels[(y*size)+x] = a
					pixels[(y*size)+(mid+x)] = d
				}
			}(y)
		}
		xWaitGroup.Wait()

		var yWaitGroup sync.WaitGroup
		yWaitGroup.Add(mid)
		for y := 0; y < mid; y++ {
			go func(y int) {
				defer yWaitGroup.Done()
				for x := 0; x < mid; x++ {
					a := pixels[(y * size) + x]
					d := a
					a = int32(float64((a + pixels[((y + mid) * size) + (x)])) / SquareRoot2)
					d = int32(float64((d - pixels[((y + mid) * size) + (x)])) / SquareRoot2)
					pixels[(y*size)+x] = a
					pixels[((y+mid)*size)+(x)] = d
				}
			}(y)
		}
		yWaitGroup.Wait()
	}
	fmt.Printf("Time to transform: %s\n", time.Since(transformCounter))

	writeCounter := time.Now()
	byteBuffer := new(bytes.Buffer)
	fmt.Printf("Time to create write buffer: %s\n", time.Since(writeCounter))
	buffErr := binary.Write(byteBuffer, binary.LittleEndian, pixels)
	if buffErr != nil {
		panic(buffErr)
	}

	_, writeErr := outputFile.Write(byteBuffer.Bytes())


	if writeErr != nil {
		panic(writeErr)
	}
	fmt.Printf("Time to write to file: %s\n", time.Since(writeCounter))

	fmt.Printf("Total Time: %s\n", time.Since(fileReadCounter))
}
