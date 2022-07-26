package logger

import "log"

var ErrorCh = make(chan error, 10000)

func init() {
	for err := range ErrorCh {
		log.Println(err.Error())
	}
}
