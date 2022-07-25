package domain

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	sender Sender
	maker  MessageMaker
}

// Do
// Синхронная отправка
func (app *App) Do(payload Payload) error {
	message, err := app.RenderTemplate(payload)
	if err != nil {
		return err
	}
	fnWrap := func(ctx context.Context) error {
		return app.sender.Send(payload.Subject, message, payload.From, payload.To)
	}
	return Retry(fnWrap, 3, time.Second*2)(context.TODO())
}

// DoAsync
// Помещает задачу по отправке в очередь и сразу возвращает
// идентификатор отправки
func DoAsync(payload Payload) string {
	go func() {}()
	return "uniqum identifier"
}

// AsyncResult
// Результат выполнения задачи по идентификатору
// Возвращает готовность задачи и результат
func AsyncResult(u string) (bool, error) {
	return false, nil
}
func (app *App) RenderTemplate(payload Payload) ([]byte, error) {
	tplPath, err := normalizePath(payload.TemplateName)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return nil, err
	}
	return app.maker.Make(tpl, payload.Data)
}
func normalizePath(rawPath string) (normalizedPath string, err error) {
	basePath, err := os.Getwd()
	targetPath := fmt.Sprintf("%v/tpls/%v", basePath, rawPath)
	pattern := fmt.Sprintf("%v/tpls/*", basePath)
	match, err := filepath.Match(pattern, targetPath)
	if !match || err != nil {
		return "", errors.New("incorrect path")
	}
	return targetPath, nil
}

func NewApp(sender Sender, maker MessageMaker) App {
	return App{
		sender: sender,
		maker:  maker,
	}
}
