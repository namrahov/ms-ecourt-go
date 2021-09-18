package client

import (
	"bytes"
	"crypto/tls"
	"net/http"
)

func isHttpStatus2xx(httpStatusCode int) bool {
	statusOK := httpStatusCode >= 200 && httpStatusCode < 300
	return statusOK
}

func sendRequest(httpMethod string, endpoint string, requestBody []byte) (*http.Response, error) {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	request, err := http.NewRequest(httpMethod, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Add("content-type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
