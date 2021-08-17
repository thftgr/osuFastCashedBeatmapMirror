package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Logger"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Route"
	"github.com/thftgr/osuFastCashedBeatmapMirror/middleWareFunc"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
)

var LogIO = bytes.Buffer{}

func init() {
	ch := make(chan struct{})

	src.LoadSetting()

	go src.StartIndex()
	go src.LoadBancho(ch)
	src.ConnectMaria()
	go Logger.LoadLogger(&LogIO)
	_ = <-ch


}
func main() {


	e := echo.New()
	e.HideBanner = true

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Output: &LogIO}),
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
	)

	e.GET("/robots.txt", Route.Robots)
	e.GET("/d/:id", Route.DownloadBeatmapSet)
	e.GET("/search", Route.Search)

	e.Logger.Fatal(e.Start(":" + src.Setting.Port))

}
