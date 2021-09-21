package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/namrahov/ms-ecourt-go/model"
	"io/ioutil"
	"net/http"
	"time"
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

func SendRequest(ctx context.Context, request *http.Request, response interface{},
	contentType string, method string, requestTimeout time.Duration) error {
	request.Header = GetHeader(ctx)
	request.Header.Del("content-type")
	request.Header.Add("content-type", contentType)

	client := http.Client{
		Timeout: requestTimeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if !isHttpStatus2xx(resp.StatusCode) {
		return fmt.Errorf("client %s returned http status %s", method, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &response)
	}
	return err
}

func GetHeader(ctx context.Context) http.Header {
	header, ok := ctx.Value(model.ContextHeader).(http.Header)
	if !ok {
		return http.Header{}
	}
	return header
}
