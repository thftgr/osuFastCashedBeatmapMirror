package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"github.com/swaggo/echo-swagger"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Logger"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/Logger"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Route"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/bootLoader"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/docs"
	"github.com/thftgr/osuFastCashedBeatmapMirror/middleWareFunc"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"log"
)
//https://github.com/swaggo/swag#declarative-comments-format
// @title Swagger Example API
// @version 1.0
// @description thftgr's fast cached osu beatmap mirror.
// @termsOfService https://github.com/thftgr/osuFastCashedBeatmapMirror

// @contact.name API Support
// @contact.email thftgr@gmail.com

// @license.name GNU General Public License v3.0
// @license.url https://github.com/thftgr/osuFastCashedBeatmapMirror/blob/master/LICENSE

// @host xiiov.com
// @BasePath /
func main() {

	e := echo.New()
	e.HideBanner = true
	go func() {
		for {
			<-Logger.Ch
			e.Logger.SetOutput(log.Writer())
			pterm.Info.Println("UPDATED ECHO LOGGER.")
		}
	}()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(
		middleware.Logger(),
		//middleware.LoggerWithConfig(middleware.LoggerConfig{Output: log.Writer()}),
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
	)
	// SWAGGER =========================================================================================================
	e.GET("/*", echoSwagger.WrapHandler)

	// 서버상태 체크용 ====================================================================================================


	e.GET("/health", Route.Health)

	e.GET("/robots.txt", Route.Robots)

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:id", Route.DownloadBeatmapSet)

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", Route.Search)

	pterm.Info.Println("ECHO STARTED AT", src.Setting.Port)
	e.Logger.Fatal(e.Start(":" + src.Setting.Port))

}
