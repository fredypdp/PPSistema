package Transactions

import (
	"fmt"
	"io"
	"time"
	"strconv"
	"net/http"
	"encoding/json"

	// "github.com/fredypdp/PPSistema/constant"
	// "github.com/fredypdp/PPSistema/Pacs008XML"
)

// Estrutura para mapear a resposta da AwesomeAPI
type currencyRateResponse struct {
	Code string `json:"code"`
	Bid  string `json:"bid"`
}

func getExchangeRate(currency string) (float64, error) {
	// Endpoint da AwesomeAPI
	url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/USD-%s", currency)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao obter taxa de câmbio: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("falha na API: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta da API: %v", err)
	}

	var rates map[string]currencyRateResponse
	if err := json.Unmarshal(body, &rates); err != nil {
		return 0, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Ajuste da chave para refletir o formato correto
	key := fmt.Sprintf("USD%s", currency)
	rate, exists := rates[key]
	if !exists {
		return 0, fmt.Errorf("taxa para moeda %s não encontrada", currency)
	}

	// Converte a taxa de string para float
	return strconv.ParseFloat(rate.Bid, 64)
}

// func taxarPagamento(document Pacs008XML.DocumentPacs008) (Pacs008XML.DocumentPacs008, *Pacs008XML.DocumentPacs008, error) {
// 	// Converte o valor original de string para float
// 	originalAmt, err := strconv.ParseFloat(document.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Val, 64)
// 	if err != nil {
// 		return Pacs008XML.DocumentPacs008{}, nil, fmt.Errorf("valor inválido: %v", err)
// 	}

// 	// Obtém a moeda e busca a taxa de câmbio
// 	currency := document.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Ccy
// 	rate, err := getExchangeRate(currency)
// 	if err != nil {
// 		return Pacs008XML.DocumentPacs008{}, nil, err
// 	}

// 	// Converte o valor para dólares
// 	amountInUSD := originalAmt / rate
// 	if amountInUSD < constant.ValorMinimoParaTaxar {
// 		return document, nil, nil // Sem taxa, retorna o documento original e nenhum TaxaPacs008
// 	}

// 	// Calcula a taxa e o valor restante
// 	valorDaTaxa := originalAmt * constant.TaxaPagamentoNacional
// 	valorTaxado := originalAmt - valorDaTaxa

// 	// Cria uma cópia do documento original para o valor restante
// 	UserPacs008 := document
// 	UserPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Val = fmt.Sprintf("%.2f", valorTaxado)

// 	// Cria um novo documento para representar a taxa
// 	TaxaPacs008 := document
// 	TaxaPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Val = fmt.Sprintf("%.2f", valorDaTaxa)
// 	TaxaPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAgt.FinInstnId.BICFI = "BICFIBeneficiario"
// 	TaxaPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.Cdtr.Nm = "PP"
// 	TaxaPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAcct.Id.Othr.Id = "+244123456789" // ID da conta da empresa
// 	TaxaPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAcct.Id.Othr.SchmeNm.Prtry = "PHONE" // Tipo de identificador

// 	// Retorna os dois documentos
// 	return UserPacs008, &TaxaPacs008, nil
// }