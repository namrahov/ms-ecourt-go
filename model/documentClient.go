package model

type DocumentDto struct {
	Id             int64          `json:"id"`
	Description    string         `json:"description"`
	DocumentStatus DocumentStatus `json:"documentStatus"`
	RequestType    RequestType    `json:"requestType"`
}
