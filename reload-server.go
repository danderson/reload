//go:build ignore

// This is just a test server that demonstrates using the reload
// handler.
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/danderson/reload"
)

func main() {
	rl := &reload.Reloader{}

	http.HandleFunc("/", index)
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		rl.Reload()
		http.Redirect(w, r, "/reload-ui", http.StatusFound)
	})
	http.HandleFunc("/reload-ui", reloadUI)
	http.Handle("/live", rl)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html>
      <head>
        <script defer type="text/javascript" src="/live"></script>
      </head>
      <body>
        <p>At the tone, the time is %v, Coordinated Universal Time</p>

        <p><a href="/reload-ui">Manual reloads</a></p>
      </body>
    </html>`,
		time.Now().UTC().Format("03:04:05"))
}

func reloadUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, `<!DOCTYPE html>
	<html>
      <head>
        <script defer type="text/javascript" src="/live"></script>
      </head>
      <body>
        <form action="/reload" method="get">
          <input type="submit" value="Trigger reload">
        </form>
      </body>
    </html>`)
}
