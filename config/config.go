package config

import (
	"github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"
)

const RootPath = "/v1/cards/delivery"
const InternalRootPath = "/internal/v1/cards/delivery"

type Args struct {
	DbHost               string    `arg:"env:DB_CARD_DELIVERY_HOST,required"`
	DbPort               string    `arg:"env:DB_CARD_DELIVERY_PORT,required"`
	DbName               string    `arg:"env:DB_CARD_DELIVERY_NAME,required"`
	DbUser               string    `arg:"env:DB_CARD_DELIVERY_USER,required"`
	DbPass               string    `arg:"env:DB_CARD_DELIVERY_PASS,required"`
	UsersEndpoint        string    `arg:"env:USERS_ENDPOINT,required"`
	DictEndpoint         string    `arg:"env:DICT_ENDPOINT,required"`
	DictReaderEndpoint   string    `arg:"env:FLEX_DICT_REST_ENDPOINT,required"`
	BusinessCardEndpoint string    `arg:"env:BUSINESS_CARD_ENDPOINT,required"`
	LogLevel             log.Level `arg:"env:LOG_LEVEL"`
	Port                 string    `arg:"env:PORT,required"`
	Hostname             string    `arg:"env:HOSTNAME,required"`
	FileStorageClient    string    `arg:"env:FILESTORAGE_ENDPOINT,required"`
	SalaryOrderClient    string    `arg:"env:SALARY_ORDER_ENDPOINT,required"`
	MailSenderQueue      string    `arg:"env:MAIL_SENDER_QUEUE"`
	RabbitMqHost         string    `arg:"env:RABBITMQ_HOST"`
	RabbitMqPort         string    `arg:"env:RABBITMQ_PORT"`
	RabbitMqUser         string    `arg:"env:RABBITMQ_USER"`
	RabbitMqPass         string    `arg:"env:RABBITMQ_PASS"`
}

var Props Args

func LoadConfig() {
	arg.Parse(&Props)
}
