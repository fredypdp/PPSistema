package rabbitmq

type QueueWithRoutingKey struct {
	QueueName  string
	RoutingKey string
}