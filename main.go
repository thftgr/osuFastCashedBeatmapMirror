package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Logger"
	"github.com/thftgr/osuFastCashedBeatmapMirror/Route"
	_ "github.com/thftgr/osuFastCashedBeatmapMirror/bootLoader"
	"github.com/thftgr/osuFastCashedBeatmapMirror/src"
	"log"
)

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
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		//middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
		middleware.Recover(),
	)
	// docs ============================================================================================================
	e.GET("/", Route.Wiki)

	// 서버상태 체크용 ====================================================================================================

	e.GET("/health", Route.Health)

	e.GET("/robots.txt", Route.Robots)

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:id", Route.DownloadBeatmapSet)

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", Route.Search)
	e.GET("/search/beatmap/:mi", Route.SearchByBeatmapId)
	e.GET("/search/beatmapset/:si", Route.SearchByBeatmapSetId)

	pterm.Info.Println("ECHO STARTED AT", src.Setting.Port)
	e.Logger.Fatal(e.Start(":" + src.Setting.Port))

}

//var (
//	upgrader = websocket.Upgrader{}
//)

//func hello(c echo.Context) error {
//	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return err
//	}
//	defer ws.Close()
//
//	for {
//		// Write
//		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
//		if err != nil {
//			c.Logger().Error(err)
//		}
//
//		// Read
//		_, msg, err := ws.ReadMessage()
//		if err != nil {
//			c.Logger().Error(err)
//		}
//		fmt.Printf("%s\n", msg)
//	}
//}
