package Transactions

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/fredypdp/PPSistema/constant"
	"github.com/fredypdp/PPSistema/rabbitmq"
	"github.com/fredypdp/PPSistema/Pacs002XML"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeStatusReports() {
	rabbitmqChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer rabbitmqChannel.Close()

	// Verifique se o arquivo .env existe e carregue-o, se disponível
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	// Acesse as variáveis de ambiente
	rabbitUser := os.Getenv("RABBITMQ_USER")

	// Canal para receber as mensagens consumidas
	reportsChannel := make(chan amqp.Delivery)
	go rabbitmq.ConsumeQueue(rabbitmqChannel, reportsChannel, constant.PaymentStatusQueueAndKey.QueueName, rabbitUser)

	// Loop para processar mensagens continuamente
	for report := range reportsChannel {
		err := Pacs002XML.ValidarFormatoXML(report.Body)
		if err != nil {
			log.Printf("Erro ao validar XML Pacs.002: %v", err)
			if ackErr := report.Ack(false); ackErr != nil {
				log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
			}
			continue
		}

		if ackErr := report.Ack(false); ackErr != nil {
			log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
		}

		log.Println("Mensagem pacs.002:")
		log.Println(string(report.Body))
	}
}