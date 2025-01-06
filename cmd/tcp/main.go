package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("must pass l or d")
	}
	if os.Args[1] == "l" {
		listen()
	} else if os.Args[1] == "d" {
		dial()
	}
}

func listen() {
	server, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	for {
		fmt.Println("listen on localhost:8080")
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		conn.Write([]byte("hello"))
		reader := bufio.NewReader(conn)
		res, _ := reader.ReadString('\n')
		if string(res) == "hello\n" {
			conn.Write([]byte("bye"))
			conn.Close()
		} else {
			conn.Write([]byte("invalid"))
			conn.Close()
		}
	}
}

func dial() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("hello\n"))
	res, err := io.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("string(res): %v\n", string(res))
}
