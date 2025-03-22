package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sumdeusvitae/httpfromtcp/internal/request"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
	}

}
