package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"slices"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("Must specify s/c. c: serve, c: client.")
	}
	switch os.Args[1] {
	case "s":
		serve()
	case "c":
		client()
	default:
		panic(fmt.Sprintf("Unknown flag: %s", os.Args[1]))
	}
}

const (
	port = 8080
)

func serve() {
	addr := fmt.Sprintf("localhost:%d", port)
	server, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on %s: %v", addr, err))
	}
	defer server.Close()
	log.Printf("Started listening on %s\n", addr)

	ctx := context.Background()

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(fmt.Errorf("Unexpected error on net.Listener.Accept: %w", err))
		}

		go handler(ctx, conn)
	}
}

func handler(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Accept %v\n", conn.RemoteAddr())

	for {
		fmt.Println("Reading...")
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		req, err := http.ReadRequest(bufio.NewReader(conn))
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			log.Println("Timeout Keep-Alive, close TCP connection")
			break
		}
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			panic(err)
		}

		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			panic(err)
		}
		log.Println(string(dump))

		resp := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		}
		if slices.Contains(strings.Split(req.Header.Get("Accept-Encoding"), ","), "gzip") {
			content := "Hello, world(gzipped)\n"
			var buffer bytes.Buffer
			writer := gzip.NewWriter(&buffer)
			io.WriteString(writer, content)
			writer.Close()

			resp.Header.Set("Content-Encoding", "gzip")
			resp.Body = io.NopCloser(&buffer)
			resp.ContentLength = int64(buffer.Len())
		} else {
			content := "Hello, world\n"
			resp.Body = io.NopCloser(strings.NewReader(content))
			resp.ContentLength = int64(len(content))
		}
		resp.Write(conn)
	}
}

func client() {
	msgs := []string{
		"HELLO",
		"PRETTY",
		"BOY",
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for i := 0; i < len(msgs); i++ {
		req, err := http.NewRequest("POST", "http://localhost:8080", strings.NewReader(msgs[i]))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Accept-Encoding", "gzip")
		req.Write(conn)
		resp, err := http.ReadResponse(
			bufio.NewReader(conn),
			req,
		)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		reader := resp.Body
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(reader)
			if err != nil {
				panic(err)
			}
		}

		body, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}
