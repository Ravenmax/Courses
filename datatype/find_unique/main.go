package main

import "fmt"

func main() {
	var array = [...]int8{2, 3, 4, 5, 1, 2, 3}
	var number int8
	for _, element := range array {
		fmt.Printf("%08b ^ %08b = %08b \n", number, element, number^element)
		number ^= element
	}
	fmt.Printf("result %08b\n", number)
}
