package main

import (
	"fmt"
	"unsafe"
)

func IsLittle() bool {
	var number = 0x0001
	pointer := (*int8)(unsafe.Pointer(&number))
	return *pointer == 1
}
func IsBig() bool {
	return !IsLittle()
}
func main() {
	if IsLittle() {
		fmt.Println("is little")
	} else {
		fmt.Println("is big")
	}
}
