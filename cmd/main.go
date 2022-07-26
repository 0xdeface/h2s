package main

import (
	"context"
	"emailer/internal/domain"
	"emailer/internal/http"
	"emailer/internal/logger"
	"emailer/internal/message_maker"
	sender "emailer/internal/sender/email"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger.StartLogger()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	Close := handleShutdown(ctx, cancel, wg)
	maker := message_maker.MessageMaker{}
	emailSender, err := sender.NewEmailSender()
	if err != nil {
		logger.GetLoggerCh() <- err
		cancel()
	}
	app := domain.NewApp(emailSender, maker, ctx)
	http.RunServer(ctx, wg, app)
	wg.Wait()
	Close()
	defer func() { close(logger.GetLoggerCh()) }()
	log.Println("Shutdown...")
}

func handleShutdown(ctx context.Context, cancel func(), wg *sync.WaitGroup) (Close func()) {
	stopCh := make(chan os.Signal, 1)
	wg.Add(1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	s := sync.Once{}
	go func() {
		select {
		case <-stopCh:
			s.Do(wg.Done)
			cancel()
		case <-ctx.Done():
			s.Do(wg.Done)
		}
	}()
	return func() {
		close(stopCh)
	}
}
