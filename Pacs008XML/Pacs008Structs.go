package Pacs008XML

import (
    "encoding/xml"
)

type DocumentPacs008 struct {
    XMLName xml.Name `xml:"Document"`
    XMLNS   string   `xml:"xmlns,attr"`
    FIToFICstmrCdtTrf fiToFICstmrCdtTrf `xml:"FIToFICstmrCdtTrf"`
}

type fiToFICstmrCdtTrf struct {
    GrpHdr      grpHdr        `xml:"GrpHdr"`
    CdtTrfTxInf []cdtTrfTxInf `xml:"CdtTrfTxInf"`
}

type grpHdr struct {
    MsgId    string    `xml:"MsgId"`
    CreDtTm  string    `xml:"CreDtTm"`
    NbOfTxs  string       `xml:"NbOfTxs"`
    SttlmInf sttlmInf  `xml:"SttlmInf"`
}

type sttlmInf struct {
    SttlmMtd string `xml:"SttlmMtd"`
}

type cdtTrfTxInf struct {
    PmtId           pmtId           `xml:"PmtId"`
    InstdAmt        instdAmt        `xml:"InstdAmt"`
    IntrBkSttlmAmt  intrBkSttlmAmt  `xml:"IntrBkSttlmAmt"`
    CdtrAgt         cdtrAgt         `xml:"CdtrAgt"`
    Cdtr            cdtr            `xml:"Cdtr"`
    CdtrAcct        cdtrAcct        `xml:"CdtrAcct"`
    DbtrAgt         dbtrAgt         `xml:"DbtrAgt"`
    Dbtr            dbtr            `xml:"Dbtr"`
    DbtrAcct        dbtrAcct        `xml:"DbtrAcct"`
    IntrmyAgt1      intrmyAgt1      `xml:"IntrmyAgt1"`
    IntrmyAgt1Acct  intrmyAgt1Acct  `xml:"IntrmyAgt1Acct"`
    ChrgBr          string          `xml:"ChrgBr"`
    ChrgsInf        chrgsInf        `xml:"ChrgsInf"`
    RmtInf          rmtInf          `xml:"RmtInf"`
}

// Restante das estruturas permanece igual
type pmtId struct {
    InstrId    string `xml:"InstrId"`
    EndToEndId string `xml:"EndToEndId"`
}

type instdAmt struct {
    Ccy string `xml:"Ccy,attr"`
    Val string `xml:",chardata"`
}

type intrBkSttlmAmt struct {
    Ccy string `xml:"Ccy,attr"`
    Val string `xml:",chardata"`
}

type intrmyAgt1 struct {
    FinInstnId finInstnId `xml:"FinInstnId"`
}

type intrmyAgt1Acct struct {
    Nm  string `xml:"Nm"`
    Ccy string `xml:"Ccy"`
    Id  id     `xml:"Id"`
}

type chrgsInf struct {
    Amt  instdAmt  `xml:"Amt"`
    Agt  cdtrAgt   `xml:"Agt"`
}

type cdtrAgt struct {
    FinInstnId finInstnId `xml:"FinInstnId"`
}

type dbtrAgt struct {
    FinInstnId finInstnId `xml:"FinInstnId"`
}

type finInstnId struct {
    Nm    string `xml:"Nm"`
    BICFI string `xml:"BICFI"`
}

type cdtr struct {
    Nm      string   `xml:"Nm"`
    PstlAdr pstlAdr  `xml:"PstlAdr"`
}

type dbtr struct {
    Nm      string   `xml:"Nm"`
    PstlAdr pstlAdr  `xml:"PstlAdr"`
}

type pstlAdr struct {
    Ctry string `xml:"Ctry"`
}

type cdtrAcct struct {
    Ccy string       `xml:"Ccy"`
    Id  id           `xml:"Id"`
    Tp  accountType  `xml:"Tp"`
}

type dbtrAcct struct {
    Ccy string       `xml:"Ccy"`
    Id  id           `xml:"Id"`
    Tp  accountType  `xml:"Tp"`
}

type id struct {
    Othr othr `xml:"Othr"`
}

type othr struct {
    Id      string     `xml:"Id"`
    SchmeNm schemeName `xml:"SchmeNm"`
}

type schemeName struct {
    Prtry string `xml:"Prtry"`
}

type accountType struct {
    Prtry string `xml:"Prtry"`
}

type rmtInf struct {
    Ustrd string `xml:"Ustrd"`
}