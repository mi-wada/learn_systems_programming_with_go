package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type MyString string

func (m *MyString) Write(p []byte) (n int, err error) {
	*m = MyString(p)
	return len(p), nil
}

func main() {
	s := "hello, world"
	sbr := fmt.Sprintf("%s\n", s)

	var ms MyString
	ms.Write([]byte(s))
	fmt.Printf("ms: %v\n", ms)

	os.Stdout.Write([]byte(sbr))
	fmt.Println()

	var bbuf bytes.Buffer
	for i := 0; i < 5; i++ {
		bbuf.Write([]byte(fmt.Sprintf("bbuf%d: %s", i, sbr)))
	}
	os.Stdout.Write(bbuf.Bytes())

	var sbuf strings.Builder
	for i := 0; i < 5; i++ {
		sbuf.Write([]byte(fmt.Sprintf("sbuf%d: %s", i, sbr)))
	}
	os.Stdout.Write([]byte(sbuf.String()))

	buffredStdout := bufio.NewWriter(os.Stdout)
	defer buffredStdout.Flush()
	for i := 0; i < 5; i++ {
		buffredStdout.Write([]byte(fmt.Sprintf("buffredStdout%d: %s", i, sbr)))
	}

	var multi1 MyString
	var multi2 MyString
	multiWriter := io.MultiWriter(&multi1, &multi2)
	multiWriter.Write([]byte(s))
	fmt.Printf("multi1: %v\n", multi1)
	fmt.Printf("multi2: %v\n", multi2)

	f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte(sbr))

	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: example.com\r\n\r\n"))
	io.Copy(os.Stdout, conn)

	csvFile, err := os.Create("test.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write([]string{"col1", "col2"})
	csvWriter.Write([]string{"hello", "world"})
	csvWriter.Write([]string{"goodbye", "world"})
	csvWriter.Flush()

	gzipFile, err := os.Create("test.txt.gz")
	if err != nil {
		panic(err)
	}
	defer gzipFile.Close()
	gzipWriter := gzip.NewWriter(gzipFile)
	gzipWriter.Header.Name = "test.txt"
	gzipWriter.Write([]byte(sbr))
	gzipWriter.Close()
}
