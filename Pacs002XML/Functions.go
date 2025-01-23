package Pacs002XML

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func CreateDocumentResponse(originalMsgId, originalMsgNmId, originalCreDtTm, groupStatus, transactionStatus, originalEndToEndId, originalTxId, TxReasonCode, InstdAgtBICFI, IntrBkSttlmAmtCcy, IntrBkSttlmAmtVal string, additionalInfo []string) (string, error) {
	doc := DocumentPacs002{
		XMLNS: "urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10",
		FIToFIPmtStsRpt: fiToFIPmtStsRpt{
			GrpHdr: grpHdr{
				MsgId:   fmt.Sprintf("RESP-%s", uuid.New().String()),
				CreDtTm: time.Now().Format("2006-01-02T15:04:05"),
			},
			OrgnlGrpInfAndSts: orgnlGrpInfAndSts{
				OrgnlMsgId:   originalMsgId,
				OrgnlMsgNmId: originalMsgNmId,
				OrgnlCreDtTm: originalCreDtTm,
				GrpSts:       groupStatus,
			},
			TxInfAndSts: []txInfAndSts{
				{
					OrgnlEndToEndId: originalEndToEndId,
					OrgnlTxId:       originalTxId,
					TxSts:           transactionStatus,
					StsRsnInf: []stsRsnInf{
						{
							Rsn: rsn{
								Cd: TxReasonCode,
							},
							AddtlInf: additionalInfo,
						},
					},
					InstdAgt: instdAgt{
						FinInstnId: finInstnId{
							BICFI: InstdAgtBICFI,
						},
					},
					OrgnlTxRef: orgnlTxRef{
						IntrBkSttlmAmt: intrBkSttlmAmt{
							Ccy: IntrBkSttlmAmtCcy,
							Val: IntrBkSttlmAmtVal,
						},
					},
				},
			},
		},
	}

	// Converter para XML
	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao gerar XML: %v", err)
	}

	// Adicionar declaração XML no início
	xmlOutput := xml.Header + string(xmlData)

	return xmlOutput, nil
}