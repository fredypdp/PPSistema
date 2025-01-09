package Pacs008XML

import (
	"encoding/xml"
)

type Document struct {
	XMLName xml.Name `xml:"Document"`
	XMLNS   string   `xml:"xmlns,attr"`
	FIToFICstmrCdtTrf FIToFICstmrCdtTrf `xml:"FIToFICstmrCdtTrf"`
}

type FIToFICstmrCdtTrf struct {
	GrpHdr     GrpHdr     `xml:"GrpHdr"`
	CdtTrfTxInf CdtTrfTxInf `xml:"CdtTrfTxInf"`
}

type GrpHdr struct {
	MsgId    string    `xml:"MsgId"`
	CreDtTm  string    `xml:"CreDtTm"`
	NbOfTxs  string    `xml:"NbOfTxs"`
	SttlmInf SttlmInf  `xml:"SttlmInf"`
}

type SttlmInf struct {
	SttlmMtd string `xml:"SttlmMtd"`
}

type CdtTrfTxInf struct {
	PmtId          PmtId          `xml:"PmtId"`
	IntrBkSttlmAmt IntrBkSttlmAmt `xml:"IntrBkSttlmAmt"`
	CdtrAgt        CdtrAgt        `xml:"CdtrAgt"`
	Cdtr           Cdtr           `xml:"Cdtr"`
	CdtrAcct       CdtrAcct       `xml:"CdtrAcct"`
	DbtrAgt        DbtrAgt        `xml:"DbtrAgt"`
	Dbtr           Dbtr           `xml:"Dbtr"`
	DbtrAcct       DbtrAcct       `xml:"DbtrAcct"`
	ChrgBr         string         `xml:"ChrgBr"`
	RmtInf         RmtInf         `xml:"RmtInf"`
}

type PmtId struct {
	InstrId    string `xml:"InstrId"`
	EndToEndId string `xml:"EndToEndId"`
}

type IntrBkSttlmAmt struct {
	Ccy string `xml:"Ccy,attr"`
	Val string `xml:",chardata"`
}

type CdtrAgt struct {
	FinInstnId FinInstnId `xml:"FinInstnId"`
}

type FinInstnId struct {
	BICFI string `xml:"BICFI"`
}

type Cdtr struct {
	Nm      string    `xml:"Nm"`
	PstlAdr PstlAdr   `xml:"PstlAdr"`
}

type CdtrAcct struct {
	Id Id `xml:"Id"`
}

type DbtrAgt struct {
	FinInstnId FinInstnId `xml:"FinInstnId"`
}

type Dbtr struct {
	Nm      string    `xml:"Nm"`
	PstlAdr PstlAdr   `xml:"PstlAdr"`
}

type DbtrAcct struct {
	Id Id `xml:"Id"`
}

type Id struct {
	Othr Othr `xml:"Othr"`
}

type Othr struct {
	Id string `xml:"Id"`
}

type PstlAdr struct {
	Ctry string `xml:"Ctry"`
}

type RmtInf struct {
	Ustrd string `xml:"Ustrd"`
}