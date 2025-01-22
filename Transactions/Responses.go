package Transactions

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/fredypdp/PPSistema/Pacs002XML"
	"github.com/fredypdp/PPSistema/Pacs004XML"
	"github.com/fredypdp/PPSistema/constant"
	"github.com/fredypdp/PPSistema/rabbitmq"
	"github.com/joho/godotenv"
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
	if rabbitUser == "" {
		log.Println("Aviso: variável RABBITMQ_USER não definida.")
	}

	// Canal para receber as mensagens consumidas
	reportsChannel := make(chan amqp.Delivery)
	defer close(reportsChannel)
	go func() {
		if err := rabbitmq.ConsumeQueue(rabbitmqChannel, reportsChannel, constant.PaymentStatusQueueAndKey.QueueName, rabbitUser); err != nil {
			log.Printf("Erro ao consumir fila de reports: %v", err)
		}
	}()

	// Loop para processar mensagens continuamente
	for report := range reportsChannel {
		err := Pacs002XML.ValidarFormatoXML(report.Body)
		if err != nil {
			log.Printf("erro ao validar XML Pacs.002: %v", err)
			if ackErr := report.Ack(false); ackErr != nil {
				log.Printf("erro ao confirmar a mensagem: %v", ackErr)
			}
			continue
		}

		if ackErr := report.Ack(false); ackErr != nil {
			log.Printf("erro ao confirmar a mensagem: %v", ackErr)
		}

		// Processa a resposta
		if err := processarPacs002(report.Body); err != nil {
			log.Printf("Erro ao fazer o pagamento: %v", err)
		}
	}
}

func ConsumeRefunds() {
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
	if rabbitUser == "" {
		log.Println("Aviso: variável RABBITMQ_USER não definida.")
	}

	// Canal para receber as mensagens consumidas
	reportsChannel := make(chan amqp.Delivery)
	defer close(reportsChannel)
	go func() {
		if err := rabbitmq.ConsumeQueue(rabbitmqChannel, reportsChannel, constant.PaymentRefundQueueAndKey.QueueName, rabbitUser); err != nil {
			log.Printf("Erro ao consumir fila de reports: %v", err)
		}
	}()

	// Loop para processar mensagens continuamente
	for report := range reportsChannel {
		err := Pacs004XML.ValidarFormatoXML(report.Body)
		if err != nil {
			log.Printf("erro ao validar XML Pacs.004: %v", err)
			if ackErr := report.Ack(false); ackErr != nil {
				log.Printf("erro ao confirmar a mensagem: %v", ackErr)
			}
			continue
		}

		if ackErr := report.Ack(false); ackErr != nil {
			log.Printf("erro ao confirmar a mensagem: %v", ackErr)
		}

		log.Println("Mensagem pacs.004:")
		log.Println(string(report.Body))
	}
}

func processarPacs002(data []byte) error {
	var xmlPacs002 Pacs002XML.DocumentPacs002
	err := xml.Unmarshal(data, &xmlPacs002)
	if err != nil {
		return fmt.Errorf("falha ao formatar XML Pacs.008: %v", err)
	}

	err = xml.Unmarshal(data, &xmlPacs002)
	if err != nil {
		log.Printf("falha ao formatar XML Pacs.002: %v", err)
	}

	pacs002, err := xml.MarshalIndent(xmlPacs002, "", "   ")
	if err != nil {
		log.Printf("erro ao serializar XML do Pacs.002 restante: %v", err)
	}

	log.Println("Mensagem pacs.002:")
	log.Println(string(pacs002))

	// Validar o montante da liquidação interna em coroas suecas: <IntrBkSttlmAmt Ccy="AOA">71.12</IntrBkSttlmAmt>
	// Validar o status individual da transação - Aprovado (Accepted): <TxSts>ACSP</TxSts>
	// Pegar uma transação no DIPT usando o OrgnlEndToEndId do pacs.002 e Validar a concistência dos dados:
		// Identificador único da mensagem: <MsgId>ME180217000010000001230</MsgId>
		// Identificador da mensagem original: <OrgnlMsgId>M180217100120000001230</OrgnlMsgId>
		// Identificador ponta a ponta da transação original: <OrgnlEndToEndId>NO34347438583</OrgnlEndToEndId>
		// Código BIC da instituição financeira: <BICFI>TESTSW2X</BICFI>
	
	// Enviar os arquivos
	// chPSPUser, err := OpenPSPChannel("pp", "pp", pspBeneficiario.RabbitMQHost, pspBeneficiario.RabbitMQPort, pspBeneficiario.Name)
	// if err != nil {
	// 	return fmt.Errorf("%v", err)
	// }
	// defer chPSPUser.Close()

	// if err := sendPacs008ToCreditor(chPSPUser, string(pacs008Taxado), pspBeneficiario.RabbitMQPmntExName, pspBeneficiario.RabbitMQNrmlPmntRtngKey); err != nil {
	// 	log.Fatalf("Erro ao enviar pacs.008: %v", err)
	// }

	// if pacs008ReceberTaxa != nil {
	// 	fmt.Println("Para o PSP do PP:")
	// 	fmt.Println(string(pacs008ReceberTaxa))
	// }

	return nil
}

// Enviar o pacs008 para o banco do beneficiário
func sendPacsToDebtor(ch *amqp.Channel, body string, exchangeName string, Routingkey string) error {
	err := ch.Publish(
		exchangeName,
		Routingkey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/xml",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}