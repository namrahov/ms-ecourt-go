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
	GetApplications(ctx context.Context, page int, count int) (*model.PageableApplicationDto, error)
}

type Service struct {
	Repo repo.IApplicationRepo
}

func (s *Service) GetApplications(ctx context.Context, page int, count int) (*model.PageableApplicationDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("GetApplications.GetApplications.start")

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * count

	applications, err := s.Repo.GetApplications(offset, count)
	if err != nil {
		logger.Errorf("ActionLog.GetApplications.error: cannot get applications %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-applications", model.Exception))
	}

	//totalCount

	pageableApplicationDto := model.PageableApplicationDto{
		List: applications,
	}

	/*parsed, _ := time.Parse("2006-01-02 15:04:05", delivery.DeliveryDate)
	delivery.DeliveryDate = parsed.Format(time.RFC3339)*/

	logger.Info("ActionLog.GetApplications.success")
	return &pageableApplicationDto, nil
}

//pageSize -> count
// page -> page

/*
   @Override
    public List<Candidate> getCandidateList(int page) {
        if(page < 1) {
            page = 1;
        }
        int offset = (page - 1) * pageSize;
        return candidateRepository.getCandidateList(offset, pageSize);
    }

    @Override
    public long getCandidatePageCount() {
        long candidateCount = candidateRepository.getCandidateCount();
        long pageCount = candidateCount/pageSize;
        if(candidateCount % pageSize > 0) {
            pageCount++;
        }
        return pageCount;
    }
*/
