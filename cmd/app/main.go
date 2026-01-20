package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/xh3sh/go-real-time-chats/internal/db"
	"github.com/xh3sh/go-real-time-chats/internal/repo"
	"github.com/xh3sh/go-real-time-chats/internal/route"
	"github.com/xh3sh/go-real-time-chats/internal/templates"
)

const SERVER_PORT = ":80"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("env файл не найден, используются стандартные параметры")
	}

	e := echo.New()

	tmpl := templates.NewTemplates()
	e.Renderer = tmpl

	e.Use(middleware.Logger())

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb, err := db.NewRedisClient(redisAddr)
	if err != nil {
		e.Logger.Fatalf("Не удалось подключиться к Redis:", err)
	}

	ctx := context.Background()

	redisPrefix := os.Getenv("REDIS_KEY_PREFIX")
	if redisPrefix == "" {
		redisPrefix = "chats"
	}
	repo := repo.NewRedisRepository(rdb, redisPrefix)
	repo.StartCleanupZMessages(ctx)

	routes := route.New(e, repo)
	routes.InitRoute(tmpl)

	e.Logger.Fatal(e.Start(SERVER_PORT))
}
