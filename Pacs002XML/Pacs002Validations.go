package Pacs002XML

import (
	"encoding/xml"
	"errors"
	"fmt"
)

var validGrpStsCodes = []string{
    "ACTC", // AcceptedTechnicalValidation
    "ACCP", // AcceptedCustomerProfile
    "ACSP", // AcceptedSettlementInProcess
    "ACSC", // AcceptedSettlementCompleted
    "ACWC", // AcceptedWithChange
    "PDNG", // Pending
    "RCVD", // Received
    "RJCT", // Rejected
}

var validTxStsCodes = []string{
    "ACSP", // AcceptedSettlementInProcess
    "ACSC", // AcceptedSettlementCompleted
    "RJCT", // Rejected
    "PDNG", // Pending
    "RCVD", // Received
}

func validarCodigoGrpSts(grpSts string) error {
    for _, codigo := range validGrpStsCodes {
        if grpSts == codigo {
            return nil
        }
    }
    return fmt.Errorf("GrpSts inválido: %s", grpSts)
}

// Validação do TxSts
func validarCodigoTxSts(txSts string) error {
    for _, codigo := range validTxStsCodes {
        if txSts == codigo {
            return nil
        }
    }
    return fmt.Errorf("TxSts inválido: %s", txSts)
}

func validarCampo(nome, valor string) error {
	if valor == "" {
		return fmt.Errorf("%s indefinido ou vazio", nome)
	}
	return nil
}

func validarIdentificadoresUnicos(document DocumentPacs002) error {
    orgnlTxIds := make(map[string]bool)
    orgnlEndToEndIds := make(map[string]bool)

    for i, tx := range document.FIToFIPmtStsRpt.TxInfAndSts {
        // Verifica OrgnlTxId
        if _, exists := orgnlTxIds[tx.OrgnlTxId]; exists {
            return fmt.Errorf("OrgnlTxId duplicado encontrado na transação %d: %s", i+1, tx.OrgnlTxId)
        }
        orgnlTxIds[tx.OrgnlTxId] = true

        // Verifica OrgnlEndToEndId
        if _, exists := orgnlEndToEndIds[tx.OrgnlEndToEndId]; exists {
            return fmt.Errorf("OrgnlEndToEndId duplicado encontrado na transação %d: %s", i+1, tx.OrgnlEndToEndId)
        }
        orgnlEndToEndIds[tx.OrgnlEndToEndId] = true
    }

    return nil
}

func ValidarFormatoXML(data []byte) error {
	var xmlPacs002 DocumentPacs002
	err := xml.Unmarshal(data, &xmlPacs002)
	if err != nil {
		return fmt.Errorf("falha ao formatar XML Pacs.002: %v", err)
	}

	// Validação do namespace
	if xmlPacs002.XMLNS != "urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10" {
		return errors.New("namespace inválido")
	}

	// Validações de namespace e campos obrigatórios
	if err := validarCampo("Namespace", xmlPacs002.XMLNS); err != nil {
		return err
	}

	grpSts := xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.GrpSts
	if err := validarCodigoGrpSts(grpSts); err != nil {
		return err
	}

	// Validação de elementos do GrpHdr e OrgnlGrpInfAndSts
    if err := validarCampo("MsgId", xmlPacs002.FIToFIPmtStsRpt.GrpHdr.MsgId); err != nil {
        return err
    }

    if err := validarCampo("CreDtTm", xmlPacs002.FIToFIPmtStsRpt.GrpHdr.CreDtTm); err != nil {
        return err
    }

    if err := validarCampo("OrgnlMsgId", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlMsgId); err != nil {
        return err
    }

    if err := validarCampo("OrgnlMsgNmId", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlMsgNmId); err != nil {
        return err
    }

    if err := validarCampo("OrgnlCreDtTm", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlCreDtTm); err != nil {
        return err
    }

    if err := validarCampo("GrpSts", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.GrpSts); err != nil {
        return err
    }

	if err := validarIdentificadoresUnicos(xmlPacs002); err != nil {
        return err
    }

	// Validação de elementos de cada TxInfAndSts
    for i, txInfAndSts := range xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts {
		txSts := txInfAndSts.TxSts
		if err := validarCodigoTxSts(txSts); err != nil {
			return err
		}

		validacoes := []struct {
			campo string
			valor string
		}{
			{fmt.Sprintf("OrgnlEndToEndId  na transação %d", i+1), txInfAndSts.OrgnlEndToEndId},
			{fmt.Sprintf("OrgnlTxId  na transação %d", i+1), txInfAndSts.OrgnlTxId},
			{fmt.Sprintf("TxSts  na transação %d", i+1), txInfAndSts.TxSts},
			{fmt.Sprintf("BICFI (InstdAgt)  na transação %d", i+1), txInfAndSts.InstdAgt.FinInstnId.BICFI},
			{fmt.Sprintf("IntrBkSttlmAmt (moeda)  na transação %d", i+1), txInfAndSts.OrgnlTxRef.IntrBkSttlmAmt.Ccy},
			{fmt.Sprintf("IntrBkSttlmAmt (valor)  na transação %d", i+1), txInfAndSts.OrgnlTxRef.IntrBkSttlmAmt.Val},
		}
	
		// Validação de campos obrigatórios
		for _, v := range validacoes {
			if err := validarCampo(v.campo, v.valor); err != nil {
				return err
			}
		}
	
		// Regras dos códigos de de status
		switch grpSts {
		case "ACTC", "ACCP", "ACSP", "ACSC", "ACWC":
			if txSts == "RJCT" {
				return errors.New("TransactionStatus não pode ser 'RJCT' quando GroupStatus é um valor aceito")
			}
		case "PDNG":
			if txSts == "RJCT" {
				return errors.New("TransactionStatus não pode ser 'RJCT' quando GroupStatus é 'PDNG'")
			}
		case "RCVD":
			if txSts != "" {
				return errors.New("TransactionStatus deve estar ausente quando GroupStatus é 'RCVD'")
			}
		case "RJCT":
			if txSts != "RJCT" {
				return errors.New("TransactionStatus deve ser 'RJCT' quando GroupStatus é 'RJCT'")
			}
		}
	}

	return nil
}