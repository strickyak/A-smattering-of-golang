// +build main

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("Usage:   go run ls.go dir1 dir2...")
	}

	for _, arg := range flag.Args() {
		fmt.Printf("%s:\n", arg)
		f, err := os.Open(arg)
		if err != nil {
			log.Fatalf("cannot Open %q: %v", arg, err)
		}

		names, err := f.Readdirnames( /*all*/ -1)
		if err != nil {
			log.Fatalf("cannot Readdirnames %q: %v", arg, err)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Printf("{%s}\n", name)
		}
	}
}
