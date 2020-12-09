package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	listen      = flag.String("listen", ":8080", "listen address")
	dir         = flag.String("dir", "htdocs", "directory to serve")
	openBrowser = flag.Bool("openBrowser", false, "open a browser while serving")
)

func main() {
	flag.Parse()
	log.Printf("listening on %s", *listen)

	cmd := ""
	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	}
	if *openBrowser && cmd != "" {
		go func() {
			time.Sleep(500 * time.Millisecond)
			exec.Command(cmd, "http://127.0.0.1"+*listen).Run()
		}()
	}

	log.Fatal(http.ListenAndServe(*listen, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Page, Page-Size, Total-Pages, query, Total-Items, Query-Duration, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			http.FileServer(http.Dir(*dir)).ServeHTTP(w, r)
		}
	})))
}
