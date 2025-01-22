package dipt

import "time"

type PSPAccount struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	TokenAccess string `json:"token_access"`
	Name string `json:"name"`
	Type string `json:"type"`
	Address string `json:"address"`
	CIL string `json:"cil"`
	Country string `json:"country"`
	ActiveStatus bool `json:"active_status"`
	RabbitMQHost string `json:"rabbitmqHost"`
	RabbitMQPort string `json:"rabbitmqPort"`
	RabbitMQUser string `json:"rabbitmqUser"`
	RabbitMQPass string `json:"rabbitmqPass"`
	BICSwift string `json:"bic_swift"`
	RabbitMQPmntRespExName string `json:"rabbitmqPmntRespExName"`
	RabbitMQStatusReportsRtngKey string `json:"rabbitmqStatusReportsRtngKey"`
	RabbitMQRefundsRtngKey string `json:"rabbitmqRefundsRtngKey"`
	RabbitMQPmntExName string `json:"rabbitmqPmntExName"`
	RabbitMQNrmlPmntRtngKey string `json:"rabbitmqNrmlPmntRtngKey"`
	RabbitMQPrclPmntRtngKey string `json:"rabbitmqPrclPmntRtngKey"`
	RabbitMQGrntPmntRtngKey string `json:"rabbitmqGrntPmntRtngKey"`
}

type Identificador struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Ident string `json:"ident"`
	IdentCountry string `json:"ident_country"`
	IdentTotalPaymentsSent *int `json:"ident_total_payments_sent"`
	IdentTotalPaymentsReceived *int `json:"ident_total_payments_received"`
	PspID int `json:"psp_id"`
	PspName string `json:"psp_name"`
	PspUserID string `json:"psp_user_id"`
	PspUserCil string `json:"psp_user_cil"`
	PspUserName string `json:"psp_user_name"`
	PspUserAddress string `json:"psp_user_address"`
	PspUserTypeAccount *string `json:"psp_user_type_account"`
}
