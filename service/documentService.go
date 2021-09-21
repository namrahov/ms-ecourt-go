package service

import (
	"bytes"
	"context"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
	"html/template"
)

type IDocumentService interface {
	GenerateAct(ctx context.Context, dto model.TodoPageData) ([]byte, error)
}

type DocumentService struct{}

func (s *DocumentService) GenerateAct(ctx context.Context, dto model.TodoPageData) ([]byte, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GenerateAct.start")

	tmpl := template.Must(template.ParseFiles("layout.html"))

	var data = make(map[string]interface{})
	data["PageTitle"] = dto.PageTitle
	data["Todos"] = dto.Todos

	file, err := GeneratePdf(ctx, data, tmpl)
	if err != nil {
		logger.Error("ActionLog.GenerateAct.error cant generate pdf", err)
		return nil, err
	}

	logger.Info("ActionLog.GenerateAct.success")
	return file, nil
}

func GeneratePdf(ctx context.Context, data map[string]interface{}, tmpl *template.Template) ([]byte, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GeneratePdf.start")
	var tmplBytes bytes.Buffer

	tmpl.Execute(&tmplBytes, data)

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		logger.Error("ActionLog.ConvertHTML2PDF.error ", err)
		return nil, err
	}

	pdfg.AddPage(wkhtml.NewPageReader(&tmplBytes))

	err = pdfg.Create()
	if err != nil {
		logger.Error("ActionLog.ConvertHTML2PDF.error ", err)
		return nil, err
	}

	logger.Info("ActionLog.GeneratePdf.end")
	return pdfg.Bytes(), nil
}
