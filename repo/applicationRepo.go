package repo

import "github.com/namrahov/ms-ecourt-go/model"

type IApplicationRepo interface {
	GetApplications(page int64, count int64) (*[]model.ApplicationResponse, error)
}

type ApplicationRepo struct{}

func (r ApplicationRepo) GetApplications(page int64, count int64) (*[]model.ApplicationResponse, error) {
	var applications []model.ApplicationResponse

	err := Db.Model(&applications).
		Limit(page).
		Offset(count).
		Select()

	//"limit {OFFSET}, {PAGE_SIZE}";

	return &applications, err
}
