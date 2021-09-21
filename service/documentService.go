package service

import (
	"context"
	"github.com/namrahov/ms-ecourt-go/client"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
	"html/template"
)

type IDocumentService interface {
	GenerateAct(ctx context.Context, tmpl *template.Template, dto model.TodoPageData) (*model.File, error)
}

type DocumentService struct {
	Html2PdfClient client.IHtml2PdfClient
}

func (s *DocumentService) GenerateAct(ctx context.Context, tmpl *template.Template, dto model.TodoPageData) (*model.File, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetApplication.start")
	/*data := make(map[string]interface{})
	data["data"] = dto*/

	file, err := s.Html2PdfClient.ConvertHtmlToPdf(ctx, tmpl, dto)
	if err != nil {
		return nil, err
	}

	logger.Info("ActionLog.GetApplication.success")

	return file, nil
}
