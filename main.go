package main

import "fmt"

func nextNode(chain *[]string) string {
	if len(*chain) == 0 {
		return ""
	}
	res := (*chain)[0]
	*chain = (*chain)[1:]
	return res
}

func main() {
	chain := []string{"a", "b", "c"}
	next := nextNode(&chain)
	fmt.Println(next)
	fmt.Printf("chain: %v\n", chain)
}
