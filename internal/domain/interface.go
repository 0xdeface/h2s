package domain

import "html/template"

type Sender interface {
	Send(subject string, message []byte, from string, to []string) error
}

type MessageMaker interface {
	Make(*template.Template, ...interface{}) ([]byte, error)
}
