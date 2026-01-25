package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	workerCount int
	httpServer  *http.Server
	ticker      *time.Ticker
	requestChan chan *http.Request
}

type Option func(*server)

func WithWorkerCount(cont int) Option {
	return func(s *server) {
		s.workerCount = cont
	}
}

func WithTicker(sec time.Duration) Option {
	return func(s *server) {
		s.ticker = time.NewTicker(sec * time.Second)
	}
}
func newServer(port int, opts ...Option) *server {
	requestChan := make(chan *http.Request, 100)
	s := &server{
		requestChan: requestChan,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *server) Start(ctx context.Context) error {

	addr := fmt.Sprintf(":%d", 3005)
	mux := http.NewServeMux()
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		s.requestChan <- r
		w.Write([]byte("Queued work..."))
	})

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	for i := 0; i < s.workerCount; i++ {
		go s.worker(ctx, i)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Cache warmer stopping...")
				return
			case <-s.ticker.C:
				fmt.Println("Warming cache...")
			}
		}
	}()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return nil
}

func (s *server) worker(ctx context.Context, id int) {
	for {
		select {
		case req := <-s.requestChan:
			time.Sleep(5 * time.Second)
			fmt.Printf("Worker %d processing request: %s\n", id, req.URL.Path)

		case <-ctx.Done():
			fmt.Println("Worker", id, "shutting down")
			return
		}
	}
}

func (s *server) Stop(ctx context.Context) error {
	fmt.Println("Shutting down...")
	// Stop accepting new HTTP requests
	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Println("HTTP shutdown error:", err)
	}
	// Stop ticker
	if s.ticker != nil {
		s.ticker.Stop()
	}
	// Wait for ongoing requests to finish
	for {
		if len(s.requestChan) == 0 {
			break
		}
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown timeout. Forcing exit.")
			return ctx.Err()
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}

	fmt.Println("Graceful shutdown complete.")
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := newServer(3005, WithWorkerCount(5), WithTicker(10))
	err := srv.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	<-done
	log.Print("Server Stopped")
	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
