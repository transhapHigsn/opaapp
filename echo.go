package main

// import (
// 	"context"
// 	"math/rand"
// 	"net/http"
// 	"time"

// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// 	"github.com/labstack/gommon/log"
// )

// func echoApp() {
// 	e := echo.New()
// 	e.HideBanner = true

// 	s := &http.Server{
// 		Addr:         ":1323",
// 		ReadTimeout:  20 * time.Minute,
// 		WriteTimeout: 20 * time.Minute,
// 	}
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())

// 	e.Logger.SetLevel(log.DEBUG)

// 	dbPool := initDb()
// 	env := &Env{dbPool: dbPool}
// 	defer dbPool.Close()

// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Hello, World!")
// 	})

// 	e.GET("/time", getSystemTime)
// 	e.GET("/db", env.getDbTime)
// 	e.GET("/closure", getDbTimeByClosure(env))

// 	e.Logger.Fatal(e.StartServer(s))
// }

// func (env *Env) getDbTime(c echo.Context) error {
// 	var now time.Time

// 	// following log command doesn't works.
// 	// have to look into that.
// 	log.Debug("Querying DB for time.")
// 	err := env.dbPool.QueryRow(context.Background(), "select now()").Scan(&now)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var content response

// 	content.Response = "DB call check"
// 	content.Timestamp = now
// 	content.Random = rand.Intn(1000)

// 	return c.JSON(http.StatusOK, &content)
// }

// func getSystemTime(c echo.Context) error {
// 	var content response
// 	content.Response = "System time check"
// 	content.Timestamp = time.Now().UTC()
// 	content.Random = rand.Intn(1000)

// 	return c.JSON(http.StatusOK, &content)
// }

// func getDbTimeByClosure(env *Env) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var now time.Time

// 		err := env.dbPool.QueryRow(context.Background(), "select now()").Scan(&now)
// 		if err != nil {
// 			panic(err)
// 		}

// 		var content response

// 		content.Response = "DB call check by closure"
// 		content.Timestamp = now
// 		content.Random = rand.Intn(1000)

// 		return c.JSON(http.StatusOK, &content)
// 	}
// }
