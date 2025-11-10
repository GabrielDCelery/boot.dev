package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp4", "localhost:42069")

	if err != nil {
		log.Fatalf("failed to resolve udp address, reason: %v\n", err)
	}

	conn, err := net.DialUDP("udp4", nil, remoteUDPAddr)

	if err != nil {
		log.Fatalf("failed to establish connection to %d", remoteUDPAddr.Port)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")

		str, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Printf("stdin closed, exitting...\n")
				break
			}

			log.Fatalf("failed to read input")
		}

		_, err = conn.Write([]byte(str))

		if err != nil {
			log.Fatalf("failed to write input to connection")
		}
	}
}
