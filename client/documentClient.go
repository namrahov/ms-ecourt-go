package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/config"
	"github.com/namrahov/ms-ecourt-go/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type IDocumentClient interface {
	GetDocument(ctx context.Context, id int64) (*model.DocumentDto, error)
}

type DocumentClient struct {
}

func (c *DocumentClient) GetDocument(ctx context.Context, id int64) (*model.DocumentDto, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Debug("Get delivery by orderId from ms-card-delivery start")

	//pat variable da id gonderirik, body-de hec ne gondermirik
	endpoint := fmt.Sprintf("%s/v1/internal/ecourt/documents/%d", config.Props.DocumentEndpoint, id)
	resp, err := sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if !isHttpStatus2xx(resp.StatusCode) {
		logger.Errorf("Card delivery GetDeliveryByOrderId client exception %s", resp.Status)
		return nil, errors.New("card delivery client GetDeliveryByOrderId returned http status " + resp.Status)
	}

	//gelen responsedan body-ni gotururuk
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var documentDto model.DocumentDto

	err = json.Unmarshal(body, &documentDto)
	if err != nil {
		return nil, err
	}

	log.Debug("Get delivery by orderId from ms-card-delivery success")
	return &documentDto, nil
}
