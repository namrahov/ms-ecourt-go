package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/repo"
	log "github.com/sirupsen/logrus"
)

type IService interface {
	GetApplications(ctx context.Context, page int64, count int64) (*model.PageableApplicationDto, error)
}

type Service struct {
	Repo repo.IApplicationRepo
}

func (s *Service) GetApplications(ctx context.Context, page int64, count int64) (*model.PageableApplicationDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetDeliveryByOrderId.start")

	applications, err := s.Repo.GetApplications(page, count)
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get delivery details for order %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-delivery", model.Exception))
	}

	pageableApplicationDto := model.PageableApplicationDto{
		List: *applications,
	}

	/*parsed, _ := time.Parse("2006-01-02 15:04:05", delivery.DeliveryDate)
	delivery.DeliveryDate = parsed.Format(time.RFC3339)*/

	logger.Info("ActionLog.GetDeliveryByOrderId.success")
	return &pageableApplicationDto, nil
}
