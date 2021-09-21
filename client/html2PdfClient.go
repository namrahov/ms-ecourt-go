package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/handler"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/util"
	logger "github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

const contentTypeJson = "application/json"
const contentTypeHeader = "Content-type"
const contentDispositionHeader = "Content-Disposition"

type IHtml2PdfClient interface {
	GenerateAct(ctx context.Context, tmpl *template.Template, dto handler.TodoPageData) (*model.File, error)
}

type HtmlToPdfClient struct{}

func (f *HtmlToPdfClient) GenerateAct(ctx context.Context, tmpl *template.Template, dto handler.TodoPageData) (*model.File, error) {
	logger.Debug("Generate act start")

	endpoint := fmt.Sprintf("%s%s", config.Props.Html2PdfEndpoint, "/v1/internal/html-2-pdf")
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}
	request.Header = util.GetHeader(ctx)

	request.Header.Add(contentTypeHeader, contentTypeJson)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if !isHttpStatus2xx(resp.StatusCode) {
		return nil, fmt.Errorf("client DownloadFile returned http status %s", resp.Status)
	}
	file := new(model.File)
	file.Type = resp.Header.Get(contentTypeHeader)
	file.Name = resp.Header.Get(contentDispositionHeader)
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	file.Size = int64(size)
	if err != nil {
		logger.Errorf("ActionLog.DownloadFile.error : when trying to convert size %v", err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("ActionLog.DownloadFile.error : when trying to read body %v", err.Error())
		return nil, err
	}
	file.Content = body
	defer resp.Body.Close()

	logger.Debug("Generate act end")
	return file, nil
}
