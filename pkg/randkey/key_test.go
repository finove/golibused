package randkey

import (
	"fmt"
	"testing"
)

func TestKey(t *testing.T) {
	testRand(6)
	testRand(8)
	testRand(10)
	testRand(12)
	fmt.Printf("genpass:%s\n", GenPassword([]string{"1234567890", "abcdefghijklmnopqkl", "!@#$%^&*()"}, []int{3, 5, 7}))
	fmt.Printf("GeneratePassword:%s\n", GeneratePassword(9, 3, 2))
}

func testRand(count int) {
	fmt.Printf("count:%d key:%s\n", count, NumbersOnly(count))
	fmt.Printf("count:%d key:%s\n", count, NumberUpper(count))
	fmt.Printf("count:%d key:%s\n", count, NumberLower(count))
	fmt.Printf("count:%d key:%s\n", count, NumberPass(count))
}
