package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fredypdp/PPSistema/Transactions"
	"github.com/fredypdp/PPSistema/constant"
	"github.com/fredypdp/PPSistema/dipt"
	"github.com/fredypdp/PPSistema/rabbitmq"
	"github.com/joho/godotenv"
)

// Criar filas e exchanges do sistema
func CreateRabbitmqPSPUsers() error {
	// Verifique se o arquivo .env existe
	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env existe, então tente carregá-lo
		err := godotenv.Load()
		if err != nil {
			return fmt.Errorf("aviso: Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	// Acesse as variáveis de ambiente
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitManagementPort := os.Getenv("RABBITMQ_MANAGEMENT_PORT")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")

	accounts, err := dipt.GetAllPSPAccounts()
	if err != nil {
		return fmt.Errorf("erro ao buscar todos os PSP: %v", err)
	}

	if len(accounts) == 0 {
		return fmt.Errorf("não foi possível encontrar o PSP do remetente: %v", err)
	}
	
	for _, psp := range accounts {
		if strings.Contains(psp.RabbitMQUser, " ") {
			return fmt.Errorf("o nome de usuário (%s) escolhido para o rabbitm contém espaços em branco", psp.RabbitMQUser)
		}

		err := rabbitmq.CreatePSPUsers(rabbitHost, rabbitManagementPort, rabbitUser, rabbitPass, psp.RabbitMQUser, psp.RabbitMQPass)
		if err != nil {
			return fmt.Errorf("erro ao criar usuário %s: %v", psp.RabbitMQUser, err)
		}
	}

	return nil
}

func CreateExchangesAndQueues() error {
	rabbitmqChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer rabbitmqChannel.Close()

	// Cria a exchange de pagamentos
	if err := rabbitmq.CreateDirectExchange(rabbitmqChannel, constant.PaymentExchangeName); err != nil {
		return fmt.Errorf("erro ao criar exchange: %v", err)
	}
	log.Printf("Exchanges de pagamento criadas")

	// Criando as filas de pagamento
	if err := rabbitmq.CreateQueues(rabbitmqChannel, constant.PaymentQueues, constant.PaymentExchangeName); err != nil {
		return fmt.Errorf("erro ao criar filas: %v", err)
	}
	log.Printf("Filas de pagamento criadas")

	// Cria a exchange de respostas
	if err := rabbitmq.CreateDirectExchange(rabbitmqChannel, constant.PaymentResponseExchangeName); err != nil {
		return fmt.Errorf("erro ao criar exchange: %v", err)
	}
	log.Printf("Exchanges de resposta criadas")

	// Criando as filas de respostas
	if err := rabbitmq.CreateQueues(rabbitmqChannel, constant.PaymentResponseQueues, constant.PaymentResponseExchangeName); err != nil {
		return fmt.Errorf("erro ao criar filas: %v", err)
	}
	log.Printf("Filas de resposta criadas")

	return nil
}

func main() {
	// if err := CreateRabbitmqPSPUsers(); err != nil {
	// 	log.Fatalf("%s", err)
	// }

	if err := CreateExchangesAndQueues(); err != nil {
		log.Fatalf("%s", err)
	}

	// Cria um WaitGroup para gerenciar as goroutines
	var wg sync.WaitGroup

	// Adiciona ao WaitGroup o número de goroutines que serão iniciadas
	wg.Add(3)

	// Goroutine para consumir pagamentos normais
	go func() {
		defer wg.Done()
		log.Println("Iniciando consumo de pagamentos normais")
		Transactions.ConsumeNormalPayments()
	}()

	// Goroutine para consumir relatórios de status
	go func() {
		defer wg.Done()
		log.Println("Iniciando consumo de relatórios de status")
		Transactions.ConsumeStatusReports()
	}()

	// Goroutine para consumir devoluções
	go func() {
		defer wg.Done()
		log.Println("Iniciando consumo de devoluções")
		Transactions.ConsumeRefunds()
	}()

	// Espera por todas as goroutines concluírem
	wg.Wait()
}