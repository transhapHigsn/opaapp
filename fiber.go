package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"opaapp/opa"
	"os"
	"strconv"
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
	Prefork:               false,
	ReadTimeout:           5 * time.Second,
	WriteTimeout:          5 * time.Second,
	DisableStartupMessage: true,
}

func fiberApp() {
	pid := os.Getpid()

	env_prefork := os.Getenv("OPAAPP_PREFORK")
	if env_prefork != "" {
		prefork, _ := strconv.ParseBool(env_prefork)
		fiberConfig.Prefork = prefork
	}

	if !fiberConfig.Prefork {
		log.Printf("pid=%d level=info Preforking disabled.", pid)
	}

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
	app.Get("/rego", runRegoPolicy)

	port := os.Getenv("OPAAPP_PORT")
	if port == "" {
		port = "3000"
	}

	listen_on := fmt.Sprintf(":%s", port)

	log.Printf("pid=%d level=info Starting up server ...", pid)
	log.Printf("pid=%d level=info Server listening on -> %s ", pid, listen_on)
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

func runRegoPolicy(c *fiber.Ctx) error {

	input := map[string]interface{}{
		"pet_list": []map[string]interface{}{
			{
				"breed":           "St. Bernard",
				"name":            "Cujo",
				"up_for_adoption": false,
			},
			{
				"breed":           "Collie",
				"name":            "Lassie",
				"up_for_adoption": true,
			},
		},
		"token": "eyJ1IjoiSFMyNTYiLCJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiYWxpY2UiLCJlbXBsb3llZSI6dHJ1ZSwibmFtZSI6IkFsaWNlIFNtaXRoIn0.vMBYEW8VK9XM7yPkKTu1C3Gy1tOq1A0d4-xYMkkRpEc",
	}

	rid := c.Locals(requestIdConfig.ContextKey)
	log.Printf("rid=%s func=runRegoPolicy level=info msg=Running rego policy.", rid)
	start := time.Now()
	result := opa.RunRegoQuery(input)
	elapsed := time.Since(start)
	log.Printf("rid=%s func=runRegoPolicy level=info msg=Policy ran successfully in %s.", rid, elapsed)

	x := result[0].Bindings["result"]

	return c.JSON(fiber.Map{
		"output": x,
	})
}
