package client

import (
	"encoding/json"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/util"
	"github.com/prometheus/common/log"
)

type IAdminClient interface {
	CheckUserPermission(dto model.UserPermissionKeyDto) (*model.AdminResponse, *model.ErrorResponse)
	GetUserById(id int64) (*model.UserDto, error)
}

type AdminClient struct{}

func (c *AdminClient) GetUserById(id int64) (*model.UserDto, error) {
	endpoint := fmt.Sprintf("%s%s%d", config.Props.AdminClient, "/users/", id)
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

func (c *AdminClient) CheckUserPermission(dto model.UserPermissionKeyDto) (*model.AdminResponse, *model.ErrorResponse) {
	requestBody, err := json.Marshal(dto)

	if err != nil {
		log.Error("ActionLog.CheckUserPermission.error when marshall object ", err.Error())
		return nil, &model.CantMarshalObjectError
	}

	endpoint := config.Props.AdminClient + "/checkUserPermissionByUserIdAndKey"
	res := SendRequestToClient(endpoint, model.POST, requestBody)

	if res == nil {
		return nil, &model.UnexpectedError
	}

	if !IsHttpStatusSuccess(res.StatusCode) {
		log.Error("ActionLog.CheckUserPermission.error when getting data from ms-admin status:", res.StatusCode)
		return nil, &model.AdminClientError
	}

	var adminResponse model.AdminResponse

	errorResponse := util.DecodeJSON(res.Body, &adminResponse)

	defer res.Body.Close()
	if errorResponse != nil {
		return nil, errorResponse
	}

	return &adminResponse, nil
}
