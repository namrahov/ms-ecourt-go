package handler

import (
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/namrahov/ms-ecourt-go/middleware"
)

func ApplicationHandler(router *mux.Router) *mux.Router {
	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	return router

}
