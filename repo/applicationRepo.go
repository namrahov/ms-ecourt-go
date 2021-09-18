package repo

import (
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
)

type IApplicationRepo interface {
	GetPageableApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) (*[]model.Application, error)
	GetTotalCount() (int, error)
	GetApplicationById(id int64) (*model.Application, error)
	GetApplications() (*[]model.Application, error)
	SaveApplication(application *model.Application) (*model.Application, error)
}

type ApplicationRepo struct {
}

func (r ApplicationRepo) GetPageableApplications(offset int, count int, applicationCriteria model.ApplicationCriteria) (*[]model.Application, error) {
	var applications []model.Application
	err := Db.Model(&applications).
		Column("application.*", "Comments", "Documents").
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
		Column("application.*", "Comments", "Documents").
		Where("id = ?", id).
		Select()
	if err != nil {
		log.Fatal(err)
	}

	return &application, nil
}

func (r ApplicationRepo) GetApplications() (*[]model.Application, error) {
	var applications []model.Application
	err := Db.Model(&applications).
		Column("application.*", "Comments", "Documents").
		Select()
	if err != nil {
		log.Fatal(err)
	}

	return &applications, nil
}

func (r ApplicationRepo) SaveApplication(application *model.Application) (*model.Application, error) {
	_, err := Db.Model(application).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		log.Fatal(err)
	}
	return application, nil
}
