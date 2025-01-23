package Pacs002XML

import (
    "encoding/xml"
)

type DocumentPacs002 struct {
    XMLName           xml.Name        `xml:"Document"`
    XMLNS             string          `xml:"xmlns,attr"`
    FIToFIPmtStsRpt   fiToFIPmtStsRpt `xml:"FIToFIPmtStsRpt"`
}

type fiToFIPmtStsRpt struct {
    GrpHdr             grpHdr             `xml:"GrpHdr"`
    OrgnlGrpInfAndSts  orgnlGrpInfAndSts  `xml:"OrgnlGrpInfAndSts"`
    TxInfAndSts        []txInfAndSts      `xml:"TxInfAndSts"`
}

type grpHdr struct {
    MsgId   string `xml:"MsgId"`
    CreDtTm string `xml:"CreDtTm"`
}

type orgnlGrpInfAndSts struct {
    OrgnlMsgId    string `xml:"OrgnlMsgId"`
    OrgnlMsgNmId  string `xml:"OrgnlMsgNmId"`
    OrgnlCreDtTm  string `xml:"OrgnlCreDtTm"`
    GrpSts        string `xml:"GrpSts"`
}

type txInfAndSts struct {
    OrgnlEndToEndId string         `xml:"OrgnlEndToEndId"`
    OrgnlTxId       string         `xml:"OrgnlTxId"`
    TxSts           string         `xml:"TxSts"`
    StsRsnInf       []stsRsnInf    `xml:"StsRsnInf,omitempty"` // Opcional
    InstdAgt        instdAgt       `xml:"InstdAgt"`
    OrgnlTxRef      orgnlTxRef     `xml:"OrgnlTxRef"`
}

type stsRsnInf struct {
    Rsn        rsn      `xml:"Rsn"`
    AddtlInf   []string `xml:"AddtlInf"` // Campo que permite múltiplas informações adicionais
}

type rsn struct {
    Cd    string `xml:"Cd"`
    Prtry string `xml:"Prtry,omitempty"` // Campo opcional para razão proprietária
}

type instdAgt struct {
    FinInstnId finInstnId `xml:"FinInstnId"`
}

type finInstnId struct {
    BICFI string `xml:"BICFI"`
}

type orgnlTxRef struct {
    IntrBkSttlmAmt intrBkSttlmAmt `xml:"IntrBkSttlmAmt"`
}

type intrBkSttlmAmt struct {
    Ccy string `xml:"Ccy,attr"`
    Val string `xml:",chardata"`
}