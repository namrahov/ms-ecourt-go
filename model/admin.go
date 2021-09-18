package model

type UserPermissionKeyDto struct {
	UserId int64  `json:"userId"`
	Key    string `json:"key"`
}

type AdminResponse struct {
	Status ResponseStatus `json:"status"`
}

type ResponseStatus string

const (
	Error   ResponseStatus = "ERROR"
	Success                = "SUCCESS"
)
