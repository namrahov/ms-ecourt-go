package model

type StatusHistory struct {
	tableName struct{} `sql:"status_history" pg:",discard_unknown_columns"`

	Id              int64         `sql:"id" json:"id"`
	Response        int64         `sql:"response" json:"response"`
	ServiceResponse int64         `sql:"service_response" json:"serviceResponse"`
	ErrorMessage    string        `sql:"error_message" json:"errorMessage"`
	CreatedAt       string        `sql:"created_at" json:"createdAt"`
	UpdatedAt       string        `sql:"updated_at" json:"updatedAt"`
	Application     []Application `json:"applications"`
}
