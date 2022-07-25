package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	sender      Sender
	maker       MessageMaker
	queue       chan QueueMessage
	queueResult map[string]ExecutionResult
}

type QueueMessage struct {
	id      string
	payload Payload
}
type ExecutionResult struct {
	err   error
	ready bool
}

// Do
// Синхронная отправка
func (app *App) Do(ctx context.Context, payload Payload) error {
	message, err := app.RenderTemplate(payload)
	if err != nil {
		return err
	}
	fnWrap := func(ctx context.Context) error {
		return app.sender.Send(payload.Subject, message, payload.From, payload.To)
	}
	return Retry(fnWrap, 3, time.Second*2)(ctx)
}

// DoAsync
// Помещает задачу по отправке в очередь и сразу возвращает
// идентификатор отправки
func (app *App) DoAsync(payload Payload) string {
	u, _ := uuid.NewUUID()
	app.queueResult[u.String()] = ExecutionResult{ready: false, err: nil}
	app.queue <- QueueMessage{
		id:      u.String(),
		payload: payload,
	}
	return u.String()
}

// AsyncResult
// Результат выполнения задачи по идентификатору
// Возвращает готовность задачи и результат
func (app *App) AsyncResult(u string) (bool, error) {
	if val, exists := app.queueResult[u]; exists {
		if val.ready {
			delete(app.queueResult, u)
		}
		return val.ready, val.err
	}
	return false, fmt.Errorf("task with id:%v not exist", u)
}

func (app *App) Worker(ctx context.Context) {
	for {
		select {
		case m := <-app.queue:
			app.queueResult[m.id] = ExecutionResult{
				err:   app.Do(ctx, m.payload),
				ready: true,
			}
		case <-ctx.Done():
			return
		}
	}
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

func NewApp(sender Sender, maker MessageMaker, ctx context.Context) App {
	app := App{
		sender:      sender,
		maker:       maker,
		queue:       make(chan QueueMessage, 40000),
		queueResult: make(map[string]ExecutionResult),
	}
	go app.Worker(ctx)
	return app
}
