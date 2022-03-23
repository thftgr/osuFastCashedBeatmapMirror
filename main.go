package main

import (
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/Route"
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCroller"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// TODO DOING DB 테이블 없으면 자동으로 생성하게
// TODO DOING 로그 디비에 넣을때 어떤 데이터 넣을지.
// TODO DOING 헤더로 프론트인지 api 인지 구분할수있게
//  	END   에러 핸들러.
//  	END   검색엔진 버그 체크하고 쿼리문 수정
//  	END   비트맵 반쵸에서 다운로드중에 클라이언트가 취소해도 서버는 계속 다운로드.
// TODO DOING 서버간 비트맵파일 해시값 비교해서 서로 다른경우 둘다 서버에서 삭제.
// TODO DOING 서버끼리 서로 비트맵파일 동기화 시킬수 있게
// TODO DOING 반쵸 비트맵 다운로드 제한 10분간 약 200건 10분 정지. (429 too many request) => 10분 내 100건 봇 감지 알고리즘
// TODO DOING 서버 자체적으로 10분당 150건 이내만 다운로드 가능하게 셋팅
// 		END	  검색 쿼리시 서버에 캐싱되어있는 비트맵인지 여부
// TODO DOING /status에 들어갈 상태값 추가.
// TODO DOING 반쵸에서 가져온 데이터 검색캐싱에 추가
// TODO DOING 디스코드 웹훅
func init() {
	ch := make(chan struct{})
	config.LoadConfig()
	src.StartIndex()
	db.ConnectMaria()
	go db.LoadIndex()
	go db.LoadCache()
	go banchoCroller.LoadBancho(ch)
	_ = <-ch
	//go banchoCroller.RunGetBeatmapDataASBancho()

	if os.Getenv("debug") != "true" {

		go banchoCroller.RunGetBeatmapDataASBancho()
	} else {
	}
	//go banchoCroller.UpdateAllPackList()
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

	e.Pre(
		middleware.RemoveTrailingSlash(),
	)

	e.Use(
		middleware.Logger(),
		//middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		//middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
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
			"CpuThreadCount":        runtime.NumCPU(),
			"RunningGoroutineCount": runtime.NumGoroutine(),
			"apiCount":              *banchoCroller.ApiCount,
		})
	})

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:seIid", Route.DownloadBeatmapSet)
	e.GET("/beatmap/:mapId", Route.DownloadBeatmapSet)
	e.GET("/beatmapset/:seIid", Route.DownloadBeatmapSet)
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
	//e.HTTPErrorHandler = httpErrorHandler.HttpErrorHandler
	//e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
	e.GET("/test", func(c echo.Context) error {
		fmt.Println(time.Now().Format("15:04:05.000"), "test")
		return c.String(200, "/test")
	},

	//middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
	//	Rate:      6,
	//	Burst:     6,
	//	ExpiresIn: time.Minute * 10,
	//})),
	//middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
	//
	//	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
	//		middleware.RateLimiterMemoryStoreConfig{
	//			// 1 / 원하는 대기시간?
	//			//   0.5 = 2s / 0.1 = 10s / 0.01 = 100s /
	//			Rate:  1,
	//			Burst: 0,
	//			//ExpiresIn: time.Millisecond, // 이게 있던없던 그냥 1초당 제한인듯함
	//		},
	//	),
	//}),
	)
	e.GET("/ws", hello)
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
	//webhook.DiscordError(&bodyStruct.ErrorStruct{
	//	Code:      "1",
	//	Path:      "2",
	//	RequestId: "3",
	//	Error:     "4",
	//	Message:   "5",
	//})
	pterm.Info.Println("ECHO STARTED AT", config.Config.Port)
	e.Logger.Fatal(e.Start(":" + config.Config.Port))

}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func hello(c echo.Context) (err error) {
	fmt.Println("wss")
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		pterm.Error.Println(err)
		return err
	}
	defer ws.Close()
	for {
		// Write
		err = ws.WriteMessage(websocket.TextMessage, Route.Search2(Route.SearchQuery{}))
		if err != nil {
			if !strings.Contains(err.Error(), "websocket: close") {
				pterm.Error.Println(err)
				c.Logger().Error(err)

			}
			break
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "websocket: close") {
				pterm.Error.Println(err)
				c.Logger().Error(err)

			}
			break

		}
		fmt.Printf("%s\n", msg)
	}
	return
}
