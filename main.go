package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/Logger"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Route"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/bootLoader"
	"github.com/thftgr/osuFastCashedBeatmapMirror/middleWareFunc"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"log"
)


func main() {

	e := echo.New()
	e.HideBanner = true

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{Output: log.Writer()}),
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
	)

	e.GET("/robots.txt", Route.Robots)
	e.GET("/d/:id", Route.DownloadBeatmapSet)
	e.GET("/search", Route.Search)

	e.Logger.Fatal(e.Start(":" + src.Setting.Port))

}
