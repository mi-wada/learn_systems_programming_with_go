package main

import "fmt"

func a() *string {
	s := "hogehoge"
	return &s
}

func main() {
	fmt.Printf("a(): %v\n", a())
}
