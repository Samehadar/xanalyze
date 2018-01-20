package main

import (
	"fmt"

	"github.com/mvryan/fasttag"
)

func main() {
	fmt.Println("Hello, world")
	words := fasttag.WordsToSlice("Hello, world")
	fmt.Println("Words: ", words)
	pos := fasttag.BrillTagger(words)
	fmt.Println("Parts of Speech: ", pos)
}
