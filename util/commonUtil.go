package util

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func HandleError(w http.ResponseWriter, err *model.ErrorResponse) {
	w.Header().Add(model.ContentTypeString, model.JSONType)
	w.WriteHeader(err.Status)
	encodeErr := json.NewEncoder(w).Encode(err)

	if encodeErr != nil {
		log.Error("ActionLog.HandleError.error happened when encode json")
	}
}

func DecodeJSON(r io.Reader, data interface{}) *model.ErrorResponse {
	log.Debug("ActionLog.DecodeJSON.start")

	err := json.NewDecoder(r).Decode(&data)

	if err != nil {
		log.Error("ActionLog.DecodeJSON.error can't parse data " + err.Error())
		return &model.CantParseDataError
	}

	log.Debug("ActionLog.DecodeJSON.success")
	return nil
}

func WriteExcelHeaders(w http.ResponseWriter, file *excelize.File) {
	w.Header().Add(model.ContentDispositionString, model.AttachmentFilename)
	w.Header().Add(model.ContentTypeString, model.ExcelType)
	err := file.Write(w)

	if err != nil {
		log.Error("ActionLog.WriteExcelHeaders.error happened when write excel file")
	}
}
