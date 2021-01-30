package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/nicuf/file-processor-api/cache"
	"github.com/nicuf/file-processor-api/handler"
	"github.com/nicuf/file-processor-api/message_queue"
	"github.com/nicuf/file-processor-api/worker"
)

func main() {
	//testCache()
	//testGetNextID()
	l := log.New(os.Stdout, "file-processor-api", log.LstdFlags)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisCache := cache.NewRedisCache(redisClient)
	redisMessageQueue := message_queue.NewRedisMessageQueue(redisClient)

	taskProcessor := worker.NewProcessor(l, redisCache)

	workerPool := worker.NewWorkerPool(2, redisMessageQueue, l, taskProcessor.RunTask)
	workerPool.StartMaster()
	workerPool.StartWorkers()

	h := handler.NewTaskHandler(l, redisCache, redisMessageQueue)

	sm := mux.NewRouter()

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/task/{uuid}", h.AddTask)

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/task/{uuid}", h.GetTaskInfo)
	getR.HandleFunc("/files/{uuid}", h.SearchFiles)
	getR.HandleFunc("/isloop/{uuid}", h.IsLoop)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{}
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.IdleTimeout = 120 * time.Second
	s.Addr = ":7070"
	s.Handler = sm
	s.ErrorLog = l

	go func() {
		l.Println("Starting server")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
