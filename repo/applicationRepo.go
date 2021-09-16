package repo

import (
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) (*[]model.Application, error)
	GetTotalCount() (int, error)
	GetApplicationById(id int64) (*model.Application, error)
}

type ApplicationRepo struct {
}

func (r ApplicationRepo) GetApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) (*[]model.Application, error) {
	var applications []model.Application
	err := Db.Model(&applications).
		Column("application.*", "Comments").
		Where("court_name like ?", "%"+applicationCriteria.CourtName+"%").
		Where("judge_name like ?", "%"+applicationCriteria.JudgeName+"%").
		Where("person like ?", "%"+applicationCriteria.Person+"%").
		Where("created_at::DATE >= ?", applicationCriteria.CreateDateFrom).
		Where("created_at::DATE <= ?", applicationCriteria.CreateDateTo).
		Limit(count).
		Offset(offset).
		Select()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &applications, err
}

func (r ApplicationRepo) GetTotalCount() (int, error) {
	var totalCount int
	var applications []model.Application
	totalCount, err := Db.Model(&applications).Count()
	if err != nil {
		log.Fatal(err)
	}
	return totalCount, nil
}

func (r ApplicationRepo) GetApplicationById(id int64) (*model.Application, error) {
	var application model.Application
	err := Db.Model(&application).
		Column("application.*", "Comments").
		Where("id = ?", id).
		Select()
	if err != nil {
		log.Fatal(err)
	}

	return &application, nil
}
