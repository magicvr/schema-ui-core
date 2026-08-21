// Command otlp-sink is a minimal OTLP/HTTP collector SINK for observability
// evidence and operator troubleshooting (VP-015 / workspace-015 GOAL-006
// D-001 §2). It accepts POSTs on any path (e.g. /v1/traces), counts request
// bytes and prints one line per request. It is NOT a collector: no parsing,
// no semantics, no storage. Do not use in production.
//
// Usage:
//
//	go run ./cmd/otlp-sink            # listens on :4318 (OTLP/HTTP default)
//	OTLP_SINK_ADDR=:4319 go run ./cmd/otlp-sink
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	addr := os.Getenv("OTLP_SINK_ADDR")
	if addr == "" {
		addr = ":4318"
	}
	var posts, bytes int64
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		n, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			log.Printf("sink: read error %v (url=%s)", err, r.URL)
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}
		posts++
		bytes += n
		log.Printf("sink: POST %s bytes=%d (total posts=%s bytes=%s)", r.URL.Path, n, strconv.FormatInt(posts, 10), strconv.FormatInt(bytes, 10))
		w.WriteHeader(http.StatusOK)
	})
	log.Printf("otlp-sink listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("sink: %v", err)
	}
}