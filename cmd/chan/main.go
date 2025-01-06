package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Printf("time.Now(): %v\n", time.Now())
	<-time.After(2 * time.Second)
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("duration.String(): %v\n", duration.String())
	fmt.Printf("time.Now(): %v\n", time.Now())
}
