package Pacs004XML

import (
	"encoding/xml"
)

type DocumentPacs004 struct {
	XMLName xml.Name `xml:"Document"`
	XMLNS   string   `xml:"xmlns,attr"`
	PmtRtr  pmtRtr   `xml:"PmtRtr"`
}

type pmtRtr struct {
	GrpHdr grpHdr `xml:"GrpHdr"`
	TxInf  txInf  `xml:"TxInf"`
}

type grpHdr struct {
	MsgId                  string   `xml:"MsgId"`
	CreDtTm                string   `xml:"CreDtTm"`
	NbOfTxs                string   `xml:"NbOfTxs"`
	TtlRtrdIntrBkSttlmAmt  amount   `xml:"TtlRtrdIntrBkSttlmAmt"`
	IntrBkSttlmDt          string   `xml:"IntrBkSttlmDt"`
	SttlmInf               sttlmInf `xml:"SttlmInf"`
}

type sttlmInf struct {
	SttlmMtd string  `xml:"SttlmMtd"`
}

type txInf struct {
	RtrId              string       `xml:"RtrId"`
	OrgnlGrpInf        orgnlGrpInf  `xml:"OrgnlGrpInf"`
	OrgnlInstrId       string       `xml:"OrgnlInstrId"`
	OrgnlEndToEndId    string       `xml:"OrgnlEndToEndId"`
	OrgnlTxId          string       `xml:"OrgnlTxId"`
	RtrdIntrBkSttlmAmt amount       `xml:"RtrdIntrBkSttlmAmt"`
	IntrBkSttlmDt      string       `xml:"IntrBkSttlmDt"`
	RtrdInstdAmt       amount       `xml:"RtrdInstdAmt"`
	ChrgBr             string       `xml:"ChrgBr"`
	RtrRsnInf          rtrRsnInf    `xml:"RtrRsnInf"`
	OrgnlTxRef         orgnlTxRef   `xml:"OrgnlTxRef"`
}

type orgnlGrpInf struct {
	OrgnlMsgId   string `xml:"OrgnlMsgId"`
	OrgnlMsgNmId string `xml:"OrgnlMsgNmId"`
}

type amount struct {
	Ccy string `xml:"Ccy,attr"`
	Val string `xml:",chardata"`
}

type rtrRsnInf struct {
	Rsn      rsn    `xml:"Rsn"`
	AddtlInf string `xml:"AddtlInf"`
}

type rsn struct {
	Cd string `xml:"Cd"`
}

type orgnlTxRef struct {
	IntrBkSttlmDt string   `xml:"IntrBkSttlmDt"`
	SttlmPrty     string   `xml:"SttlmPrty"`
	RmtInf        rmtInf   `xml:"RmtInf"`
	Dbtr          party    `xml:"Dbtr"`
	DbtrAcct      account  `xml:"DbtrAcct"`
	DbtrAgt       agent    `xml:"DbtrAgt"`
	CdtrAgt       agent    `xml:"CdtrAgt"`
	Cdtr          party    `xml:"Cdtr"`
	CdtrAcct      account  `xml:"CdtrAcct"`
}

type party struct {
	Nm      string   `xml:"Nm"`
	PstlAdr pstlAdr  `xml:"PstlAdr"`
}

type pstlAdr struct {
	StrtNm string `xml:"StrtNm"`
	BldgNb string `xml:"BldgNb"`
	PstCd  string `xml:"PstCd"`
	TwnNm  string `xml:"TwnNm"`
	Ctry   string `xml:"Ctry"`
}

type account struct {
    Id id     `xml:"Id"`
    Tp accTp  `xml:"Tp"`  // Novo campo para tipo de conta
}

type id struct {
    Othr othr `xml:"Othr"`
}

type accTp struct {
    Prtry string `xml:"Prtry"`  // Tipo da conta (ex: CACC para conta corrente)
}

type othr struct {
    Id      string     `xml:"Id"`
    SchmeNm schemeName `xml:"SchmeNm"`  // Novo campo para identificar o tipo
}

type schemeName struct {
    Prtry string `xml:"Prtry"`  // PHONE, EMAIL ou CIL
}

type agent struct {
	FinInstnId finInstnId `xml:"FinInstnId"`
}

type finInstnId struct {
	BICFI string `xml:"BICFI"`
}

type rmtInf struct {
	Ustrd string `xml:"Ustrd"`
}