package model

type DocumentDto struct {
	Id             int64          `json:"id"`
	Description    string         `json:"description"`
	DocumentStatus DocumentStatus `json:"documentStatus"`
	RequestType    RequestType    `json:"requestType"`
}

type DocumentStatus string

const (
	Waiting  DocumentStatus = "WAITING"
	Uploaded                = "UPLOADED"
	Rejected                = "REJECTED"
)

type RequestType string

const (
	SendAccountInfo      RequestType = "SendAccountInfo"
	SendCollateralInfo               = "SendCollateralInfo"
	SendCompensationInfo             = "SendCompensationInfo"
	SendCustomerInfo                 = "SendCustomerInfo"
	SendDepositInfo                  = "SendDepositInfo"
	SendLoanInfo                     = "SendLoanInfo"
	SendSuretyInfo                   = "SendSuretyInfo"
)
