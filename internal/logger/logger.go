package logger

import (
	"log"
	"sync"
)

var (
	once  sync.Once
	errCh chan error
)

func StartLogger() {
	go func() {
		for err := range GetLoggerCh() {
			log.Println(err.Error())
		}
	}()
}
func GetLoggerCh() chan error {
	once.Do(func() {
		errCh = make(chan error, 100000)
	})
	return errCh
}
