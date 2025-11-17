package internal

import (
	"net/http"
	"sync"
)

// --- SSE broadcaster (simple)
type sseHub struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

func NewSSEHub() *sseHub {
	return &sseHub{clients: map[chan string]struct{}{}}
}

func (h *sseHub) AddClient(ch chan string) {
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
}

func (h *sseHub) RemoveClient(ch chan string) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

func (h *sseHub) Broadcast(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		// Non-blocking send so one slow client doesn't block others
		select {
		case ch <- msg:
		default:
		}
	}
}

// SSE handler: upgrades an HTTP connection to an event stream.
func (h *sseHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	msgCh := make(chan string, 1)
	h.AddClient(msgCh)
	defer h.RemoveClient(msgCh)

	// Send a message(keeps connection alive in some proxies)
	_, _ = w.Write([]byte(": connected\n\n"))

	_, _ = w.Write([]byte("retry: 200\n\n"))
	fl.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			// SSE format: "data: <payload>\n\n"
			_, _ = w.Write([]byte("data: " + msg + "\n\n"))
			fl.Flush()
		}
	}
}
