package message_maker

import (
	"bytes"
	"html/template"
)

type MessageMaker struct {
}

func (m MessageMaker) Make(template *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := template.Execute(&buf, data)
	return buf.Bytes(), err
}
