package handler

import (
	"encoding/json"
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/namrahov/ms-ecourt-go/client"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/middleware"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/repo"
	"github.com/namrahov/ms-ecourt-go/service"
	"github.com/namrahov/ms-ecourt-go/service/permission"
	"github.com/namrahov/ms-ecourt-go/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type applicationHandler struct {
	Service           service.IService
	PermissionService permission.IService
}

func ApplicationHandler(router *mux.Router) *mux.Router {
	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	h := &applicationHandler{
		Service: &service.Service{
			ApplicationRepo: &repo.ApplicationRepo{},
			CommentRepo:     &repo.CommentRepo{},
			AdminClient:     &client.AdminClient{},
		},
		PermissionService: &permission.Service{
			AdminClient: &client.AdminClient{},
		},
	}

	router.HandleFunc(config.RootPath+"/applications", h.getApplications).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/get/filter-info", h.getFilterInfo).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}/change-status", h.changeStatus).Methods("GET")

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

func (h *applicationHandler) changeStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.Header.Get(model.UserIdHeader), 10, 64)

	if err != nil {
		log.Error("ActionLog.generateReport.error happened when get user id from header ", err)
		util.HandleError(w, &model.InvalidHeaderError)
		return
	}

	hasPermission := h.PermissionService.HasPermission(userId, model.GenerateReportPermissionKey)

	if !hasPermission {
		log.Error("ActionLog.generateReport.error access is denied for userId:", userId)
		util.HandleError(w, &model.AccessDeniedError)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request model.ChangeStatusRequest
	err = util.DecodeBody(w, r, &request)
	if err != nil {
		return
	}

	err = h.Service.ChangeStatus(r.Context(), userId, id, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Add("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(result)
}
