package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("time.Now(): %v\n", time.Now())
	<-time.After(10 * time.Second)
	fmt.Printf("time.Now(): %v\n", time.Now())
}
