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
	GetApplication(ctx context.Context, id int64) (*model.Application, error)
	GetFilterInfo(ctx context.Context) (*model.FilterInfo, error)
}

type Service struct {
	Repo repo.IApplicationRepo
}

func (s *Service) GetApplications(ctx context.Context, page int, count int, applicationCriteria model.ApplicationCriteria) (*model.PageableApplicationDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetApplications.start")

	offset := page * count

	applications, err := s.Repo.GetPageableApplications(offset, count, applicationCriteria)
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get paging applications %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-paging-applications", model.Exception))
	}

	totalCount, err := s.Repo.GetTotalCount()
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get total count %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-total-count", model.Exception))
	}

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

	logger.Info("ActionLog.GetApplications.success")
	return &pageableApplicationDto, nil
}

func (s *Service) GetApplication(ctx context.Context, id int64) (*model.Application, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetApplication.start")

	application, err := s.Repo.GetApplicationById(id)
	if err != nil {
		logger.Errorf("ActionLog.GetApplication.error: cannot get application for application id %d, %v", id, err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-application", model.Exception))
	}

	logger.Info("ActionLog.GetApplication.success")

	return application, nil
}

func (s *Service) GetFilterInfo(ctx context.Context) (*model.FilterInfo, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetFilterInfo.start")

	applications, err := s.Repo.GetApplications()
	if err != nil {
		logger.Errorf("ActionLog.GetFilterInfo.error: cannot get applications %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-applications", model.Exception))
	}

	var courts []string
	var judges []string
	setCourts := make(map[string]bool)
	setJudges := make(map[string]bool)

	for _, application := range *applications {
		setCourts[application.CourtName] = true
		setJudges[application.JudgeName] = true
	}

	for k := range setCourts {
		courts = append(courts, k)
	}

	for k := range setJudges {
		judges = append(judges, k)
	}

	filterInfo := model.FilterInfo{
		Courts: courts,
		Judges: judges,
	}

	logger.Info("ActionLog.GetFilterInfo.success")
	return &filterInfo, nil
}
