package Transactions

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/fredypdp/PPSistema/Pacs008XML"
	"github.com/fredypdp/PPSistema/constant"
	"github.com/fredypdp/PPSistema/dipt"
	"github.com/fredypdp/PPSistema/rabbitmq"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	log.Fatalf("%s: %s", msg, err)
}

func OpenPSPChannel(rabbitUser string, rabbitPass string, rabbitHost string, rabbitPort string, pspName string) (*amqp.Channel, error) {
	// Verifique se o arquivo .env existe
    if _, err := os.Stat(".env"); err == nil {
        // O arquivo .env existe, então tente carregá-lo
        err := godotenv.Load()
        if err != nil {
            log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
        }
    }

	// Constrói a URL de conexão dinamicamente
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		text := fmt.Sprintf("Falha ao conectar-se com o RabbitMQ do PSP '%s'", pspName)
    	failOnError(err, text)
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		text := fmt.Sprintf("Falha ao abrir o canal do PSP '%s'", pspName)
		failOnError(err, text)
		panic(err)
	}

	return ch, nil
}

func ConsumeNormalPayments() {
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
	pagamentosChannel := make(chan amqp.Delivery)
	go rabbitmq.ConsumeQueue(rabbitmqChannel, pagamentosChannel, constant.NormalPaymentQueueAndKey.QueueName, rabbitUser)

	// Loop para processar mensagens continuamente
	for pagamento := range pagamentosChannel {
		err := Pacs008XML.ValidarFormatoXML(pagamento.Body)
		if err != nil {
			log.Printf("Erro ao validar XML Pacs.008: %v", err)
			if ackErr := pagamento.Ack(false); ackErr != nil {
				log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
			}
			continue
		}

		if ackErr := pagamento.Ack(false); ackErr != nil {
			log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
		}

		// Processa o pagamento
		if err := processarPacs008(pagamento.Body); err != nil {
			log.Printf("Erro ao fazer o pagamento: %v", err)
		}
	}
}

// Aplicar taxas ao pagamento e enviar para o beneficiário
func processarPacs008(data []byte) error {
	var xmlPacs008 Pacs008XML.Document
	err := xml.Unmarshal(data, &xmlPacs008)
	if err != nil {
		return fmt.Errorf("falha ao formatar XML Pacs.008: %v", err)
	}

	UserPacs008, TaxaPacs008, err := taxarPagamento(xmlPacs008)
	if err != nil {
		return fmt.Errorf("erro ao aplicar taxa: %v", err)
	}

	// Serializa os documentos para XML
	pacs008Taxado, err := xml.MarshalIndent(UserPacs008, "", "   ")
	if err != nil {
		return fmt.Errorf("erro ao serializar XML restante: %v", err)
	}
	pacs008ReceberTaxa, err := xml.MarshalIndent(TaxaPacs008, "", "   ")
	if err != nil {
		return fmt.Errorf("erro ao serializar XML da taxa: %v", err)
	}

	// Encontrar o PSP do pagador
	accountsPagador, err := dipt.GetPSPAccountByBIC_Swift(UserPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.DbtrAgt.FinInstnId.BICFI)
	if err != nil {
		return fmt.Errorf("erro ao buscar PSP do remetente: %v", err)
	}

	if len(accountsPagador) == 0 {
		return fmt.Errorf("não foi possível encontrar o PSP do remetente")
	}

	// Encontrar o PSP do beneficiário
	accountsBeneficiario, err := dipt.GetPSPAccountByBIC_Swift(UserPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAgt.FinInstnId.BICFI)
	if err != nil {
		return fmt.Errorf("erro ao buscar PSP do beneficiário: %v", err)
	}

	if len(accountsBeneficiario) == 0 {
		return fmt.Errorf("não foi possível encontrar o PSP do beneficiário")
	}

	pspBeneficiario := accountsBeneficiario[0]

	ch, err := OpenPSPChannel("pp", "pp", pspBeneficiario.RabbitMQHost, pspBeneficiario.RabbitMQPort, pspBeneficiario.Name)
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	if err := sendPacs008ToCreditor(ch, string(pacs008Taxado), pspBeneficiario.RabbitMQPmntExName, pspBeneficiario.RabbitMQNrmlPmntRtngKey); err != nil {
		log.Fatalf("Erro ao enviar pacs.008: %v", err)
	}

	fmt.Println("Para o PSP do cliente:")
	fmt.Println(string(pacs008Taxado))

	if pacs008ReceberTaxa != nil {
		fmt.Println("Para o banco do PP:")
		fmt.Println(string(pacs008ReceberTaxa))
	}

	return nil
}

// Enviar o pacs008 para o banco do beneficiário
func sendPacs008ToCreditor(ch *amqp.Channel, body string, exchangeName string, Routingkey string) error {
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