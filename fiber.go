package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

var requestIdConfig = requestid.Config{
	Generator: func() string {
		return uuid.New().String()
	},
	Header:     fiber.HeaderXRequestID,
	ContextKey: "requestid",
}

var fiberConfig = fiber.Config{
	Prefork:               true,
	ReadTimeout:           5 * time.Second,
	WriteTimeout:          5 * time.Second,
	DisableStartupMessage: true,
}

func fiberApp() {
	app := fiber.New(fiberConfig)

	app.Use(requestid.New(requestIdConfig))

	app.Use(logger.New(logger.Config{
		TimeZone:   "Asia/Kolkata",
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} pid=${pid} requestid=${locals:requestid} path=${path} method=${method} status=${status} latency=${latency}\n",
	}))
	app.Use(recover.New())

	dbPool := initDb()
	env := &Env{dbPool: dbPool}
	defer dbPool.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/time", getSystemTimeFiber)
	app.Get("/closure", getDbTimeByClosureFiber(env))
	app.Get("/db", env.getDbTimeFiber)

	port := os.Getenv("OPAAPP_PORT")
	if port == "" {
		port = "3000"
	}

	pid := os.Getpid()
	listen_on := fmt.Sprintf(":%s", port)

	log.Printf("pid=%d Starting up server ...", pid)
	log.Printf("pid=%d Server listening on -> %s ", pid, listen_on)
	log.Fatal(app.Listen(listen_on))
}

func getSystemTimeFiber(c *fiber.Ctx) error {
	var content response
	content.Response = "System time check"
	content.Timestamp = time.Now().UTC()
	content.Random = rand.Intn(1000)

	return c.JSON(content)
}

func getDbTimeByClosureFiber(env *Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var now time.Time

		rid := c.Locals(requestIdConfig.ContextKey)

		log.Printf("rid=%s func=getDbTimeByClosureFiber level=info msg=Fetching time from db.", rid)
		err := env.dbPool.QueryRow(context.Background(), "select now()").Scan(&now)
		if err != nil {
			panic(err)
		}
		log.Printf("rid=%s func=getDbTimeByClosureFiber level=info msg=Fetched time from db.", rid)

		var content response

		content.Response = "DB call check by closure"
		content.Timestamp = now
		content.Random = rand.Intn(1000)

		return c.JSON(content)
	}
}

func (env *Env) getDbTimeFiber(c *fiber.Ctx) error {
	var now time.Time

	rid := c.Locals(requestIdConfig.ContextKey)

	log.Printf("rid=%s func=getDbTimeFiber level=info msg=Fetching time from db.", rid)
	err := env.dbPool.QueryRow(context.Background(), "select now()").Scan(&now)
	if err != nil {
		panic(err)
	}
	log.Printf("rid=%s func=getDbTimeFiber level=info msg=Fetched time from db.", rid)

	var content response

	content.Response = "DB call check"
	content.Timestamp = now
	content.Random = rand.Intn(1000)

	return c.JSON(content)
}
