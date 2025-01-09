package Pacs002XML

import (
	"encoding/xml"
	"errors"
	"fmt"
)

func validarCampo(nome, valor string) error {
	if valor == "" {
		return fmt.Errorf("%s indefinido ou vazio", nome)
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

	// Validação dos elementos
	validacoes := []struct {
		campo string
		valor string
	}{
		{"MsgId", xmlPacs002.FIToFIPmtStsRpt.GrpHdr.MsgId},
		{"CreDtTm", xmlPacs002.FIToFIPmtStsRpt.GrpHdr.CreDtTm},
		{"OrgnlMsgId", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlMsgId},
		{"OrgnlMsgNmId", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlMsgNmId},
		{"OrgnlCreDtTm", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.OrgnlCreDtTm},
		{"GrpSts", xmlPacs002.FIToFIPmtStsRpt.OrgnlGrpInfAndSts.GrpSts},
		{"OrgnlEndToEndId", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.OrgnlEndToEndId},
		{"OrgnlTxId", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.OrgnlTxId},
		{"TxSts", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.TxSts},
		{"BICFI (InstdAgt)", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.InstdAgt.FinInstnId.BICFI},
		{"IntrBkSttlmAmt (moeda)", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.OrgnlTxRef.IntrBkSttlmAmt.Ccy},
		{"IntrBkSttlmAmt (valor)", xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.OrgnlTxRef.IntrBkSttlmAmt.Val},
	}

	// Validação de campos obrigatórios
	for _, v := range validacoes {
		if err := validarCampo(v.campo, v.valor); err != nil {
			return err
		}
	}

	if xmlPacs002.FIToFIPmtStsRpt.TxInfAndSts.InstdAgt.FinInstnId.BICFI == "" {
		return fmt.Errorf("o elemento 'BICFI' não pode estar vazio para o InstdAgt")
	}

	return nil
}