package service

import (
	"context"
	"github.com/namrahov/ms-ecourt-go/repo"
)

type IService interface {
	GetApplications(ctx context.Context, page int64, count int64) (*model.Delivery, error)
}

type Service struct {
	Repo repo.IApplicationRepo
}
