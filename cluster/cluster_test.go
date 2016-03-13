package main

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestCrc(t *testing.T) {

	b := make([]byte, 16)
	count := make(map[int]int)

	for i := 0; i < 9999; i++ {
		rand.Read(b)
		c := Crc16(b) % 4

		count[c] = count[c] + 1
	}

	fmt.Println(count)
}
