package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func readFromString() {
	src := `col1,col2
hello,world
goodbye,world
`
	reader := bufio.NewReader(strings.NewReader(src))

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Print(line)
	}
}

func readFromFile() {
	f, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Print(line)
	}
}

func readFromFileAsCSV() {
	f, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(line)
	}
}

func scan() {
	src := `col1,col2
hello,world
goodbye,world
`
	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func fscan() {
	src := "123 1.234 1.0e4 test"
	reader := strings.NewReader(src)
	var i int
	var j, k float64
	var l string
	_, err := fmt.Fscan(reader, &i, &j, &k, &l)
	if err != nil {
		panic(err)
	}
	fmt.Printf("i=%#v j=%#v k=%#v l=%#v\n", i, j, k, l)
}

func tee() {
	header := strings.NewReader("HEADER\n")
	content := strings.NewReader("Hello!!!\n")
	footer := strings.NewReader("FOOTER\n")

	reader := io.MultiReader(header, content, footer)
	teeReader := io.TeeReader(reader, os.Stdout)
	bytes, err := io.ReadAll(teeReader)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(string(bytes))
}

func pipe() {
	reader, writer := io.Pipe()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer writer.Close()
		writer.Write([]byte("hello"))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bytes, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	}()

	wg.Wait()
}

func main() {
	pipe()
}
