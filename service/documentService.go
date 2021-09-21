package service

import (
	"context"
	"github.com/namrahov/ms-ecourt-go/client"
	"github.com/namrahov/ms-ecourt-go/handler"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
	"html/template"
)

type IDocumentService interface {
	GenerateAct(ctx context.Context, tmpl *template.Template, dto handler.TodoPageData) (*model.File, error)
}

type DocumentService struct {
	Html2PdfClient client.IHtml2PdfClient
}

func (s *DocumentService) GenerateAct(ctx context.Context, tmpl *template.Template, dto handler.TodoPageData) (*model.File, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetApplication.start")

	file, err := s.Html2PdfClient.GenerateAct(ctx, tmpl, dto)
	if err != nil {
		return nil, err
	}

	logger.Info("ActionLog.GetApplication.success")

	return file, nil
}
