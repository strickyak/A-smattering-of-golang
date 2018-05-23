// +build main

/*
Serve (by HTTP) a single flat directory "." as a toplevel web site.

Usage:  cd /my/web/dir && go run .../flatweb.go --bind="localhost:8080"
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var Bind = flag.String("bind", "localhost:8080", "hostname:port to bind webserver to")

func Serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[0] != '/' {
		log.Fatalf("What sort of URL.Path doesn't begin with slash: %q", r.URL.Path)
	}
	path := r.URL.Path[1:]
	if path == "" {
		path = "index.html"
	}
	if strings.Contains(path, "/") {
		http.Error(w, "403 Forbidden (flat names only)", 403)
		return
	}
	if path[0] == '.' {
		http.Error(w, "403 Forbidden (no dot files)", 403)
		return
	}
	http.ServeFile(w, r, path)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", Serve)
	err := http.ListenAndServe(*Bind, nil)
	log.Fatalf("Cannot ListenAndServe: %v: %q", err, *Bind)
}
