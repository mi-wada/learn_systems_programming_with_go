package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func readAll() {
	f, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

func read() {
	f, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := make([]byte, 1024)
	var res string
	for {
		_, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		res += string(buf)
	}
	fmt.Println(res)
}

func readTCP() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: example.com\r\n\r\n"))
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}

func readStdin() {
	res := make(chan string)
	go func() {
		var stdin string
		for {
			buf := make([]byte, 5)
			_, err := os.Stdin.Read(buf)
			stdin += string(buf)
			if err == io.EOF {
				fmt.Println("EOF")
				res <- stdin
				break
			}
		}
	}()
	for {
		select {
		case stdin := <-res:
			fmt.Println("stdin:", stdin)
			return
		case <-time.After(1 * time.Second):
			fmt.Println("waiting...")
		}
	}

}

func limitReader() {
	f, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	lReader := io.LimitReader(f, 10)
	bytes, err := io.ReadAll(lReader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func sectionReader() {
	f, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sReader := io.NewSectionReader(f, 5, 10)
	bytes, err := io.ReadAll(sReader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func endian() {
	// 32bit big endian 10000
	data := []byte{0x0, 0x0, 0x27, 0x10}
	var i int32
	binary.Read(bytes.NewReader(data), binary.BigEndian, &i)
	fmt.Println(i)
}

func copyOldToNew() {
	oldFile, err := os.Open("test.old.txt")
	if err != nil {
		panic(err)
	}
	defer oldFile.Close()

	newFile, err := os.Create("test.new.txt")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		panic(err)
	}
}

func genRand1024BinFile() {
	f, err := os.Create("test.rand.bin")
	if err != nil {
		panic(err)
	}

	io.CopyN(f, rand.Reader, 1024)
}

func main() {
	src := strings.NewReader("abcdef")

	buf := make([]byte, 1024)
	writer := bytes.NewBuffer(buf)
	copyN(writer, src, 3)
	fmt.Printf("writer.String(): %v\n", writer.String())
	_, err := copyN(writer, src, 4)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("writer.String(): %v\n", writer.String())
	_, err = copyN(writer, src, 1)
	fmt.Printf("err: %v\n", err)
}

func copyN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	written, err = io.Copy(dst, io.LimitReader(src, n))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		err = io.EOF
	}
	return
}
