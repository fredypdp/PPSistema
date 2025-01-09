package Pacs008XML

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func validarCampo(nome, valor string) error {
	if valor == "" {
		return fmt.Errorf("%s indefinido ou vazio", nome)
	}
	return nil
}

func ValidarFormatoXML(data []byte) error {
	var xmlPacs008 Document
	err := xml.Unmarshal(data, &xmlPacs008)
	if err != nil {
		return fmt.Errorf("falha ao formatar XML Pacs.008: %v", err)
	}

	// Validação do namespace
	if xmlPacs008.XMLNS != "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08" {
		return errors.New("namespace inválido")
	}

	// Validação dos elementos
	validacoes := []struct {
		campo string
		valor string
	}{
		{"MsgId", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId},
		{"CreDtTm", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm},
		{"NbOfTxs", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.NbOfTxs},
		{"SttlmMtd", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd},
		{"InstrId", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.PmtId.InstrId},
		{"EndToEndId", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.PmtId.EndToEndId},
		{"IntrBkSttlmAmt (moeda)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Ccy},
		{"IntrBkSttlmAmt (valor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Val},
		{"BICFI (Credor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAgt.FinInstnId.BICFI},
		{"Nm (Credor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.Cdtr.Nm},
		{"Ctry (Credor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.Cdtr.PstlAdr.Ctry},
		{"Id (Credor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.CdtrAcct.Id.Othr.Id},
		{"BICFI (Devedor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.DbtrAgt.FinInstnId.BICFI},
		{"Nm (Devedor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.Dbtr.Nm},
		{"Ctry (Devedor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.Dbtr.PstlAdr.Ctry},
		{"Id (Devedor)", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.DbtrAcct.Id.Othr.Id},
		{"ChrgBr", xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.ChrgBr},
	}

	// Validação de campos obrigatórios
	for _, v := range validacoes {
		if err := validarCampo(v.campo, v.valor); err != nil {
			return err
		}
	}

	if xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd != "INDA" {
		return fmt.Errorf("o elemento 'SttlmMtd' deve conter o valor 'INDA'")
	}

	if err := validarIntrBkSttlmAmt(xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Ccy, xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf.IntrBkSttlmAmt.Val); err != nil {
		return err
	}

	return nil
}

type Country struct {
	Name     string `json:"Name"`
	Code     string `json:"Code"`
	Currency string `json:"Currency"`
}

func validarIntrBkSttlmAmt(moeda string, valor string) error {
	// Tenta converter a string para float64
	num, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		return fmt.Errorf("erro ao converter valor para número decimal: %v", err)
	}

	// Verifica se o número é maior que zero
	if num <= 0 {
		return fmt.Errorf("o valor deve ser maior que zero")
	}

	// Abre o arquivo countries.json
	file, err := os.Open("constant/countries.json")
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo JSON: %v", err)
	}
	defer file.Close()

	// Decodifica o JSON em um slice de Country
	var countries []Country
	if err := json.NewDecoder(file).Decode(&countries); err != nil {
		return fmt.Errorf("erro ao decodificar o arquivo JSON: %v", err)
	}

	// Verifica se a moeda está na lista de currencies
	validCurrency := false
	for _, country := range countries {
		if country.Currency == moeda {
			validCurrency = true
			break
		}
	}

	if !validCurrency {
		return fmt.Errorf("a moeda '%s' não é válida", moeda)
	}

	return nil
}