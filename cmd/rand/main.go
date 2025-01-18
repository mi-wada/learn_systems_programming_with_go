package main

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("By math/rand")
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 10; i++ {
		fmt.Println(r.Int31())
	}

	fmt.Println("By crypto/rand")
	for i := 0; i < 10; i++ {
		b := make([]byte, 4)
		cryptorand.Read(b)
		n := int(b[0]) | int(b[1])<<8 | int(b[2])<<16 | int(b[3])<<24
		fmt.Println(n)
	}
}
