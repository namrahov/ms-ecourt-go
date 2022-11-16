package handler

import (
	"crypto/sha1"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type documentHandler struct {
	DocumentService   service.IDocumentService
	PermissionService permission.IService
}

func DocumentHandler(router *mux.Router) *mux.Router {
	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	h := &documentHandler{
		DocumentService: &service.DocumentService{
			ApplicationRepo: &repo.ApplicationRepo{},
		},
		PermissionService: &permission.Service{
			AdminClient: &client.AdminClient{},
		},
	}

	router.HandleFunc(config.RootPath+"/documents/generate-act", h.generateAct).Methods("POST")
	router.HandleFunc(config.RootPath+"/documents/get-report", h.getReport).Methods("GET")
	router.HandleFunc(config.RootPath+"/documents/upload", h.uploadExcel).Methods("POST")
	router.HandleFunc(config.RootPath+"/documents/upload-file", h.uploadFile).Methods("POST")

	return router
}

func (h *documentHandler) generateAct(w http.ResponseWriter, r *http.Request) {
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

	dto := model.TodoPageData{
		PageTitle: "My list",
		Todos: []model.Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	file, errNew := h.DocumentService.GenerateAct(r.Context(), dto)
	if errNew != nil {
		http.Error(w, errNew.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		http.Error(w, "No pdf exists", http.StatusNotFound)
		return
	}

	_, _ = w.Write(file)
}

func (h *documentHandler) getReport(w http.ResponseWriter, r *http.Request) {
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

	file := new(excelize.File)

	file, clientErr := h.DocumentService.GenerateReportOfLightApplication(r.Context())

	if clientErr != nil {
		http.Error(w, clientErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add(model.ContentDispositionString, model.AttachmentFilename)
	w.Header().Add(model.ContentTypeString, model.ExcelType)
	fileErr := file.Write(w)

	if fileErr != nil {
		log.Error("ActionLog.WriteExcelHeaders.error happened when write excel file")
	}
}

func (h *documentHandler) uploadExcel(w http.ResponseWriter, r *http.Request) {
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

	xlsx, meta, e := r.FormFile("file")
	if e != nil {
		log.Error("ActionLog.uploadExcel.error can't read file ", e.Error())
		return
	}

	const maxUploadSize = 25 * 1024
	if meta != nil {
		if meta.Size > maxUploadSize {
			log.Error("ActionLog.uploadExcel.error file size so much ", meta.Size)
			return
		}
	}

	erro := h.DocumentService.UploadExcel(r.Context(), xlsx)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "multipart/form-data")
	w.WriteHeader(204)
}

func (h *documentHandler) uploadFile(w http.ResponseWriter, req *http.Request) {
	mf, fh, err := req.FormFile("nf")
	if err != nil {
		fmt.Println(err)
	}

	defer mf.Close()

	// crate sha for file name
	ext := strings.Split(fh.Filename, ".")[1]

	hs := sha1.New()

	io.Copy(hs, mf)

	fName := fmt.Sprintf("%x", hs.Sum(nil)) + "." + ext

	//create new file
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := filepath.Join(wd, "public", fName)

	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	defer nf.Close()
	//copy

	seek, err := mf.Seek(0, 0)
	if err != nil {
		fmt.Println(seek, "  ", err)
		return
	}

	io.Copy(nf, mf)

}
