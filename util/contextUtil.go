package util

import (
	"context"
	"github.com/namrahov/ms-ecourt-go/model"
	"net/http"
)

func GetHeader(ctx context.Context) http.Header {
	header, ok := ctx.Value(model.ContextHeader).(http.Header)
	if !ok {
		return http.Header{}
	}
	return header
}
