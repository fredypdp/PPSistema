package constant

import "github.com/fredypdp/PPSistema/rabbitmq"

// Exchange names
var (
	PaymentExchangeName          = "payment_exchange"
	PaymentResponseExchangeName  = "payment_response_exchange"
)

// Payment queues
var PaymentQueues = []rabbitmq.QueueWithRoutingKey{
	{QueueName: "normal_payment", RoutingKey: "ppn_pay"},
	{QueueName: "parcel_payment", RoutingKey: "ppp_pay"},
	{QueueName: "garant_payment", RoutingKey: "ppg_pay"},
}

// Specific payment queues
var (
	NormalPaymentQueueAndKey = PaymentQueues[0]
	ParcelPaymentQueueAndKey = PaymentQueues[1]
	GarantPaymentQueueAndKey = PaymentQueues[2]
)

// Payment response queues
var PaymentResponseQueues = []rabbitmq.QueueWithRoutingKey{
	{QueueName: "payment_status_reports", RoutingKey: "status_reports"},
	{QueueName: "payment_refunds", RoutingKey: "refunds"},
}

// Specific payment response queue
var (
	PaymentStatusQueueAndKey = PaymentResponseQueues[0]
	PaymentRefundQueueAndKey = PaymentResponseQueues[1]
)