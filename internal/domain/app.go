package domain

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type App struct {
	sender Sender
	maker  MessageMaker
}

func (app *App) Do(payload Payload) error {
	message, err := app.RenderTemplate(payload)
	if err != nil {
		return err
	}
	return app.sender.Send(payload.Subject, message, payload.From, payload.To)
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
