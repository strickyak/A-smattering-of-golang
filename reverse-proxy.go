// +build main

/*
Listen on given HTTP and/or HTTPS ports and reverse-proxy all connections to given target.

Example Usage:

	$ sh make-tls-key.sh
	$ go run flatweb.go --bind=localhost:7000 &  ## The example target webserver.
	$ go build reverse-proxy.go
	$ sudo ./reverse-proxy --target=http://localhost:7000/ &  ## The reverse proxy
	$ wget --no-check-certificate -v -O /dev/stdout  https://localhost/flatweb.go

See source for other flags.
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var target = flag.String("target", "http://localhost:8888/", "Target to reverse proxy connections to")
var bind = flag.String("bind", "localhost:80", "hostname:port to bind webserver to; empty to not listen")
var bindTLS = flag.String("tls_bind", "localhost:443", "hostname:port to bind TLS webserver to; empty to not listen")
var certFileTLS = flag.String("tls_certfile", ".tls.crt", "Certificate file for TLS")
var keyFileTLS = flag.String("tls_keyfile", ".tls.key", "Key file for TLS")

func main() {
	flag.Parse()

	u, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("Cannot parse target URL: %v: %q", err, *target)
	}
	rp := httputil.NewSingleHostReverseProxy(u)

	var useful bool
	if *bind != "" {
		useful = true
		go func() {
			log.Printf("Plain Listening on %q", *bind)
			err := http.ListenAndServe(*bind, rp)
			log.Fatalf("Cannot ListenAndServe: %v: %q", err, *bind)
		}()
	}
	if *bindTLS != "" {
		useful = true
		go func() {
			log.Printf("TLS Listening on %q", *bindTLS)
			err := http.ListenAndServeTLS(*bindTLS, *certFileTLS, *keyFileTLS, rp)
			log.Fatalf("Cannot ListenAndServeTLS: %v: %q", err, *bindTLS)
		}()
	}
	if !useful {
		log.Fatal("Nothing useful is being done.")
	}
	time.Sleep(10 * 366 * 24 * time.Hour) // Sleep over 10 years.
}
