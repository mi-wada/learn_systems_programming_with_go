package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	fileName = "test.os.txt"
)

func create() {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.WriteString(file, "Hello, World!\n")
}

func append() {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.WriteString(file, "Hello, World!\n")
}

func cat() {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(os.Stdout, file)
}

func remove() {
	err := os.Remove(fileName)
	if err != nil {
		panic(err)
	}
}

func benchmarkSync() {
	f, _ := os.Create("test.bench.txt")
	a := time.Now()
	f.Write([]byte("Hello, World!\n"))
	b := time.Now()
	f.Sync()
	c := time.Now()
	f.Close()
	d := time.Now()
	fmt.Println("Write:", b.Sub(a))
	fmt.Println("Sync:", c.Sub(b))
	fmt.Println("Close:", d.Sub(c))
}

const dirName = "test.dir"

func mkdir() {
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		panic(err)
	}
}

func stat(fileName string) {
	file, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist\n", fileName)
			return
		}

		panic(err)
	}
	if file.IsDir() {
		fmt.Printf("%s is dir\n", fileName)
	} else {
		fmt.Printf("%s is not dir\n", fileName)
	}
}

func rmdir() {
	err := os.Remove(dirName)
	if err != nil {
		panic(err)
	}
}

func main() {
	create()
	fmt.Println("After create:")
	cat()
	append()
	fmt.Println("After append:")
	cat()
	remove()
	stat(fileName)

	benchmarkSync()

	mkdir()
	stat(dirName)
	rmdir()
	stat(dirName)

	f, err := os.Create("test.truncate.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data := make([]byte, 101)
	f.Write(data)
	f.Truncate(100)
	fStat, err := os.Stat("README.md")
	if err != nil {
		panic(err)
	}
	fmt.Println(fStat.Size(), fStat.ModTime(), fStat.Name(), fStat.Size())

}
