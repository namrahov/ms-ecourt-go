package model

type PageableApplicationDto struct {
	List       []*ApplicationResponse `json:"list"`
	TotalCount int64                  `json:"totalCount"`
}

type ApplicationResponse struct {
	tableName struct{} `sql:"application" pg:",discard_unknown_columns"`

	Id             int64         `sql:"id"  json:"id"`
	RequestId      int64         `sql:"request_id" json:"requestId"`
	CheckedId      int64         `sql:"checked_id" json:"checkedId"`
	Person         string        `sql:"person" json:"person"`
	CustomerType   CustomerType  `sql:"customer_type" json:"customerType"`
	CustomerName   string        `sql:"customer_name" json:"customerName"`
	FilePath       string        `sql:"file_path" json:"filePath"`
	CourtName      string        `sql:"court_name" json:"courtName"`
	JudgeName      string        `sql:"judge_name" json:"judgeName"`
	DecisionNumber string        `sql:"decision_number" json:"decisionNumber"`
	Note           string        `sql:"note" json:"note"`
	Status         Status        `sql:"status" json:"status"`
	Deadline       string        `sql:"deadline" json:"deadline"`
	AssigneeId     int64         `sql:"assignee_id" json:"assigneeId"`
	Priority       Priority      `sql:"priority" json:"priority"`
	AssigneeName   string        `sql:"assignee_name" json:"assigneeName"`
	Comments       []CommentDto  `sql:"-" json:"comments"`
	Documents      []DocumentDto `sql:"-" json:"documents"`
	BeginDate      string        `sql:"begin_date" json:"beginDate"`
	EndDate        string        `sql:"end_date" json:"endDate"`
	CreatedAt      string        `sql:"created_at" json:"createdAt"`
	BankDetails    []BankDetail  `sql:"-" json:"bankDetails"`
}

type CustomerType string

const (
	Person   CustomerType = "PERSON"
	Taxpayer              = "TAXPAYER"
)

type Status string

const (
	Received   Status = "RECEIVED"
	Inprogress        = "IN_PROGRESS"
	Sent              = "SENT"
	Hold              = "HOLD"
)

type Priority string

const (
	Standard Priority = "STANDARD"
	High              = "HIGH"
)

type BankDetail struct {
	CustomerNo   string `sql:"customer_no" json:"customerNo"`
	CustomerName string `sql:"customer_name" json:"customerName"`
	RequestDate  string `sql:"request_date" json:"requestDate"`
}
