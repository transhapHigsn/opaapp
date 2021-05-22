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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/lightstep/otel-launcher-go/launcher"
	fiberOtel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
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
	// initTracer()

	pid := os.Getpid()
	access_token, ok := os.LookupEnv("LIGHTSTEP_ACCESS_TOKEN")
	if !ok {
		log.Fatalf("pid=%d level=danger msg=Unable to find access token.", pid)
	}

	host, _ := os.Hostname()
	environment, ok := os.LookupEnv("OPAAPP_ENV")
	if !ok {
		log.Fatalf("pid=%d level=danger msg=environment variable lookup failure.", pid)
	}

	ls := launcher.ConfigureOpentelemetry(
		launcher.WithAccessToken(access_token),
		launcher.WithServiceName("opaapp"),
		launcher.WithServiceVersion("v0.0.1"),
		launcher.WithResourceAttributes(map[string]string{
			"host.hostname":       host,
			"container.name":      "my-container-name",
			"cloud.region":        "ap-south-1",
			"service.environment": environment,
			"process.pid":         strconv.Itoa(pid),
		}),
	)
	defer ls.Shutdown()

	env_prefork := os.Getenv("OPAAPP_PREFORK")
	if env_prefork != "" {
		prefork, _ := strconv.ParseBool(env_prefork)
		log.Printf("pid=%d level=info msg=prefork value from env file: %s", pid, env_prefork)
		fiberConfig.Prefork = prefork
	}

	if !fiberConfig.Prefork {
		log.Printf("pid=%d level=info msg=Preforking disabled.", pid)
	}

	app := fiber.New(fiberConfig)

	otelMiddleware := fiberOtel.New(fiberOtel.Config{
		// name for root span in trace on request
		SpanName: "http/request",
		// array of span options for root span
		TracerStartAttributes: []trace.SpanOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
		},
		// key name for context store in fiber.Ctx
		LocalKeyName: "otel-context",
	})
	app.Use(otelMiddleware)

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
	app.Get("/span", spanCheck)

	opaapp_port := os.Getenv("OPAAPP_PORT")
	port := os.Getenv("PORT")

	var application_port string
	if port != "" {
		application_port = port
	} else if opaapp_port != "" {
		application_port = opaapp_port
	} else {
		port = "3000"
	}

	listen_on := fmt.Sprintf(":%s", application_port)

	log.Printf("pid=%d level=info msg=Starting up server ...", pid)
	log.Printf("pid=%d level=info msg=Server listening on -> %s ", pid, listen_on)
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
		ctx := fiberOtel.FromCtx(c)
		_, span := fiberOtel.Tracer.Start(ctx, "db-time-finder-closure")
		defer span.End()

		rid := c.Locals(requestIdConfig.ContextKey).(string)
		span.SetAttributes(attribute.String("request.id", rid))

		span.AddEvent("create-response")
		content := env.createResponseForGetTime(ctx, "DB call check by closure")
		span.AddEvent("return-response")
		return c.JSON(content)
	}
}

func (env *Env) getDbTimeFiber(c *fiber.Ctx) error {

	ctx := fiberOtel.FromCtx(c)
	_, span := fiberOtel.Tracer.Start(ctx, "db-time-finder")
	defer span.End()

	rid := c.Locals(requestIdConfig.ContextKey).(string)
	span.SetAttributes(attribute.String("request.id", rid))

	span.AddEvent("create-response")
	content := env.createResponseForGetTime(ctx, "DB call check.")
	span.AddEvent("return-response")
	return c.JSON(content)
}

func runRegoPolicy(c *fiber.Ctx) error {

	ctx := fiberOtel.FromCtx(c)
	_, span := fiberOtel.Tracer.Start(ctx, "rego-context")
	defer span.End()

	span.AddEvent("load-request-id")
	rid := c.Locals(requestIdConfig.ContextKey).(string)

	span.SetAttributes(attribute.String("request.id", rid))

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

	span.AddEvent("run-rego-query")
	result := opa.RunRegoQuery(ctx, input)

	x := result[0].Bindings["result"]
	span.AddEvent("return-rego-result")

	return c.JSON(fiber.Map{
		"output": x,
	})
}

func spanCheck(c *fiber.Ctx) error {
	ctx := fiberOtel.FromCtx(c)

	// use retrieved context
	_, span := fiberOtel.Tracer.Start(ctx, "nested-route-tracer")
	span.AddEvent("get-post")
	span.AddEvent("get-comments")
	span.AddEvent("get-author")
	defer span.End()

	return c.JSON(fiber.Map{
		"output": "hello",
	})
}

func (env *Env) getTime(ctx context.Context) time.Time {
	_, span := fiberOtel.Tracer.Start(
		ctx,
		"env.getTime",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "now"),
		),
	)

	defer span.End()

	var now time.Time
	err := env.dbPool.QueryRow(context.Background(), "select now()").Scan(&now)
	if err != nil {
		panic(err)
	}
	return now
}

func (env *Env) createResponseForGetTime(ctx context.Context, reply string) response {
	_, span := fiberOtel.Tracer.Start(
		ctx,
		"env.createResponseForGetTime",
		trace.WithSpanKind(trace.SpanKindClient),
		// trace.WithAttributes(
		// 	attribute.String("db.system", "postgresql"),
		// 	attribute.String("db.operation", "now"),
		// ),
	)

	defer span.End()
	var content response
	content.Timestamp = env.getTime(ctx)
	content.Random = rand.Intn(10000)
	content.Response = reply

	return content
}
