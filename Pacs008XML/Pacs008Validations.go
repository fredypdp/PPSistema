package Pacs008XML

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fredypdp/PPSistema/Pacs002XML"
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

func validarValorMoeda(moeda string, valor string, numbTransacao int) error {
	// Tenta converter a string para float64
	num, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		return fmt.Errorf("erro ao converter valor para número decimal na transação %d: %v", numbTransacao, err)
	}

	// Verifica se o número é maior que zero
	if num <= 0 {
		return fmt.Errorf("o valor deve ser maior que zero")
	}

	// Abre o arquivo countries.json
	file, err := os.Open("constant/countries.json")
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo JSON na transação %d: %v", numbTransacao, err)
	}
	defer file.Close()

	// Decodifica o JSON em um slice de Country
	var countries []Country
	if err := json.NewDecoder(file).Decode(&countries); err != nil {
		return fmt.Errorf("erro ao decodificar o arquivo JSON na transação %d: %v", numbTransacao, err)
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
		return fmt.Errorf("a moeda '%s' não é válida, na transação %d", moeda, numbTransacao)
	}

	return nil
}

func validarIdentificadoresUnicos(document DocumentPacs008) error {
    instrIds := make(map[string]bool)
    endToEndIds := make(map[string]bool)

    for i, tx := range document.FIToFICstmrCdtTrf.CdtTrfTxInf {
        // Verifica InstrId
        if _, exists := instrIds[tx.PmtId.InstrId]; exists {
			errText := fmt.Errorf("InstrId duplicado encontrado na transação %d: %s", i+1, tx.PmtId.InstrId)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				document.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				document.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				tx.PmtId.EndToEndId,
				tx.PmtId.InstrId,
				"AM05",
				tx.CdtrAgt.FinInstnId.BICFI,
				tx.IntrBkSttlmAmt.Ccy,
				tx.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
        }
        instrIds[tx.PmtId.InstrId] = true

        // Verifica EndToEndId
        if _, exists := endToEndIds[tx.PmtId.EndToEndId]; exists {
			errText := fmt.Errorf("EndToEndId duplicado encontrado na transação %d: %s", i+1, tx.PmtId.EndToEndId)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				document.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				document.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				tx.PmtId.EndToEndId,
				tx.PmtId.InstrId,
				"AM05",
				tx.CdtrAgt.FinInstnId.BICFI,
				tx.IntrBkSttlmAmt.Ccy,
				tx.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
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


	if err := validarIdentificadoresUnicos(xmlPacs008); err != nil {
        return err
    }

	// Validação de elementos de cada CdtTrfTxInf
    for i, txInf := range xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf {
		numbTransacao := i+1
		// Validação dos elementos
		validacoesDeCampo := []struct {
			campo string
			valor string
		}{
			{fmt.Sprintf("Identificador único da mensagem na transação %d", numbTransacao), xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId},
			{fmt.Sprintf("Data e hora de criação da mensagem na transação %d", numbTransacao), xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm},
			{fmt.Sprintf("Número de transações contidas na mensagem na transação %d", numbTransacao), xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.NbOfTxs},
			{fmt.Sprintf("Método de liquidação na transação %d", numbTransacao), xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd},
			{fmt.Sprintf("Identificador único da instrução na transação %d", numbTransacao), txInf.PmtId.InstrId},
			{fmt.Sprintf("Identificador end-to-end na transação %d", numbTransacao), txInf.PmtId.EndToEndId},
			{fmt.Sprintf("Moeda do valor bruto do pagamento na transação %d", numbTransacao), txInf.InstdAmt.Ccy},
			{fmt.Sprintf("Valor bruto do pagamento na transação %d", numbTransacao), txInf.InstdAmt.Val},
			{fmt.Sprintf("Moeda do valor da liquidação na transação %d", numbTransacao), txInf.IntrBkSttlmAmt.Ccy},
			{fmt.Sprintf("Valor da liquidação na transação %d", numbTransacao), txInf.IntrBkSttlmAmt.Val},
			{fmt.Sprintf("Nome do remetente na transação %d", numbTransacao), txInf.Dbtr.Nm},
			{fmt.Sprintf("Código do país do remetente na transação %d", numbTransacao), txInf.Dbtr.PstlAdr.Ctry},
			{fmt.Sprintf("Moeda da conta do remetente na transação %d", numbTransacao), txInf.DbtrAcct.Ccy},
			{fmt.Sprintf("Identificador da conta do remetente na transação %d", numbTransacao), txInf.DbtrAcct.Id.Othr.Id},
			{fmt.Sprintf("Tipo do identificador da conta do remetente na transação %d", numbTransacao), txInf.DbtrAcct.Id.Othr.SchmeNm.Prtry},
			{fmt.Sprintf("Tipo da conta do remetente na transação %d", numbTransacao), txInf.DbtrAcct.Tp.Prtry},
			{fmt.Sprintf("Nome do banco do remetente na transação %d", numbTransacao), txInf.DbtrAgt.FinInstnId.Nm},
			{fmt.Sprintf("Código BIC do banco do remetente na transação %d", numbTransacao), txInf.DbtrAgt.FinInstnId.BICFI},	
			{fmt.Sprintf("Nome do beneficiário na transação %d", numbTransacao), txInf.Cdtr.Nm},
			{fmt.Sprintf("Código do país do beneficiário na transação %d", numbTransacao), txInf.Cdtr.PstlAdr.Ctry},
			{fmt.Sprintf("Moeda da conta do beneficiário na transação %d", numbTransacao), txInf.CdtrAcct.Ccy},
			{fmt.Sprintf("Identificador da conta do beneficiário na transação %d", numbTransacao), txInf.CdtrAcct.Id.Othr.Id},
			{fmt.Sprintf("Tipo do identificador da conta do beneficiário na transação %d", numbTransacao), txInf.CdtrAcct.Id.Othr.SchmeNm.Prtry},
			{fmt.Sprintf("Tipo da conta do beneficiário na transação %d", numbTransacao), txInf.CdtrAcct.Tp.Prtry},
			{fmt.Sprintf("Nome do banco do beneficiário na transação %d", numbTransacao), txInf.CdtrAgt.FinInstnId.Nm},
			{fmt.Sprintf("Código BIC do banco do beneficiário na transação %d", numbTransacao), txInf.CdtrAgt.FinInstnId.BICFI},
			{fmt.Sprintf("Nome do primeiro banco intermediário na transação %d", numbTransacao), txInf.IntrmyAgt1.FinInstnId.Nm},
			{fmt.Sprintf("BICFI do primeiro banco intermediário na transação %d", numbTransacao), txInf.IntrmyAgt1.FinInstnId.BICFI},
			{fmt.Sprintf("Nome do PP no primeiro banco intermediário na transação %d", numbTransacao), txInf.IntrmyAgt1Acct.Nm},
			{fmt.Sprintf("Páis do PP no primeiro banco intermediário na transação %d", numbTransacao), txInf.IntrmyAgt1Acct.Ccy},
			{fmt.Sprintf("ID do PP no primeiro banco intermediário na transação %d", numbTransacao), txInf.IntrmyAgt1Acct.Id.Othr.Id},
			{fmt.Sprintf("Responsável por pagar as taxas na transação %d", numbTransacao), txInf.ChrgBr},
			{fmt.Sprintf("Moeda do valor da taxa na transação %d", numbTransacao), txInf.ChrgsInf.Amt.Ccy},
			{fmt.Sprintf("Valor da taxa na transação %d", numbTransacao), txInf.ChrgsInf.Amt.Val},
			{fmt.Sprintf("Nome do banco que cobra as taxas na transação %d", numbTransacao), txInf.ChrgsInf.Agt.FinInstnId.Nm},
			{fmt.Sprintf("BICFI do banco que cobra as taxas na transação %d", numbTransacao), txInf.ChrgsInf.Agt.FinInstnId.BICFI},
		}

		// Validação de campos obrigatórios
		for _, v := range validacoesDeCampo {
			if err := validarCampo(v.campo, v.valor); err != nil {
				xml, errXml := Pacs002XML.CreateDocumentResponse(
					xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
					"pacs.008.001.08",
					xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
					"RJCT",
					"RJCT",
					txInf.PmtId.EndToEndId,
					txInf.PmtId.InstrId,
					"AM05",
					txInf.CdtrAgt.FinInstnId.BICFI,
					txInf.IntrBkSttlmAmt.Ccy,
					txInf.IntrBkSttlmAmt.Val,
					[]string{err.Error()},
				)
				if errXml != nil {
					log.Fatalf("Erro: %v", errXml)
				}
			
				fmt.Println(xml)
				return err
			}
		}

		if xmlPacs008.XMLNS != "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.08" {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{"namespace inválido"},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errors.New("namespace inválido")
		}
	
		_, err = strconv.Atoi(xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.NbOfTxs)
		if err != nil {
			errText := fmt.Errorf("erro ao converter o valor do elemento NbOfTxs para número na transação %d: %s", numbTransacao, err)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
		}
		
		if xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.SttlmInf.SttlmMtd != "INDA" {
			errText := fmt.Errorf("o elemento 'SttlmMtd' deve conter o valor 'INDA' na transação %d", numbTransacao)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
		}

		valorPagamentoTotal, _ := strconv.ParseFloat(txInf.InstdAmt.Val, 64)
        valorPagamentoLiquido, _ := strconv.ParseFloat(txInf.IntrBkSttlmAmt.Val, 64)
        valorPagamentoTaxa, _ := strconv.ParseFloat(txInf.ChrgsInf.Amt.Val, 64)
        valorDaTaxa := valorPagamentoTotal * constant.TaxaPagamentoNacional

		if valorPagamentoLiquido != valorPagamentoTotal - valorDaTaxa {
			errText := fmt.Errorf("valor líquido do pagamento está incorreto na transação %d", numbTransacao)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
		}

		if valorPagamentoTaxa != valorDaTaxa {
			errText := fmt.Errorf("valor da taxa está incorreto na transação %d", numbTransacao)
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{errText.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return errText
		}
	
		// Validar valores monetários
		if err := validarValorMoeda(txInf.IntrBkSttlmAmt.Ccy, txInf.IntrBkSttlmAmt.Val, numbTransacao); err != nil {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{err.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return err
		}
	
		if err := validarValorMoeda(txInf.InstdAmt.Ccy, txInf.InstdAmt.Val, numbTransacao); err != nil {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{err.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return err
		}

		if err := validarValorMoeda(txInf.ChrgsInf.Amt.Ccy, txInf.ChrgsInf.Amt.Val, numbTransacao); err != nil {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{err.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return err
		}
	
		// Validar tipos de identificador
		if err := validarTipoIdentificador(txInf.CdtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{err.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return err
		}
	
		if err := validarTipoIdentificador(txInf.DbtrAcct.Id.Othr.SchmeNm.Prtry); err != nil {
			xml, errXml := Pacs002XML.CreateDocumentResponse(
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.MsgId,
				"pacs.008.001.08",
				xmlPacs008.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
				"RJCT",
				"RJCT",
				txInf.PmtId.EndToEndId,
				txInf.PmtId.InstrId,
				"AM05",
				txInf.CdtrAgt.FinInstnId.BICFI,
				txInf.IntrBkSttlmAmt.Ccy,
				txInf.IntrBkSttlmAmt.Val,
				[]string{err.Error()},
			)
			if errXml != nil {
				log.Fatalf("Erro: %v", errXml)
			}

			fmt.Println(xml)
			return err
		}
	}

    return nil
}