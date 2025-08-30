package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"project_go-task-queue/pkg/handler"
	"project_go-task-queue/pkg/worker"
	"strconv"
	"syscall"
)

func main() {
	queueSize := 64
	if qs := os.Getenv("QUEUE_SIZE"); qs != "" {
		if v, err := strconv.Atoi(qs); err == nil {
			queueSize = v
		}
	}

	workers := 4
	if w := os.Getenv("WORKERS"); w != "" {
		if v, err := strconv.Atoi(w); err == nil {
			workers = v
		}
	}

	tq := worker.NewTaskQueue(workers, queueSize)
	tq.Start()

	h := handler.NewHandler(tq)

	mux := http.NewServeMux()
	mux.HandleFunc("/enqueue", h.Enqueue)
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/tasks", h.Tasks)

	server := &http.Server{Addr: ":8080", Handler: mux}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server started at :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received")
	tq.Stop()
	log.Println("All workers done, exiting")
}
