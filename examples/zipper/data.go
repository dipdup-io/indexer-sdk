package main

import "fmt"

// ZipData -
type ZipData struct {
	key   int
	value string
}

// Key -
func (z ZipData) Key() int {
	return z.key
}

// String -
func (z ZipData) String() string {
	return fmt.Sprintf("Key: %d | Value: %s", z.key, z.value)
}
