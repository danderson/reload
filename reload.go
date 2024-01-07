// Package reload provides an HTTP handler that implements simple live
// reloading for Go web servers.
//
// To use, register this handler in your HTTP server, then load it via
// a <script> tag in your HTML. The page will be reloaded whenever the
// server restarts, and when explicitly requested by invoking the
// handler's Reload function.
package reload

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Reloader is a live-reloading HTTP handler.
type Reloader struct {
	mu      sync.Mutex
	cookie  []byte
	refresh chan struct{}
}

//go:embed reload.js
var js []byte

// ServeHTTP serves the live reloading script and websocket handler.
func (rl *Reloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		rl.socket(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Cache-Control", "no-store")
	cookie, _ := rl.getState()
	fmt.Fprintf(w, "(%s)(%q)", js, cookie)
}

// Reload triggers an immediate reload of all clients.
func (rl *Reloader) Reload() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.reloadLocked()
}

func (rl *Reloader) socket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// upgrader has already responded to client
		return
	}
	defer conn.Close()

	cookie, refresh := rl.getState()

	for {
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, cookie); err != nil {
			return
		}
		cookie = nil

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}

		select {
		case <-time.After(5 * time.Second):
		case <-refresh:
			cookie, refresh = rl.getState()
		case <-r.Context().Done():
			return
		}
	}
}

func (rl *Reloader) getState() (cookie []byte, refresh <-chan struct{}) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.cookie == nil {
		rl.reloadLocked()
	}
	return rl.cookie, rl.refresh
}

func (rl *Reloader) reloadLocked() {
	if rl.refresh != nil {
		close(rl.refresh)
	}
	rl.cookie = []byte(strconv.FormatInt(time.Now().UnixNano(), 36))
	rl.refresh = make(chan struct{})
}
