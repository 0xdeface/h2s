package domain

import "html/template"

type Sender interface {
	Send(subject string, message []byte, from string, to []string) error
}

type MessageMaker interface {
	Make(*template.Template, interface{}) ([]byte, error)
}
type Payload struct {
	TemplateName string      `json:"templateName"`
	To           []string    `json:"to"`
	From         string      `json:"from"`
	Subject      string      `json:"subject"`
	Data         interface{} `json:"data"`
}
