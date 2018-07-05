// +build main

package main

import (
	//"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	//"time"
)

//var Dest = flag.String("dest", "", "send probe to `host:port`")
var Listen = flag.String("listen", ":2237", "listen on `:port`")

//var Delay = flag.Duration("delay", 1.0*time.Second, "seconds to delay")

func sum(bb []byte) uint64 {
	var z uint64
	for i, b := range bb {
		z += uint64(i) * uint64(b)
	}
	return z
}

/*
var history map[int]uint64

func probeLoop(conn *net.UDPConn, destAddy *net.UDPAddr) {
	n := 20
	for {
		time.Sleep(*Delay)
		bb := make([]byte, n)
		rand.Read(bb)

		sz, err := conn.WriteToUDP(bb, destAddy)
		if err != nil {
			log.Printf("cannot WriteToUDP to %v: %v", destAddy, err)
		} else {
			z := sum(bb)
			history[n] = z
			log.Printf("Sent %d bytes to %v = %d", sz, destAddy, z)
		}
		n++
		if n > 100 {
			n = 20
		}
	}
}
*/

func main() {
	flag.Parse()

	listenAddy, err := net.ResolveUDPAddr("udp4", *Listen)
	if err != nil {
		log.Fatalf("cannot resolve udp4 listen address: %q", *Listen)
	}
	/*
		destAddy, err := net.ResolveUDPAddr("udp4", *Dest)
		if err != nil {
			log.Fatalf("cannot resolve udp4 dest address: %q", *Dest)
		}
	*/
	conn, err := net.ListenUDP("udp4", listenAddy)
	if err != nil {
		log.Fatalf("cannot open socket: %v", err)
	}

	// go probeLoop(conn, destAddy)

	for {
		bb := make([]byte, 3*512)
		sz, addy, err := conn.ReadFromUDP(bb)
		if err != nil {
			log.Fatalf("cannot ReadFromUDP: %v", err)
		}
		// z := sum(bb)
		// zz := history[sz]
		// log.Printf("Got %d bytes from %v = %d %v", sz, addy, sum(bb), (z == zz))
		log.Printf("Got %d bytes from %v", sz, addy)
		HexDump(bb[:sz])
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
