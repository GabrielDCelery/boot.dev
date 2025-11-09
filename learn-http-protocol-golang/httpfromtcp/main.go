package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fileName := "messages.txt"

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed to open file %s, reason: %v\n", fileName, err)
	}

	defer file.Close()

	buffer := make([]byte, 8)

	for {
		n, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			log.Fatalf("failed to read chunk into buffer, reason: %v\n", err)
		}

		// even if we encounter an EOF error according to the docs we might have successfully read bytes into the buffer
		// so first we need to handle the lefover chunk
		if n > 0 {
			// We need to do the slicing because according to the docs the rest of the buffer might be used as scratch space https://pkg.go.dev/io@go1.25.4#Reader
			fmt.Printf("read: %s\n", string(buffer[:n]))
		}

		if err == io.EOF {
			break
		}

		clear(buffer)
	}
}
