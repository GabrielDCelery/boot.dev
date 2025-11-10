package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	listener, err := net.Listen("tcp4", ":42069")

	if err != nil {
		log.Fatalf("failed to create listener, reason: %v\n", err)
	}

	fmt.Printf("started listener on %s\n", listener.Addr())

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalf("failed to accept connection, reason: %v\n", err)
		}

		go func(conn net.Conn) {
			fmt.Printf("connection has been accepted on %s\n", conn.LocalAddr())

			linesChan := getLinesChannel(conn)

			for line := range linesChan {
				fmt.Printf("%s\n", line)
			}

			conn.Close()

			fmt.Printf("connection has been closed on %s\n", conn.LocalAddr())
		}(conn)
	}
}
