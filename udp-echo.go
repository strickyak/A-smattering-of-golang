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

	var prev int64
	var t int64
	for {
		bb := make([]byte, 3*512)
		sz, addy, err := conn.ReadFromUDP(bb)
		if err != nil {
			log.Fatalf("cannot ReadFromUDP: %v", err)
		}
		if len(bb) > 10 {
			t = int64(bb[2]) + int64((bb[3]<<8)) + int64((bb[4]<<16)) + int64((bb[5]<<24)) + int64((bb[6]<<32)) + int64((bb[7]<<40))
		}

		log.Printf("Got %d bytes from %v .... [[[ %d ]]]", sz, addy, t - prev)
		prev = t
		//HexDump(bb[:sz])

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
	os.Stdout.Write(w.Bytes())
}
