package model

const (
	HeaderKeyCustomerID = "DP-Customer-ID"
	HeaderKeyUserID     = "DP-User-ID"
	HeaderKeyUserAgent  = "User-Agent"
	HeaderKeyUserIP     = "X-Forwarded-For"
	HeaderKeyRequestID  = "requestid"
)

const (
	LoggerKeyRequestID  = "REQUEST_ID"
	LoggerKeyOperation  = "OPERATION"
	LoggerKeyCustomerID = "CUSTOMER_ID"
	LoggerKeyUserID     = "USER_ID"
	LoggerKeyUserIP     = "USER_IP"
	LoggerKeyUserAgent  = "USER_AGENT"
	ContextLogger       = "contextLogger"
	ContextHeader       = "contextHeader"
)

const Exception = "error.card-delivery"
const AzeIsdCode = "994"

// Time format for databases
const (
	CustomTimeFormat = "2006-01-02T15:04:05"
	JsonTimeFormat   = "2006-01-02T15:04:05Z"
	DbTimeFormat     = "2006-01-02 15:04:05"
)

const (
	PashaBankCustomerId    string = "00000001"
	PortBakuBranchCode     string = "006"
	ShareholdersBranchCode string = "012"
	HeadOfficeBranchCode   string = "001"
)

const (
	BusinessCardsSheet = "Business cards"
	SalaryCardsSheet   = "Salary cards"
)
