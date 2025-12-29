package main

import (
	"log"

	"github.com/vanshitkumar/sappend/cmd"
)

func main() {
	if err := cmd.RootCommand(); err != nil {
		log.Fatal(err)
	}
}
