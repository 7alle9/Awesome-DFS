package main

import (
	"fmt"
	"maps"
)

func test(m map[string]int) {
	fmt.Println(&m)
}

func main() {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	fmt.Println(maps.Keys(m))
}
