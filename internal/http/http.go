package http

import (
	"context"
	"emailer/internal/domain"
	"emailer/internal/logger"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
	"time"
)

func runHttpServer(ctx context.Context, wg *sync.WaitGroup, app domain.App) {
	wg.Add(1)
	defer wg.Done()
	engine := html.New("./internal/http/views", ".html")
	httpSrv := fiber.New(fiber.Config{
		Views:       engine,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	})
	httpSrv.Use(recover.New())
	handlers := handler{App: app, errCh: logger.GetLoggerCh()}
	httpSrv.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{"Name": 888})
	})
	httpSrv.Post("/send", handlers.send)
	httpSrv.Post("/send-async", handlers.sendAsync)
	httpSrv.Get("/result", handlers.getResult)
	httpSrv.Post("/test", handlers.test)
	port := ":8090"
	if p, exists := os.LookupEnv("HTTP_PORT"); exists {
		port = p
	}
	go func() {
		log.Printf("Start server at %v\n", port)
		if err := httpSrv.Listen(port); err != nil {
			log.Fatalln(err)
		}
	}()
	<-ctx.Done()
	log.Println("Stopping http server")
	err := httpSrv.Shutdown()
	if err != nil {
		handlers.errCh <- err
	}
}

func RunServer(ctx context.Context, wg *sync.WaitGroup, app domain.App) {
	go runHttpServer(ctx, wg, app)
}

type handler struct {
	domain.App
	errCh chan error
}

func (h *handler) test(ctx *fiber.Ctx) error {
	var data domain.Payload
	err := json.Unmarshal(ctx.Body(), &data)
	if err != nil {
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	message, err := h.RenderTemplate(data)
	if err != nil {
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	_, err = ctx.WriteString(string(message))
	return err
}

func (h *handler) send(ctx *fiber.Ctx) error {
	var data domain.Payload
	err := json.Unmarshal(ctx.Body(), &data)
	if err != nil {
		err := fmt.Errorf("send handler unmarshal err, received %v; %w", string(ctx.Body()), err)
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	err = h.Do(ctx.Context(), data)
	if err != nil {
		err := fmt.Errorf("send handler Do err, received %+v\n; %w", data, err)
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"success": "true"})
}
func (h *handler) getResult(ctx *fiber.Ctx) error {
	u := ctx.Query("uuid", "")
	id, err := uuid.Parse(u)
	if err != nil {
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	ready, err := h.AsyncResult(id.String())
	return ctx.JSON(fiber.Map{"success": ready, "error": err})
}

func (h *handler) sendAsync(ctx *fiber.Ctx) error {
	var data domain.Payload
	err := json.Unmarshal(ctx.Body(), &data)
	if err != nil {
		err := fmt.Errorf("sendAsync handler unmarshal err, received %+v\n; %w", string(ctx.Body()), err)
		h.errCh <- err
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	u := h.DoAsync(data)
	return ctx.JSON(fiber.Map{"id": u, "success": "true"})
}
