package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Must specify s or c")
	}
	switch os.Args[1] {
	case "s":
		server()
	case "c":
		client()
	default:
		panic(fmt.Sprintf("Unknown flag: %s", os.Args[1]))
	}
}

const addr = "localhost:9999"

func server() {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Printf("Listening on %s\n", addr)
	for {
		buf := make([]byte, 1024)
		_, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		_, err = conn.WriteTo([]byte("Hello from server!"), remoteAddr)
		if err != nil {
			panic(err)
		}

		log.Println(string(buf))
	}
}

func client() {
	conn, err := net.Dial("udp4", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("Hello from client!"))

	buf := make([]byte, 1024)
	conn.Read(buf)
	log.Println(string(buf))
}
