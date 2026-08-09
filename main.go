package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//go:embed game
var gameFiles embed.FS

func init() {
	// Some OS MIME databases don't include these; register explicitly
	// so http.FileServer always sends the correct Content-Type.
	mime.AddExtensionType(".wasm", "application/wasm")
	mime.AddExtensionType(".pck", "application/octet-stream")
}

// godotHeaders wraps an http.Handler and injects the two headers that browsers
// require before they expose SharedArrayBuffer (used by Godot's WASM threads).
type godotHeaders struct{ next http.Handler }

func (g godotHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	g.next.ServeHTTP(w, r)
}

func main() {
	port := availablePort(8080)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	// Sub-FS rooted at "game/" so URLs map directly to file names
	sub, err := fs.Sub(gameFiles, "game")
	if err != nil {
		log.Fatalf("embed sub-filesystem: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", godotHeaders{http.FileServer(http.FS(sub))})

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // streaming large WASM/PCK files; no write deadline
		IdleTimeout:  60 * time.Second,
	}

	// ── Banner ────────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  ⚽  Super Liquid Soccer")
	fmt.Printf("  🌐  %s\n", url)
	fmt.Println("  ⌨   Ctrl+C to quit")
	fmt.Println()

	// ── Open browser after a short delay ─────────────────────────────────────
	go func() {
		time.Sleep(350 * time.Millisecond)
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠  Could not open browser: %v\n", err)
			fmt.Fprintf(os.Stderr, "     Open %s manually.\n", url)
		}
	}()

	// ── Graceful shutdown on Ctrl+C ───────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("\n  Shutting down…")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}

// availablePort returns preferred if it's free, otherwise any free port.
func availablePort(preferred int) int {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", preferred))
	if err == nil {
		ln.Close()
		return preferred
	}
	ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("no available port:", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

// openBrowser launches the default system browser pointing at url.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		// "start" is a shell built-in; invoke via cmd.exe
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default: // Linux, BSD, …
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
