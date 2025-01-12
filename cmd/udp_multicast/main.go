package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
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

const (
	addr     = "244.0.0.100:6000"
	interval = 1 * time.Second
)

func server() {
	log.Printf("Start tick server %s", addr)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	start := time.Now()
	wait := start.Truncate(interval).Add(interval).Sub(start)
	time.Sleep(wait)
	ticker := time.Tick(interval)
	for now := range ticker {
		conn.Write([]byte(now.Format(time.RFC3339Nano)))
		log.Println("Tick: ", now.Format(time.RFC3339Nano))
	}
}

func client() {
	log.Printf("Listen tick server at %s", addr)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	listerner, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		panic(err)
	}
	defer listerner.Close()

	buf := make([]byte, 1024)
	for {
		length, remoteAddr, err := listerner.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		log.Printf("Received: %s from %s\n", string(buf[:length]), remoteAddr)
	}
}
