package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

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
			defer conn.Close()

			fmt.Printf("connection has been accepted on %s\n", conn.LocalAddr())

			req, err := request.RequestFromReader(conn)

			if err != nil {
				log.Fatalf("failed to parse request, reason: %v", err)
			}

			fmt.Println("Request line:")
			fmt.Printf("- Methd: %s\n", req.RequestLine.Method)
			fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
			fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

			fmt.Println("Headers")
			for k, v := range req.Headers {
				fmt.Printf("- %s - %v\n", k, v)
			}
			response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 2\r\n\r\nOK"
			conn.Write([]byte(response))

			fmt.Printf("connection has been closed on %s\n", conn.LocalAddr())
		}(conn)
	}
}
