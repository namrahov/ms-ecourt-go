package permission

import (
	"github.com/namrahov/ms-ecourt-go/client"
	"github.com/namrahov/ms-ecourt-go/model"
)

type IService interface {
	HasPermission(userId int64, key string) bool
}

type Service struct {
	AdminClient client.IAdminClient
}

func (s *Service) HasPermission(userId int64, key string) bool {
	var adminResponse, err = s.AdminClient.CheckUserPermission(model.UserPermissionKeyDto{UserId: userId, Key: key})

	if err != nil {
		return false
	}

	return adminResponse.Status == model.Success
}
