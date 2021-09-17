package client

import (
	"bytes"
	"crypto/tls"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

func IsHttpStatusSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func SendRequestToClient(endpoint string, httpMethod string, requestBody []byte) *http.Response {
	// Disable G402 "TLS InsecureSkipVerify set true". No need to check server certificate for in-cluster service call
	// #nosec G402
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	request, err := http.NewRequest(httpMethod, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Error("ActionLog.SendRequestToClient.error when prepare new request ", err.Error())
		return nil
	}

	request.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		logger.Error("ActionLog.SendRequestToClient.error when do request ", err.Error())
		return nil
	}

	return res
}
