package handler

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type healthHandler struct{}

func HandleHealthRequest(router *mux.Router) {

	h := &healthHandler{}
	router.HandleFunc("/readiness", h.Health)
	router.HandleFunc("/health", h.Health)
	router.Handle("/metrics", promhttp.Handler())

}

func (*healthHandler) Health(w http.ResponseWriter, r *http.Request) {

}
