package repo

import "github.com/namrahov/ms-ecourt-go/model"

type IApplicationRepo interface {
	GetApplications(offset int, count int) ([]*model.ApplicationResponse, error)
}

type ApplicationRepo struct{}

func (r ApplicationRepo) GetApplications(offset int, count int) ([]*model.ApplicationResponse, error) {
	var applications []*model.ApplicationResponse

	err := Db.Model(&applications).
		//Limit(offset).
		//Offset(count).
		Select()

	//"limit {OFFSET}, {PAGE_SIZE}";

	return applications, err
}
