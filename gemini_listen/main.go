package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func runGemini(prompt string) (string, error) {
	cmd := exec.Command("gemini", "--yolo", "--prompt", prompt, "-m", "gemini-2.5-flash")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %v", strings.TrimSpace(string(out)), err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	port := 8765
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil {
			port = p
		}
	} else if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"ok":false,"error":"Failed to read body"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			parsed = map[string]interface{}{"raw": string(body)}
		}

		message, _ := parsed["message"].(string)

		if message == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"ok":false,"error":"No message provided"}`))
			return
		}

		reply, err := runGemini(message)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp, _ := json.Marshal(map[string]interface{}{
				"ok":    false,
				"error": err.Error(),
			})
			w.Write(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(map[string]interface{}{
			"ok":    true,
			"reply": reply,
		})
		w.Write(resp)
	})

	// Handle all other routes as 404
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on http://127.0.0.1:%d/event", port)
		if wd, err := os.Getwd(); err == nil {
			log.Printf("Gemini CLI working directory: %s", wd)
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
