package main

import (
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/Route"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/middleWareFunc"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"os"
	"runtime"
)

// TODO DONG DB 테이블 없으면 자동으로 생성하게
// TODO DONG 로그 디비에 넣을때 어떤 데이터 넣을지.
// TODO DONG 서버끼리 서로 비트맵파일 동기화 시킬수 있게
// TODO DONG 헤더로 프론트인지 api 인지 구분할수있게
//  	END  에러 핸들러.
//  	END  검색엔진 버그 체크하고 쿼리문 수정
//  	END  비트맵 반쵸에서 다운로드중에 클라이언트가 취소해도 서버는 계속 다운로드.
// TODO DONG 반쵸 비트맵 다운로드 제한 10분간 약 200건 10분 정지. (429 too many request)
// TODO DONG 검색 쿼리시 서버에 캐싱되어있는 비트맵인지 여부
func init() {
	ch := make(chan struct{})
	config.LoadConfig()
	src.StartIndex()
	db.ConnectMaria()
	go db.LoadIndex()
	go src.LoadBancho(ch)
	_ = <-ch
	if os.Getenv("debug") != "true" {
		go src.RunGetBeatmapDataASBancho()
	} else {
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true
	go func() {
		for {
			<-logger.Ch
			e.Logger.SetOutput(log.Writer())
			pterm.Info.Println("UPDATED ECHO LOGGER.")
		}
	}()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(

		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
		middleware.Recover(),
	)
	// docs ============================================================================================================
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, `https://nerinyan.stoplight.io/studio/nerinyan-api`)
	})

	// 서버상태 체크용 ====================================================================================================
	e.GET("/health", Route.Health)
	e.GET("/robots.txt", Route.Robots)
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Cpu Thread Count":     runtime.NumCPU(),
			"Running Thread Count": runtime.NumGoroutine(),
		})
	})

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:id", Route.DownloadBeatmapSet, middleWareFunc.BanchoBeatmapDownloadLimiter)
	//TODO 맵아이디, 맵셋아이디 지원
	//e.GET("/d/:id", Route.DownloadBeatmapSet, middleWareFunc.LoadBalancer)

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", Route.Search)
	e.GET("/search/beatmap/:mi", Route.SearchByBeatmapId)
	e.GET("/search/beatmapset/:si", Route.SearchByBeatmapSetId)

	// 서버 데이터 강제 업데이트용. ==========================================================================================
	// TODO 맵 굳이 한개씩 강제업데이트할 이유가 없음. 맵셋으로 업데이트만 지원
	//e.GET("/update/beatmapset/:id", func(c echo.Context) error {
	//
	//	//src.ManualUpdateBeatmapSet()
	//	return nil
	//})

	// 개발중 || 테스트중 ===================================================================================================
	//e.GET("/dev/search", func(c echo.Context) error {
	//	return c.JSON(http.StatusOK, db.SearchIndex(c.QueryParam("q")))
	//})

	//e.GET("/dev/test2", func(c echo.Context) error {
	//
	//	b := []int{1}
	//
	//	return c.String(http.StatusOK, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(b)), ","), "[]"))
	//})

	// 개발중 || 테스트중 ===================================================================================================

	pterm.Info.Println("ECHO STARTED AT", config.Setting.Port)
	e.Logger.Fatal(e.Start(":" + config.Setting.Port))

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
