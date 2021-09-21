package model

type Document struct {
	tableName struct{} `sql:"document" pg:",discard_unknown_columns"`

	Id             int64          `sql:"id" json:"id"`
	Description    string         `sql:"description" json:"description"`
	DocumentStatus DocumentStatus `sql:"document_status" json:"documentStatus"`
	RequestType    RequestType    `sql:"request_type" json:"requestType"`
	ApplicationId  int64          `sql:"application_id" json:"-"`
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

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}
