package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	log.Fatalf("%s: %s", msg, err)
}

func OpenChannel() (*amqp.Channel, error) {
	// Verifique se o arquivo .env existe
    if _, err := os.Stat(".env"); err == nil {
        // O arquivo .env existe, então tente carregá-lo
        err := godotenv.Load()
        if err != nil {
            log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
        }
    }
	
	// Acesse as variáveis de ambiente
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")
	rabbitPort := os.Getenv("RABBITMQ_PORT")

	// Constrói a URL de conexão dinamicamente
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		failOnError(err, "Falha ao conectar-se com o RabbitMQ")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Falha ao abrir o canal")
		panic(err)
	}

	return ch, nil
}

func ConsumeQueue(ch *amqp.Channel, out chan<- amqp.Delivery, queueName string, consumerName string) error {
	msgs, err := ch.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publish(ch *amqp.Channel, body string, exchangeName string, Routingkey string) error {
	err := ch.Publish(
		exchangeName,
		Routingkey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Estrutura para permissões
type UserPublisherPermissions struct {
	Configure string `json:"configure"`
	Write     string `json:"write"`
	Read      string `json:"read"`
}

// Criar usuários dos PSP do DIPT
func CreatePSPUsers(rabbitHost, rabbitManagementPort, adminUser, adminPass, newUser, newPass string) error {
	// Criação de um cliente HTTP para enviar as requisições
	client := &http.Client{}

	// URL para criar o usuário
	url := fmt.Sprintf("http://%s:%s/api/users/%s", rabbitHost, rabbitManagementPort, newUser)

	// Definindo o payload para criar o usuário
	userPayload := map[string]interface{}{
		"password": newPass,
		"tags":     "none",
	}

	// Convertendo o payload para o formato JSON
	body, _ := json.Marshal(userPayload)

	// Criando a requisição HTTP para criar o usuário via a API do RabbitMQ
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Definindo a autenticação básica usando o usuário administrador
	req.SetBasicAuth(adminUser, adminPass)
	req.Header.Set("Content-Type", "application/json")

	// Enviando a requisição e obtendo a resposta
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Se o código de status for 409 (Conflito), significa que o usuário já existe
	if resp.StatusCode == http.StatusConflict {
		fmt.Printf("Usuário %s já existe, prosseguindo...\n", newUser)
		return nil
	}

	// Verificando o código de status da resposta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("falha ao criar usuário: %s", resp.Status)
	}

	// Configuração das permissões do usuário recém-criado
	perms := UserPublisherPermissions{
		Configure: "",
		Write:     ".*",
		Read:      "",
	}

	// Convertendo as permissões para o formato JSON
	permsBody, _ := json.Marshal(perms)

	// Construindo a URL para configurar as permissões do usuário
	permsURL := fmt.Sprintf("http://%s:%s/api/permissions/%s/%s", rabbitHost, rabbitManagementPort, "%2F", newUser)

	// Criando a requisição HTTP para configurar as permissões do usuário
	permsReq, err := http.NewRequest("PUT", permsURL, bytes.NewBuffer(permsBody))
	if err != nil {
		return err
	}

	// Definindo a autenticação básica novamente para configurar as permissões
	permsReq.SetBasicAuth(adminUser, adminPass)
	permsReq.Header.Set("Content-Type", "application/json")

	// Enviando a requisição para configurar as permissões do usuário
	permsResp, err := client.Do(permsReq)
	if err != nil {
		return err
	}
	defer permsResp.Body.Close()

	// Verificando se a configuração de permissões foi bem-sucedida
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("falha ao configurar permissões: %s", permsResp.Status)
	}

	// Se tudo ocorrer bem, retorna nil, indicando sucesso
	return nil
}

// Cria uma exchange do tipo "direct"
func CreateDirectExchange(ch *amqp.Channel, exchangeName string) error {
	err := ch.ExchangeDeclare(
		exchangeName,
		"direct", // Tipo da exchange
		true,     // Durável
		false,    // Auto-delete
		false,    // Interna
		false,    // Não espera confirmação do servidor
		nil,      // Argumentos extras
	)
	if err != nil {
		log.Printf("Erro ao criar exchange %s: %v", exchangeName, err)
		return err
	}
	
	return nil
}

// Cria as filas de pagamento e as conecta a uma exchange
func CreateQueues(ch *amqp.Channel, queues []QueueWithRoutingKey, exchangeName string) error {
	for _, queue := range queues {
		// Declara a fila
		_, err := ch.QueueDeclare(
			queue.QueueName,
			true,  // Durável
			false, // Auto-delete
			false, // Exclusiva
			false, // Não espera confirmação do servidor
			nil,   // Argumentos extras
		)
		if err != nil {
			log.Printf("Erro ao declarar fila %s: %v", queue.QueueName, err)
			return err
		}

		// Conecta a fila à exchange
		err = ch.QueueBind(
			queue.QueueName,
			queue.RoutingKey,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Erro ao conectar fila %s à exchange %s com routing key %s: %v",
				queue.QueueName, exchangeName, queue.RoutingKey, err)
			return err
		}
	}
	
	return nil
}