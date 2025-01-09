package dipt

type PSPAccount struct {
	ID int `json:"id"`
	CreatedAt string `json:"created_at"`
	TokenAccess string `json:"token_access"`
	Name string `json:"name"`
	Type string `json:"type"`
	Address string `json:"address"`
	CIL string `json:"cil"`
	Country string `json:"country"`
	ActiveStatus bool   `json:"active_status"`
	RabbitMQHost string `json:"rabbitmqHost"`
	RabbitMQPort string `json:"rabbitmqPort"`
	RabbitMQUser string `json:"rabbitmqUser"`
	RabbitMQPass string `json:"rabbitmqPass"`
	BICSwift string `json:"bic_swift"`
	RabbitMQPmntExName string `json:"rabbitmqPmntExName"`
	RabbitMQNrmlPmntRtngKey string `json:"rabbitmqNrmlPmntRtngKey"`
}