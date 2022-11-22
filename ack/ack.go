package ack

import (
	"fmt"
	"net/http"

	"github.com/haleyrc/ygg/server"
)

func NewAPI() []server.Route {
	return []server.Route{
		{Methods: []string{"POST"}, Path: "/syn", Func: ackJSON},
	}
}

func NewSite() []server.Route {
	return []server.Route{
		{Methods: []string{"GET"}, Path: "/syn", Func: ackPage},
	}
}

func ackPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ACK")
}

func ackJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"ACK"}`)
}
