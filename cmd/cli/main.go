package main

import (
	"flag"
	"fmt"
	"os"
)

// spidernet <url> --depth=<int>

func main() {
	depth := flag.Int("depth", 1, "depth of the crawler")

	flag.Parse()

	if depth == nil {
		fmt.Println("ex: spidernet <url> --depth=<int>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := flag.Arg(0)

	fmt.Println(url)
	fmt.Println(*depth)

}
