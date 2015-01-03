package platypus

import (
	"encoding/xml"
)

type RequestContainer struct {
	XMLName struct{} `xml:"PLATXML"`
	Header  string   `xml:"header"`
	Body    Body     `xml:"body"`
}

type ResponseContainer struct {
	XMLName struct{} `xml:"PLAT_XML"`
	Header  string   `xml:"header"`
	Body    Body     `xml:"body"`
}

type Body struct {
	Data DataBlock `xml:"data_block"`
}

type DataBlock struct {
	Protocol     string             `xml:"protocol"`
	Object       string             `xml:"object"`
	Action       string             `xml:"action"`
	Username     string             `xml:"username"`
	Password     string             `xml:"password"`
	Logintype    string             `xml:"logintype"`
	Properties   string             `xml:"properties"`
	Parameters   interface{}        `xml:"parameters"`
	ResponseCode string             `xml:"response_code,omitempty"`
	ResponseText string             `xml:"response_text,omitempty"`
	Success      int                `xml:"is_success,omitempty"`
	Attributes   AttributeDataBlock `xml:"attributes,omitempty"`
}

type AttributeDataBlock struct {
	Block []AttributeBlock `xml:"data_block,omitempty"`
}

type AttributeBlock struct {
	AttributeList []Attribute `xml:",any,omitempty"`
}

type Attribute struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func unwrapAttributeBlock(ab AttributeBlock) map[string]string {
	m := make(map[string]string)

	for _, v := range ab.AttributeList {
		m[v.XMLName.Local] = v.Value
	}

	return m
}
