package Pacs008XML

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/fredypdp/PPSistema/constant"
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

type Country struct {
	Name     string `json:"Name"`
	Code     string `json:"Code"`
	Currency string `json:"Currency"`
}

func validarValorMoeda(moeda string, valor string) error {
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

func validarIdentificadoresUnicos(document DocumentPacs008) error {
    instrIds := make(map[string]bool)
    endToEndIds := make(map[string]bool)

    for i, tx := range document.FIToFICstmrCdtTrf.CdtTrfTxInf {
        // Verifica InstrId
        if _, exists := instrIds[tx.PmtId.InstrId]; exists {
            return fmt.Errorf("InstrId duplicado encontrado na transação %d: %s", i+1, tx.PmtId.InstrId)
        }
        instrIds[tx.PmtId.InstrId] = true

        // Verifica EndToEndId
        if _, exists := endToEndIds[tx.PmtId.EndToEndId]; exists {
            return fmt.Errorf("EndToEndId duplicado encontrado na transação %d: %s", i+1, tx.PmtId.EndToEndId)
        }
        endToEndIds[tx.PmtId.EndToEndId] = true
    }

    return nil
}

func ValidarFormatoXML(data []byte) error {
    var xmlPacs008 DocumentPacs008
    err := xml.Unmarshal(data, &xmlPacs008)
    if err != nil {
        return fmt.Errorf("falha ao formatar XML Pacs.008: %v", err)
    }

    // Validação do namespace
    if xmlPacs008.XMLNS != "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08" {
        return errors.New("namespace inválido")
    }

	// Validação de elementos do GrpHdr
    if err := validarCampo("MsgId", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId); err != nil {
        return err
    }
    if err := validarCampo("CreDtTm", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm); err != nil {
        return err
    }
    if err := validarCampo("NbOfTxs", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.NbOfTxs); err != nil {
        return err
    }

	_, err = strconv.Atoi(xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.NbOfTxs)
	if err != nil {
		return fmt.Errorf("erro ao converter o valor do elemento NbOfTxs para número: %s", err)
	}

    if err := validarCampo("SttlmMtd", xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd); err != nil {
        return err
    }

	if err := validarIdentificadoresUnicos(xmlPacs008); err != nil {
        return err
    }

	// Validação de elementos de cada CdtTrfTxInf
    for i, txInf := range xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf {
		// Validação dos elementos
		validacoes := []struct {
			campo string
			valor string
		}{
			{fmt.Sprintf("InstrId na transação %d", i+1), txInf.PmtId.InstrId},
			{fmt.Sprintf("EndToEndId na transação %d", i+1), txInf.PmtId.EndToEndId},
			{fmt.Sprintf("InstdAmt (moeda) na transação %d", i+1), txInf.InstdAmt.Ccy},
			{fmt.Sprintf("InstdAmt (valor) na transação %d", i+1), txInf.InstdAmt.Val},
			{fmt.Sprintf("IntrBkSttlmAmt (moeda) na transação %d", i+1), txInf.IntrBkSttlmAmt.Ccy},
			{fmt.Sprintf("IntrBkSttlmAmt (valor) na transação %d", i+1), txInf.IntrBkSttlmAmt.Val},
			{fmt.Sprintf("BICFI (Credor) na transação %d", i+1), txInf.CdtrAgt.FinInstnId.BICFI},
			{fmt.Sprintf("Nm (Credor) na transação %d", i+1), txInf.Cdtr.Nm},
			{fmt.Sprintf("Ctry (Credor) na transação %d", i+1), txInf.Cdtr.PstlAdr.Ctry},
			{fmt.Sprintf("Id (Credor) na transação %d", i+1), txInf.CdtrAcct.Id.Othr.Id},
			{fmt.Sprintf("BICFI (Devedor) na transação %d", i+1), txInf.DbtrAgt.FinInstnId.BICFI},
			{fmt.Sprintf("Nm (Devedor) na transação %d", i+1), txInf.Dbtr.Nm},
			{fmt.Sprintf("Ctry (Devedor) na transação %d", i+1), txInf.Dbtr.PstlAdr.Ctry},
			{fmt.Sprintf("Id (Devedor) na transação %d", i+1), txInf.DbtrAcct.Id.Othr.Id},
			{fmt.Sprintf("ChrgBr na transação %d", i+1), txInf.ChrgBr},
			{fmt.Sprintf("ChrgsInf (moeda) na transação %d", i+1), txInf.ChrgsInf.Amt.Ccy},
			{fmt.Sprintf("ChrgsInf (valor) na transação %d", i+1), txInf.ChrgsInf.Amt.Val},
			{fmt.Sprintf("ChrgsInf na transação %d", i+1), txInf.ChrgsInf.Agt.FinInstnId.Nm},
			{fmt.Sprintf("ChrgsInf na transação %d", i+1), txInf.ChrgsInf.Agt.FinInstnId.BICFI},
			{fmt.Sprintf("IntrmyAgt1 BICFI na transação %d", i+1), txInf.IntrmyAgt1.FinInstnId.BICFI},
			{fmt.Sprintf("IntrmyAgt1Acct Id na transação %d", i+1), txInf.IntrmyAgt1Acct.Id.Othr.Id},
		}
	
		// Validação de campos obrigatórios
		for _, v := range validacoes {
			if err := validarCampo(v.campo, v.valor); err != nil {
				return err
			}
		}
	
		// Validações adicionais
		if xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd != "INDA" {
			return fmt.Errorf("o elemento 'SttlmMtd' deve conter o valor 'INDA'")
		}

		valorPagamentoTotal, _ := strconv.ParseFloat(txInf.InstdAmt.Val, 64)
        valorPagamentoLiquido, _ := strconv.ParseFloat(txInf.IntrBkSttlmAmt.Val, 64)
        valorPagamentoTaxa, _ := strconv.ParseFloat(txInf.ChrgsInf.Amt.Val, 64)
        valorDaTaxa := valorPagamentoTotal * constant.TaxaPagamentoNacional

		if valorPagamentoLiquido != valorPagamentoTotal - valorDaTaxa {
			text := fmt.Sprintf("valor líquido do pagamento está incorreto na transação %d", i+1)
			return fmt.Errorf("%s, no elemento IntrBkSttlmAmt", text)
		}

		if valorPagamentoTaxa != valorDaTaxa {
			text := fmt.Sprintf("valor da taxa está incorreto na transação %d", i+1)
			return fmt.Errorf("%s, no elemento ChrgsInf>Amt", text)
		}
	
		// Validar valores monetários
		if err := validarValorMoeda(txInf.IntrBkSttlmAmt.Ccy, txInf.IntrBkSttlmAmt.Val); err != nil {
			return err
		}
	
		if err := validarValorMoeda(txInf.InstdAmt.Ccy, txInf.InstdAmt.Val); err != nil {
			return err
		}

		if err := validarValorMoeda(txInf.ChrgsInf.Amt.Ccy, txInf.ChrgsInf.Amt.Val); err != nil {
			return err
		}
	
		// Validar tipos de identificador
		if err := validarTipoIdentificador(txInf.CdtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
			return err
		}
	
		if err := validarTipoIdentificador(txInf.DbtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
			return err
		}
	}

    return nil
}