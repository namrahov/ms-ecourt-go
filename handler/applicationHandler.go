package handler

import (
	"encoding/json"
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/middleware"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/repo"
	"github.com/namrahov/ms-ecourt-go/service"
	"github.com/namrahov/ms-ecourt-go/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type applicationHandler struct {
	Service service.IService
}

func ApplicationHandler(router *mux.Router) *mux.Router {
	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	h := &applicationHandler{
		Service: &service.Service{
			Repo: &repo.ApplicationRepo{},
		},
	}

	router.HandleFunc(config.InternalRootPath, h.getApplications).Methods("GET")

	return router

}

func (h *applicationHandler) getApplications(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	count, err := strconv.ParseInt(r.URL.Query().Get("count"), 10, 64)

	if err != nil {
		log.Error("ActionLog.generateReport.error happened when get user id from header ", err)
		util.HandleError(w, &model.InvalidHeaderError)
		return
	}

	result, err := h.Service.GetApplications(r.Context(), page, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}
