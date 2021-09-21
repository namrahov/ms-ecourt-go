package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/prometheus/common/log"
	logger "github.com/sirupsen/logrus"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

const contentTypeJson = "application/json"
const contentTypeHeader = "Content-type"
const contentDispositionHeader = "Content-Disposition"

type IHtml2PdfClient interface {
	ConvertHtmlToPdf(ctx context.Context, tmpl *template.Template, dto model.TodoPageData) (*model.File, error)
}

type HtmlToPdfClient struct{}

func (f *HtmlToPdfClient) ConvertHtmlToPdf(ctx context.Context, tmpl *template.Template, dto model.TodoPageData) (*model.File, error) {
	logger.Debug("Generate act start")

	endpoint := fmt.Sprintf("%s%s", config.Props.Html2PdfEndpoint, "/v1/internal/html-2-pdf")

	bodyBuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuf)
	header := textproto.MIMEHeader{}
	header.Add(contentTypeHeader, contentTypeJson)
	header.Add(contentDispositionHeader, "form-data; name=\"data\";")
	r, err := bodyWriter.CreatePart(header)
	if err != nil {
		log.Errorf("Error while creating form request")
		return nil, err
	}
	_, err = r.Write([]byte(fmt.Sprintf("{\"path\":\"%s\",\"rewrite\":\"true\"}", dto.PageTitle)))
	if err != nil {
		log.Errorf("Error while writing request data")
		return nil, err
	}

	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", bodyWriter.Boundary())
	_ = bodyWriter.Close()
	var fileInfo []model.FileInfo
	request, err := http.NewRequest("POST", endpoint, bodyBuf)
	err = SendRequest(ctx, request, &fileInfo, contentType, "UploadFile", 60*time.Second)
	if err != nil {
		log.Errorf("Error while sending request")
		return nil, err
	}
	log.Debug("Upload file to file storage success")

	/*request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(nil))
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
	defer resp.Body.Close()*/

	logger.Debug("Generate act end")
	file := new(model.File)
	return file, nil
}
