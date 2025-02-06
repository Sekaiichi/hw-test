package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	fmt.Println(reverseText("Hello, OTUS!"))
}

func reverseText(text string) string {
	return reverse.String(text)
}
