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

func server() {
	listener, err := net.Listen("unix", "test.socketfile")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte("Hello"))

		buf := make([]byte, 5)
		conn.Read(buf)
		if string(buf) == "World" {
			log.Println("World received")
		} else {
			log.Panicf("Unexpected message: %s", string(buf))
		}

		conn.Write([]byte("Bye"))
		conn.Close()
	}
}

func client() {
	conn, err := net.Dial("unix", "test.socketfile")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 5)
	conn.Read(buf)
	if string(buf) == "Hello" {
		log.Println("Hello received")
	} else {
		log.Panicf("Unexpected message: %s", string(buf))
	}

	conn.Write([]byte("World"))

	buf = make([]byte, 3)
	conn.Read(buf)
	if string(buf) == "Bye" {
		log.Println("Bye received")
	} else {
		log.Panicf("Unexpected message: %s", string(buf))
	}
}
