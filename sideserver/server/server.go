package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var helloWorld = map[string]string{"hello": "world"}

type Server struct {
	server http.Server
	cc     chan int
}

type appContext struct {
	contextChanel chan int
}

func (a *appContext) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	a.contextChanel <- 1
}

func (a *appContext) TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(helloWorld)
	_, err := w.Write(js)
	if err != nil {
		fmt.Println("SHIT")
	}
}

func NewServer() *Server {
	r := mux.NewRouter()
	channel := make(chan int, 1)
	c := &appContext{contextChanel: channel}
	r.HandleFunc("/", c.TestHandler)
	r.HandleFunc("/shutdown", c.ShutdownHandler)
	return &Server{
		server: http.Server{
			Handler: r,
			Addr:    "127.0.0.1:8888",
		},
		cc: channel,
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	<-s.cc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Println("SHITTTTTTT Shutting down")
	}
	log.Println("Shutdown")
}
