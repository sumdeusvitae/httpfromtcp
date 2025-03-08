package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		// Handle error
		fmt.Println(err)
		return
	}
	defer ln.Close() // Close listener when done
	for {
		conn, err := ln.Accept()
		if err != nil {
			// Handle error
			continue
		}
		fmt.Println("Connection accepted")
		// Handle connection in a new goroutine
		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Read and write data to the connection
	lines := getLinesChannel(conn)
	for line := range lines {
		fmt.Println(line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		var line string
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err == io.EOF {
				if len(line) > 0 {
					ch <- line
				}
				fmt.Println("Connection closed")
				close(ch)
				return
			}
			if i := bytes.IndexByte(b[:n], '\n'); i >= 0 {
				line += string(b[:i])
				ch <- line
				line = ""
				line += string(b[i+1 : n])
			} else {
				line += string(b[:n])
			}
		}
	}()
	return ch
}
