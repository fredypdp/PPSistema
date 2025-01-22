package dipt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetAllPSPAccounts() ([]PSPAccount, error) {
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
	
	
	url := fmt.Sprintf("%s/rest/v1/psp_accounts?select=*", apiUrl)

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

	var accounts []PSPAccount
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&accounts); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return accounts, nil
}


func GetPSPAccountByID(id int) (*PSPAccount, error) {
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
	
	url := fmt.Sprintf("%s/rest/v1/psp_accounts?id=eq.%s&select=*", apiUrl, strconv.Itoa(id))

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

	var accounts []PSPAccount
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&accounts); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &accounts[0], nil
}

func GetPSPAccountByBIC_Swift(bic_swift string) (*PSPAccount, error) {
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
	
	url := fmt.Sprintf("%s/rest/v1/psp_accounts?bic_swift=eq.%s&select=*", apiUrl, bic_swift)

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

	var accounts []PSPAccount
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&accounts); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta: %w", err)
	}

	return &accounts[0], nil
}