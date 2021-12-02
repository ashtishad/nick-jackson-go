package handlers

import (
	"log"
	"net/http"
)

type bye struct {
	l *log.Logger
}

func Newbye(l *log.Logger) *bye {
	return &bye{l}
}

func (g *bye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.l.Println("bye!")
	w.Write([]byte("bye!!"))
}
