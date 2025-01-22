package Pacs004XML

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

func validarTipoIdentificador(schemeType string) error {
    tiposValidos := map[string]bool{
        "PHONE": true,
        "EMAIL": true,
        "CIL":   true,
    }
    
    if !tiposValidos[schemeType] {
        return fmt.Errorf("tipo de identificador inválido: %s. Deve ser PHONE, EMAIL ou CIL", schemeType)
    }
    
    return nil
}

func ValidarFormatoXML(data []byte) error {
	var xmlPacs004 DocumentPacs004
	err := xml.Unmarshal(data, &xmlPacs004)
	if err != nil {
		return fmt.Errorf("falha ao formatar XML Pacs.004: %v", err)
	}

	// Validação do namespace
	if xmlPacs004.XMLNS != "urn:iso:std:iso:20022:tech:xsd:pacs.004.001.09" {
		return errors.New("namespace inválido")
	}

	// Validação dos elementos
	validacoes := []struct {
		campo string
		valor string
	}{
		{"MsgId", xmlPacs004.PmtRtr.GrpHdr.MsgId},
		{"CreDtTm", xmlPacs004.PmtRtr.GrpHdr.CreDtTm},
		{"NbOfTxs", xmlPacs004.PmtRtr.GrpHdr.NbOfTxs},
		{"TtlRtrdIntrBkSttlmAmt (moeda)", xmlPacs004.PmtRtr.GrpHdr.TtlRtrdIntrBkSttlmAmt.Ccy},
		{"TtlRtrdIntrBkSttlmAmt (valor)", xmlPacs004.PmtRtr.GrpHdr.TtlRtrdIntrBkSttlmAmt.Val},
		{"IntrBkSttlmDt", xmlPacs004.PmtRtr.GrpHdr.IntrBkSttlmDt},
		{"SttlmMtd", xmlPacs004.PmtRtr.GrpHdr.SttlmInf.SttlmMtd},
		{"RtrId", xmlPacs004.PmtRtr.TxInf.RtrId},
		{"OrgnlMsgId", xmlPacs004.PmtRtr.TxInf.OrgnlGrpInf.OrgnlMsgId},
		{"OrgnlMsgNmId", xmlPacs004.PmtRtr.TxInf.OrgnlGrpInf.OrgnlMsgNmId},
		{"OrgnlInstrId", xmlPacs004.PmtRtr.TxInf.OrgnlInstrId},
		{"OrgnlEndToEndId", xmlPacs004.PmtRtr.TxInf.OrgnlEndToEndId},
		{"OrgnlTxId", xmlPacs004.PmtRtr.TxInf.OrgnlTxId},
		{"RtrdIntrBkSttlmAmt (moeda)", xmlPacs004.PmtRtr.TxInf.RtrdIntrBkSttlmAmt.Ccy},
		{"RtrdIntrBkSttlmAmt (valor)", xmlPacs004.PmtRtr.TxInf.RtrdIntrBkSttlmAmt.Val},
		{"ChrgBr", xmlPacs004.PmtRtr.TxInf.ChrgBr},
		{"Rsn Cd", xmlPacs004.PmtRtr.TxInf.RtrRsnInf.Rsn.Cd},
		{"BICFI (Devedor)", xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.DbtrAgt.FinInstnId.BICFI},
		{"BICFI (Credor)", xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.CdtrAgt.FinInstnId.BICFI},
		{"Nm (Devedor)", xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.Dbtr.Nm},
		{"Nm (Credor)", xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.Cdtr.Nm},
	}

	// Validação de campos obrigatórios
	for _, v := range validacoes {
		if err := validarCampo(v.campo, v.valor); err != nil {
			return err
		}
	}

	// Validação do método de liquidação
	if xmlPacs004.PmtRtr.GrpHdr.SttlmInf.SttlmMtd != "CLRG" {
		return fmt.Errorf("o elemento 'SttlmMtd' deve conter o valor 'CLRG'")
	}

	// Validação dos valores monetários
	if err := validarValorMonetario(xmlPacs004.PmtRtr.GrpHdr.TtlRtrdIntrBkSttlmAmt); err != nil {
		return fmt.Errorf("erro na validação do TtlRtrdIntrBkSttlmAmt: %v", err)
	}

	if err := validarValorMonetario(xmlPacs004.PmtRtr.TxInf.RtrdIntrBkSttlmAmt); err != nil {
		return fmt.Errorf("erro na validação do RtrdIntrBkSttlmAmt: %v", err)
	}

	// Validação da mensagem original
	if xmlPacs004.PmtRtr.TxInf.OrgnlGrpInf.OrgnlMsgNmId != "pacs.008.001.08" {
		return fmt.Errorf("a mensagem original deve ser do tipo pacs.008.001.08")
	}

	// Adicionar validação dos tipos de identificador
	if err := validarTipoIdentificador(xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.CdtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
		return err
	}
	
	if err := validarTipoIdentificador(xmlPacs004.PmtRtr.TxInf.OrgnlTxRef.DbtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
		return err
	}

	return nil
}

type Country struct {
	Name     string `json:"Name"`
	Code     string `json:"Code"`
	Currency string `json:"Currency"`
}

func validarValorMonetario(amount amount) error {
	// Tenta converter a string para float64
	num, err := strconv.ParseFloat(amount.Val, 64)
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
		if country.Currency == amount.Ccy {
			validCurrency = true
			break
		}
	}

	if !validCurrency {
		return fmt.Errorf("a moeda '%s' não é válida", amount.Ccy)
	}

	return nil
}