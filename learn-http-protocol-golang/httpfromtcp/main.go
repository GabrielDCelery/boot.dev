package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fileName := "messages.txt"

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed to open file %s, reason: %v\n", fileName, err)
	}

	defer file.Close()

	currentLine := ""
	buffer := make([]byte, 8)

	for {
		n, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			log.Fatalf("failed to read chunk into buffer, reason: %v\n", err)
		}

		// NOTE: even if we encounter an EOF error according to the docs we might have successfully read bytes into the buffer so first we need to handle the lefover chunk
		if n > 0 {
			// NOTE: We need to do the slicing because according to the docs the rest of the buffer might be used as scratch space https://pkg.go.dev/io@go1.25.4#Reader
			chunk := string(buffer[:n])
			parts := strings.Split(chunk, "\n")

			for i, part := range parts {
				currentLine += part

				// if we are looking at the last part and there are still leftover bytes in the file then we continue reading from the file
				if i == len(parts)-1 && err != io.EOF {
					break
				}

				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			}

			clear(buffer)
		}

		if err == io.EOF {
			break
		}
	}
}
