package dipt

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/joho/godotenv"
)

func GetIdentificadorByTelefone(telefoneNumero string) (*Identificador, error) {
	// Verifique se o arquivo .env existe
	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env existe, então tente carregá-lo
		err := godotenv.Load()
		if err != nil {
			log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	apiKey := os.Getenv("SUPABASE_APIKEY")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIKEY não definida.")
	}
	
	apiUrl := os.Getenv("SUPABASE_APIURL")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIURL não definida.")
	}
	
	url := fmt.Sprintf("%s/rest/v1/idents_tel?ident=eq.%s&select=*", apiUrl, telefoneNumero)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var Identificadores []Identificador
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&Identificadores); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &Identificadores[0], nil
}

func GetIdentificadorByEmail(email string) (*Identificador, error) {
	// Verifique se o arquivo .env existe
	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env existe, então tente carregá-lo
		err := godotenv.Load()
		if err != nil {
			log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	apiKey := os.Getenv("SUPABASE_APIKEY")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIKEY não definida.")
	}
	
	apiUrl := os.Getenv("SUPABASE_APIURL")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIURL não definida.")
	}
	
	url := fmt.Sprintf("%s/rest/v1/idents_email?ident=eq.%s&select=*", apiUrl, email)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var Identificadores []Identificador
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&Identificadores); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &Identificadores[0], nil
}
func GetIdentificadorByCil(cil string) (*Identificador, error) {
	// Verifique se o arquivo .env existe
	if _, err := os.Stat(".env"); err == nil {
		// O arquivo .env existe, então tente carregá-lo
		err := godotenv.Load()
		if err != nil {
			log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v", err)
		}
	}

	apiKey := os.Getenv("SUPABASE_APIKEY")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIKEY não definida.")
	}
	
	apiUrl := os.Getenv("SUPABASE_APIURL")
	if apiKey == "" {
		log.Println("Aviso: variável SUPABASE_APIURL não definida.")
	}
	
	url := fmt.Sprintf("%s/rest/v1/idents_cil?ident=eq.%s&select=*", apiUrl, cil)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var Identificadores []Identificador
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&Identificadores); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &Identificadores[0], nil
}