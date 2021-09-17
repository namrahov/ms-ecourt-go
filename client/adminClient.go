package client

import (
	"fmt"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/util"
	"github.com/prometheus/common/log"
)

type IAdminClient interface {
	GetUserById(id int64) (*model.UserDto, error)
}

type AdminClient struct{}

func (c *AdminClient) GetUserById(id int64) (*model.UserDto, error) {
	endpoint := fmt.Sprintf("%s%s%d", "https://ufctest.pshb.local:30444/ms-admin", "/users/", id)
	res := SendRequestToClient(endpoint, model.GET, nil)

	if res == nil {
		return nil, nil
	}

	if !IsHttpStatusSuccess(res.StatusCode) {
		log.Error("ActionLog.CheckUserPermission.error when getting data from ms-admin status:", res.StatusCode)
		return nil, nil
	}

	var userDto model.UserDto

	errorResponse := util.DecodeJSON(res.Body, &userDto)

	defer res.Body.Close()
	if errorResponse != nil {
		return nil, nil
	}

	return &userDto, nil
}
