// +build main

// Create a filename containing every allowed character in its name.
// In UNIX filesystems, only two bytes are forbidden: '\000' and '/'.
package main

import (
	"log"
	"os"
)

func main() {
	var buf []byte
	for i := 1; i < 256; i++ {
		if i == '/' {
			continue
		}
		buf = append(buf, byte(i))
	}

	filename := string(buf)
	log.Printf("Creating len=%d  %q", len(filename), filename)
	fd, err := os.Create(filename)
	if err != nil {
		log.Printf("Create ERROR: %v", err)
	}
	err = fd.Close()
	if err != nil {
		log.Printf("Close ERROR: %v", err)
	}
}
