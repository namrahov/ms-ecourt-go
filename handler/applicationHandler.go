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

	router.HandleFunc(config.RootPath+"/applications", h.getApplications).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/get/filter-info", h.getFilterInfo).Methods("GET")

	return router
}

func (h *applicationHandler) getApplications(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		log.Errorf("getApplications.error in parsing page or count: %v\n", err)
		return
	}

	courtName := r.URL.Query().Get("courtName")
	judgeName := r.URL.Query().Get("judgeName")
	person := r.URL.Query().Get("person")
	createDateFrom := r.URL.Query().Get("createDateFrom")
	createDateTo := r.URL.Query().Get("createDateTo")

	var applicationCriteria model.ApplicationCriteria
	applicationCriteria.CourtName = courtName
	applicationCriteria.JudgeName = judgeName
	applicationCriteria.Person = person
	if createDateTo == "" && createDateFrom == "" {
		createDateTo = "2300-01-01"
		createDateFrom = "1000-01-01"
	}
	applicationCriteria.CreateDateFrom = createDateFrom
	applicationCriteria.CreateDateTo = createDateTo

	result, err := h.Service.GetApplications(r.Context(), page, count, applicationCriteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *applicationHandler) getApplication(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Service.GetApplication(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *applicationHandler) getFilterInfo(w http.ResponseWriter, r *http.Request) {

	result, err := h.Service.GetFilterInfo(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
