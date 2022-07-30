package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, strconv.Itoa(os.Getpid()))
	})
	e.Logger.Fatal(e.Start(":80"))
}
