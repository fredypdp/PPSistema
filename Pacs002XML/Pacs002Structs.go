package Pacs002XML

import (
    "encoding/xml"
)

type DocumentPacs002 struct {
    XMLName xml.Name          `xml:"Document"`
    XMLNS   string            `xml:"xmlns,attr"`
    FIToFIPmtStsRpt FIToFIPmtStsRpt `xml:"FIToFIPmtStsRpt"`
}

type FIToFIPmtStsRpt struct {
    GrpHdr         GrpHdr         `xml:"GrpHdr"`
    OrgnlGrpInfAndSts OrgnlGrpInfAndSts `xml:"OrgnlGrpInfAndSts"`
    TxInfAndSts    TxInfAndSts    `xml:"TxInfAndSts"`
}

type GrpHdr struct {
    MsgId   string `xml:"MsgId"`
    CreDtTm string `xml:"CreDtTm"`
}

type OrgnlGrpInfAndSts struct {
    OrgnlMsgId    string `xml:"OrgnlMsgId"`
    OrgnlMsgNmId  string `xml:"OrgnlMsgNmId"`
    OrgnlCreDtTm  string `xml:"OrgnlCreDtTm"`
    GrpSts        string `xml:"GrpSts"`
}

type TxInfAndSts struct {
    OrgnlEndToEndId string      `xml:"OrgnlEndToEndId"`
    OrgnlTxId       string      `xml:"OrgnlTxId"`
    TxSts           string      `xml:"TxSts"`
    InstdAgt        InstdAgt    `xml:"InstdAgt"`
    OrgnlTxRef      OrgnlTxRef  `xml:"OrgnlTxRef"`
}

type InstdAgt struct {
    FinInstnId FinInstnId `xml:"FinInstnId"`
}

type FinInstnId struct {
    BICFI string `xml:"BICFI"`
}

type OrgnlTxRef struct {
    IntrBkSttlmAmt IntrBkSttlmAmt `xml:"IntrBkSttlmAmt"`
    IntrBkSttlmDt  string         `xml:"IntrBkSttlmDt"`
}

type IntrBkSttlmAmt struct {
    Ccy string `xml:"Ccy,attr"`
    Val string `xml:",chardata"`
}