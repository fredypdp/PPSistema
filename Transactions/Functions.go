package Transactions

import (
	"fmt"
	"log"
	"os"

	"github.com/fredypdp/PPSistema/dipt"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) error {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func OpenPSPChannel(rabbitUser, rabbitPass, rabbitHost, rabbitPort, pspName string) (*amqp.Channel, error) {
	// Verifica se o arquivo .env existe
	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env existe, então tente carregá-lo
		if err := godotenv.Load(); err != nil {
			log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	// Constrói a URL de conexão dinamicamente
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitUser, rabbitPass, rabbitHost, rabbitPort)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		text := fmt.Sprintf("Falha ao conectar-se com o RabbitMQ do PSP '%s'", pspName)
		return nil, fmt.Errorf("%s: %w", text, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		text := fmt.Sprintf("Falha ao abrir o canal do PSP '%s'", pspName)
		return nil, fmt.Errorf("%s: %w", text, err)
	}

	return ch, nil
}

func pegarIdentificadorPorTipo(ident string, tipo string) (*dipt.Identificador, error) {
	var identificador *dipt.Identificador
	var err error

	switch tipo {
	case "PHONE":
		identificador, err = dipt.GetIdentificadorByTelefone(ident)
	case "EMAIL":
		identificador, err = dipt.GetIdentificadorByEmail(ident)
	case "CIL":
		identificador, err = dipt.GetIdentificadorByCil(ident)
	default:
		return nil, fmt.Errorf("tipo de identificador não suportado: %s", tipo)
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar PSP do remetente: %v", err)
	}

	if identificador == nil {
		return nil, fmt.Errorf("não foi possível encontrar o PSP do remetente")
	}

	return identificador, nil
}