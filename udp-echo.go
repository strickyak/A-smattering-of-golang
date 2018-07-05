// +build main

package main

import (
	//"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var Listen = flag.String("listen", ":18512", "listen on `:port`")
var Delay = flag.Duration("delay", 1.0*time.Second, "seconds to delay echo")

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

	for {
		bb := make([]byte, 3*512)
		sz, addy, err := conn.ReadFromUDP(bb)
		if err != nil {
			log.Fatalf("cannot ReadFromUDP: %v", err)
		}
		log.Printf("Got %d bytes from %v", sz, addy)
		HexDump(bb[:sz])

		if sz > 24 && bb[0] == 22 && bb[1] == 202 {
			go DelayAndEcho(conn, addy, bb[:sz])
		} else {
			log.Printf("Invalid Packet: %d, %d, %d", sz, bb[0], bb[1])
		}
	}
}

func DelayAndEcho(conn *net.UDPConn, destAddy *net.UDPAddr, bb []byte) {
	time.Sleep(*Delay)

	sz, err := conn.WriteToUDP(bb, destAddy)
	if err != nil {
		log.Printf("cannot WriteToUDP to %v: %v", destAddy, err)
	} else {
		log.Printf("(Echo %d bytes of %d to %q)", sz, len(bb), destAddy.String())
	}
}

func HexDump(bb []byte) {
	for i := 0; i < len(bb); i += 16 {
		fmt.Printf("%4d: ", i)
		for j := 0; j < 16; j++ {
			if i+j < len(bb) {
				c := bb[i+j]
				fmt.Printf(" %02x", c)
			} else {
				fmt.Printf("   ")
			}
			if j&3 == 3 {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("  ")
		for j := 0; j < 16 && i+j < len(bb); j++ {
			c := bb[i+j]
			if ' ' <= c && c <= '~' {
				fmt.Printf("%c", c)
			} else if c >= 128 {
				fmt.Printf("^")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}
