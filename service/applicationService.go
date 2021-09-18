package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/client"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/repo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IService interface {
	GetApplications(ctx context.Context, page int, count int, applicationCriteria model.ApplicationCriteria) (*model.PageableApplicationDto, error)
	GetApplication(ctx context.Context, id int64) (*model.Application, error)
	GetFilterInfo(ctx context.Context) (*model.FilterInfo, error)
	ChangeStatus(ctx context.Context, userId int64, id int64, request model.ChangeStatusRequest) *model.ErrorResponse
}

type Service struct {
	ApplicationRepo repo.IApplicationRepo
	CommentRepo     repo.ICommentRepo
	AdminClient     client.IAdminClient
}

func (s *Service) GetApplications(ctx context.Context, page int, count int, applicationCriteria model.ApplicationCriteria) (*model.PageableApplicationDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetApplications.start")

	offset := page * count

	applications, err := s.ApplicationRepo.GetPageableApplications(offset, count, applicationCriteria)
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get paging applications %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-paging-applications", model.Exception))
	}

	totalCount, err := s.ApplicationRepo.GetTotalCount()
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

	application, err := s.ApplicationRepo.GetApplicationById(id)
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

	applications, err := s.ApplicationRepo.GetApplications()
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

func (s *Service) ChangeStatus(ctx context.Context, userId int64, id int64, request model.ChangeStatusRequest) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.ChangeStatus.start")

	user, err := s.AdminClient.GetUserById(userId)
	if err != nil {
		log.Error("ActionLog.ChangeStatus.warn while call admin client for userId:", userId)
		return &model.ErrorResponse{Code: err.Error(), Status: http.StatusForbidden}
	}

	application, err := s.ApplicationRepo.GetApplicationById(id)
	if err != nil {
		logger.Errorf("ActionLog.ChangeStatus.error: cannot get application for application id %d, %v", id, err)
		return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
	}

	comment := model.Comment{
		Commentator:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Description:   request.Description,
		CommentType:   model.Internal,
		ApplicationId: application.Id,
	}

	if request.Status == model.Hold {
		err := s.CommentRepo.SaveComment(&comment)
		if err != nil {
			logger.Errorf("ActionLog.SaveComment.error: could not save comment details for order %d - %v", application.Id, err)
			return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
		}
	}

	err = ValidationApplicationStatus(application.Status, request.Status)
	if err != nil {
		log.Warn(fmt.Sprintf("ActionLog.ValidationApplicationStatus.error: %s -> %s is not possible", application.Status, request.Status))
		return &model.ErrorResponse{
			Code:   fmt.Sprintf("%s.Invalid status change from %s to %s", model.Exception, application.Status, request.Status),
			Status: http.StatusForbidden,
		}
	}

	application.Status = request.Status
	_, err = s.ApplicationRepo.SaveApplication(application)
	if err != nil {
		return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
	}

	logger.Info("GetApplications.ChangeStatus.end")
	return nil
}

func ValidationApplicationStatus(applicationStatus model.Status, requestStatus model.Status) error {
	if !canBeChangeTo(applicationStatus, requestStatus) {
		log.Error(fmt.Sprintf("ActionLog.ValidationApplicationStatus.error: %s -> %s is not possible", applicationStatus, requestStatus))
		return errors.New(fmt.Sprintf("Invalid stattus from %s to %s", applicationStatus, requestStatus))
	}
	return nil
}

func canBeChangeTo(applicationStatus model.Status, requestStatus model.Status) bool {
	permissions := statusChangePermissions(applicationStatus)
	for _, permission := range permissions {
		if permission == requestStatus {
			return true
		}
	}
	return false
}

func statusChangePermissions(applicationStatus model.Status) []model.Status {
	var permissions []model.Status
	switch applicationStatus {
	case model.Received:
		permissions = append(permissions, model.Inprogress, model.Hold)
		break
	case model.Inprogress:
		permissions = append(permissions, model.Sent, model.Hold)
		break
	case model.Sent:
		break
	case model.Hold:
		permissions = append(permissions, model.Sent, model.Inprogress)
		break
	}
	return permissions
}
