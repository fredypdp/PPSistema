package Transactions

import (
    "encoding/xml"
    "fmt"
    "log"
    "os"

    "github.com/fredypdp/PPSistema/Pacs008XML"
    "github.com/fredypdp/PPSistema/Pacs002XML"
    "github.com/fredypdp/PPSistema/constant"
    "github.com/fredypdp/PPSistema/dipt"
    "github.com/fredypdp/PPSistema/rabbitmq"
    "github.com/joho/godotenv"
    amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeNormalPayments() {
    rabbitmqChannel, err := rabbitmq.OpenChannel()
    if err != nil {
        panic(err)
    }
    defer rabbitmqChannel.Close()

    if _, err := os.Stat(".env"); err == nil {
        if err := godotenv.Load(); err != nil {
            log.Printf("Não foi possível carregar o arquivo .env: %v", err)
        }
    }

    rabbitUser := os.Getenv("RABBITMQ_USER")
    if rabbitUser == "" {
        log.Println("Aviso: variável RABBITMQ_USER não definida.")
    }

    pagamentosChannel := make(chan amqp.Delivery)
    defer close(pagamentosChannel)
    
    go func() {
        if err := rabbitmq.ConsumeQueue(rabbitmqChannel, pagamentosChannel, constant.NormalPaymentQueueAndKey.QueueName, rabbitUser); err != nil {
            log.Printf("Erro ao consumir fila de reports: %v", err)
        }
    }()

    for pagamento := range pagamentosChannel {
        err := Pacs008XML.ValidarFormatoXML(pagamento.Body)
        if err != nil {
            log.Printf("Erro ao validar XML Pacs.008: %v", err)
            if ackErr := pagamento.Ack(false); ackErr != nil {
                log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
            }
            continue
        }

        if ackErr := pagamento.Ack(false); ackErr != nil {
            log.Printf("Erro ao confirmar a mensagem: %v", ackErr)
        }

        if err := processarPacs008(pagamento.Body); err != nil {
            log.Printf("Erro ao fazer o pagamento: %v", err)
        }
    }
}

// Função auxiliar para calcular valores com base nas taxas
func calcularValoresComTaxas(valorBruto float64, taxaPercentual float64) (valorLiquido float64, valorTaxa float64) {
    valorTaxa = valorBruto * (taxaPercentual / 100)
    valorLiquido = valorBruto - valorTaxa
    return valorLiquido, valorTaxa
}

func processarPacs008(dataXML []byte) error {
    var xmlPacs008 Pacs008XML.DocumentPacs008
    err := xml.Unmarshal(dataXML, &xmlPacs008)
    if err != nil {
        return fmt.Errorf("falha ao formatar XML Pacs.008: %v", err)
    }

    for i, txInf := range xmlPacs008.FIToFICstmrCdtTrf.CdtTrfTxInf {
        pspPagador, err := dipt.GetPSPAccountByBIC_Swift(txInf.CdtrAgt.FinInstnId.BICFI)
        if err != nil {
            errText := fmt.Errorf("erro ao encontrar PSP do remetente %v na transação %d", err, i+1)
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
    
        if pspPagador == nil {
            errText := fmt.Errorf("não foi possível encontrar o PSP do remetente %v na transação %d", err, i+1)
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
    
        pspBeneficiario, err := dipt.GetPSPAccountByBIC_Swift(txInf.DbtrAgt.FinInstnId.BICFI)
        if err != nil {
            errText := fmt.Errorf("erro ao encontrar PSP do beneficiário %v na transação %d", err, i+1)
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
    
        if pspBeneficiario == nil {
            errText := fmt.Errorf("não foi possível encontrar o PSP do beneficiário %v na transação %d", err, i+1)
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
    
        log.Printf("transação %d", i+1)
        log.Println(string(dataXML))

        // chPSPIntermediario, err := OpenPSPChannel("pp", "pp", "host", "5672", txInf.IntrmyAgt1.FinInstnId.Nm)
        // if err != nil {
        //     return fmt.Errorf("erro ao conectar-se com intermediário: %v", err)
        // }
        // defer chPSPIntermediario.Close()
    
        // if err := sendPacs008ToIntermediary(chPSPIntermediario, string(dataXML), "payment_exchange", "normal_payment"); err != nil {
        //     return fmt.Errorf("erro ao enviar documento ao intermediário (Tx %d): %v", i+1, err)
        // }
    }

    return nil
}

func sendPacs008ToIntermediary(ch *amqp.Channel, body string, exchangeName string, routingKey string) error {
    err := ch.Publish(
        exchangeName,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/xml",
            Body:        []byte(body),
        },
    )
    if err != nil {
        return err
    }
    return nil
}