package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/util"
	"github.com/prometheus/common/log"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
)

type IAdminClient interface {
	CheckUserPermission(dto model.UserPermissionKeyDto) (*model.AdminResponse, *model.ErrorResponse)
	GetUserById(id int64) (*model.UserDto, error)
}

type AdminClient struct{}

func (c *AdminClient) GetUserById(id int64) (*model.UserDto, error) {
	logger.Debug("Get orders by status start")
	endpoint := fmt.Sprintf("%s%s%d", config.Props.AdminClient, "/users/", id)
	userDto, err := sendGetUserDtoRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	logger.Debug("Get orders by status end")
	return userDto, nil

	/*endpoint := fmt.Sprintf("%s%s%d", config.Props.AdminClient, "/users/", id)
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

	return &userDto, nil*/
}

func (c *AdminClient) CheckUserPermission(dto model.UserPermissionKeyDto) (*model.AdminResponse, *model.ErrorResponse) {
	requestBody, err := json.Marshal(dto)

	if err != nil {
		log.Error("ActionLog.CheckUserPermission.error when marshall object ", err.Error())
		return nil, &model.CantMarshalObjectError
	}

	endpoint := config.Props.AdminClient + "/checkUserPermissionByUserIdAndKey"
	resp, err := sendRequest("POST", endpoint, requestBody)

	if resp == nil {
		return nil, &model.UnexpectedError
	}

	if !isHttpStatus2xx(resp.StatusCode) {
		log.Error("ActionLog.CheckUserPermission.error when getting data from ms-admin status:", resp.StatusCode)
		return nil, &model.AdminClientError
	}

	var adminResponse model.AdminResponse

	errorResponse := util.DecodeJSON(resp.Body, &adminResponse)

	defer resp.Body.Close()
	if errorResponse != nil {
		return nil, errorResponse
	}

	return &adminResponse, nil
}

func sendGetUserDtoRequest(endpoint string, requestBody []byte) (*model.UserDto, error) {
	resp, err := sendRequest("GET", endpoint, requestBody)
	if err != nil {
		return nil, err
	}
	if !isHttpStatus2xx(resp.StatusCode) {
		logger.Errorf("sendGetUserDtoRequest client exception %s", resp.Status)
		return nil, errors.New("sendGetUserDtoRequest returned http status " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	var userDto model.UserDto
	err = json.Unmarshal(body, &userDto)
	if err != nil {
		logger.Errorf("Error when unmarshaling body. %s", err.Error())
		return nil, err
	}

	return &userDto, nil
}
