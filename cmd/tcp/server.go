package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	fmt.Println("Starting server on port " + PORT + "...")
	listen, err := net.Listen(TYPE, ":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// log data
	fmt.Println(string(buffer[:]))

	// close conn
	conn.Close()
}
