package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		currentLine := ""

		buffer := make([]byte, 8)

		for {
			n, err := f.Read(buffer)

			if err != nil && err != io.EOF {
				log.Fatalf("failed to read chunk into buffer, reason: %v\n", err)
			}

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

					ch <- currentLine
					currentLine = ""
				}

				clear(buffer)
			}

			if err == io.EOF {
				break
			}
		}

		close(ch)
	}()

	return ch
}

func main() {
	fileName := "messages.txt"

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed to open file %s, reason: %v\n", fileName, err)
	}

	defer file.Close()

	linesChan := getLinesChannel(file)

	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
	}
}
