// +build main

package main

import (
	"bytes"
	//"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var Listen = flag.String("listen", ":18512", "listen on `:port`")
var Delay = flag.Duration("delay", 1.0*time.Second, "seconds to delay echo")
var Audio = flag.String("audio", "", "audio filename to capture")

func sum(bb []byte) uint64 {
	var z uint64
	for i, b := range bb {
		z += uint64(i) * uint64(b)
	}
	return z
}

func main() {
	flag.Parse()

	listenAddy, err := net.ResolveUDPAddr("udp4", *Listen)
	if err != nil {
		log.Fatalf("cannot resolve udp4 listen address: %q", *Listen)
	}

	conn, err := net.ListenUDP("udp4", listenAddy)
	if err != nil {
		log.Fatalf("cannot open socket: %v", err)
	}

	var audio *os.File
	if *Audio != "" {
		audio, err = os.Create(*Audio)
		if err != nil {
			log.Fatalf("cannot create audio file %q: %v", *Audio, err)
		}
	}

	var prev int64
	var t int64
	var realprev int64
	var realt int64
	for {
		bb := make([]byte, 3*512)
		sz, addy, err := conn.ReadFromUDP(bb)
		if err != nil {
			log.Fatalf("cannot ReadFromUDP: %v", err)
		}
		if len(bb) > 10 {
			t = int64(uint64(bb[2]) + (uint64(bb[3]) << 8) + (uint64(bb[4]) << 16) + (uint64(bb[5]) << 24) + (uint64(bb[6]) << 32) + (uint64(bb[7]) << 40))
		}

		realt = time.Now().UnixNano()
		log.Printf("Got %d bytes from %v .... [[[ %d ]]] %d", sz, addy, t-prev, realt-realprev)
		// HexDump(bb[:sz])
		prev = t
		realprev = realt

		if sz > 24 && bb[0] == 22 && bb[1] == 202 {
			go DelayAndEcho(conn, addy, bb[:sz])
		} else {
			log.Printf("Invalid Packet: %d, %d, %d", sz, bb[0], bb[1])
		}

		if audio != nil {
			_, err := audio.Write(bb[24:])
			if err != nil {
				log.Fatalf("cannot write audio: %v", err)
			}
		}
	}
}

var prevsend int64

func DelayAndEcho(conn *net.UDPConn, destAddy *net.UDPAddr, bb []byte) {
	time.Sleep(*Delay)

	sz, err := conn.WriteToUDP(bb, destAddy)
	t := time.Now().UnixNano()
	if err != nil {
		log.Printf("cannot WriteToUDP to %v: %v", destAddy, err)
	} else {
		log.Printf("(Echo %d bytes of %d to %q) %d", sz, len(bb), destAddy.String(), t-prevsend)
	}
	prevsend = t
}

func HexDump(bb []byte) {
	w := bytes.NewBuffer(nil)

	for i := 0; i < len(bb); i += 16 {
		fmt.Fprintf(w, "%4d: ", i)
		for j := 0; j < 16; j++ {
			if i+j < len(bb) {
				c := bb[i+j]
				fmt.Fprintf(w, " %02x", c)
			} else {
				fmt.Fprintf(w, "   ")
			}
			if j&3 == 3 {
				fmt.Fprintf(w, " ")
			}
		}
		fmt.Fprintf(w, "  ")
		for j := 0; j < 16 && i+j < len(bb); j++ {
			c := bb[i+j]
			if ' ' <= c && c <= '~' {
				fmt.Fprintf(w, "%c", c)
			} else if c >= 128 {
				fmt.Fprintf(w, "^")
			} else {
				fmt.Fprintf(w, ".")
			}
		}
		fmt.Fprintf(w, "\n")
	}
	log.Println(w.String())
}
