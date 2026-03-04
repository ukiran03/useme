package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Error: No input file")
	}
	for _, arg := range args {
		mount, err := GetMountPoint(arg)
		if err != nil {
			log.Printf("%s: %v\n", arg, err)
			continue
		}
		fmt.Printf("%s: %s\n", arg, mount)
	}
}
