package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

var kv *KeyVal

func StartKeyValServer(address string, storageFilePath string) {
	// Echo instance
	e := echo.New()
	kv = NewKeyVal(storageFilePath, time.Hour*24*30)

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.GET("keyval/*", handleGetKeyVal)
	e.POST("keyval/*", handlePostKeyVal)

	// Start server
	e.Logger.Fatal(e.Start(address))
}

func handleGetKeyVal(ctx echo.Context) error {
	key := ctx.Param("*")
	value, ok := kv.Get(key)

	if !ok {
		return echo.ErrNotFound
	}

	return ctx.String(200, value)
}

func handlePostKeyVal(ctx echo.Context) error {
	key := ctx.Param("*")
	value := ctx.FormValue("value")

	kv.Put(key, value)

	return ctx.String(200, "ok")
}
