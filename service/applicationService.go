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
	GetApplications(ctx context.Context, page int, count int, applicationCriteria model.ApplicationCriteria) (*model.PageableApplicationDto, error)
}

type Service struct {
	Repo repo.IApplicationRepo
}

func (s *Service) GetApplications(ctx context.Context, page int, count int, applicationCriteria model.ApplicationCriteria) (*model.PageableApplicationDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetApplications.start")

	offset := page * count

	applications, err := s.Repo.GetApplications(offset, count, applicationCriteria)
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get applications %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-applications", model.Exception))
	}

	totalCount, err := s.Repo.GetTotalCount()

	lastPageNumber := totalCount / count

	var hasNextPage bool
	check := (totalCount - (page * count)) / count
	if check > 0 {
		hasNextPage = true
	} else {
		hasNextPage = false
	}

	pageableApplicationDto := model.PageableApplicationDto{
		List:           applications,
		HasNextPage:    hasNextPage,
		LastPageNumber: lastPageNumber,
		TotalCount:     totalCount,
	}

	/*parsed, _ := time.Parse("2006-01-02 15:04:05", delivery.DeliveryDate)
	delivery.DeliveryDate = parsed.Format(time.RFC3339)*/

	logger.Info("ActionLog.GetApplications.success")
	return &pageableApplicationDto, nil
}
