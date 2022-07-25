package http

import (
	"context"
	"emailer/internal/domain"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"log"
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
	httpSrv.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{"Name": 888})
	})
	httpSrv.Post("/", func(ctx *fiber.Ctx) error {
		var data domain.Payload
		err := json.Unmarshal(ctx.Body(), &data)
		if err != nil {
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}
		err = app.Do(data)
		if err != nil {
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.JSON(fiber.Map{"success": "true"})
	})
	httpSrv.Post("/test", func(ctx *fiber.Ctx) error {
		var data domain.Payload
		err := json.Unmarshal(ctx.Body(), &data)
		if err != nil {
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}
		message, err := app.RenderTemplate(data)
		if err != nil {
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}
		ctx.WriteString(string(message))
		return nil
	})

	port := ":8090"
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
		log.Println(err)
	}
}

func RunServer(ctx context.Context, wg *sync.WaitGroup, app domain.App) {
	go runHttpServer(ctx, wg, app)
}
