package main

import (
	"github.com/Nerinyan/Nerinyan-APIV2/Logger"
	"github.com/Nerinyan/Nerinyan-APIV2/banchoCroller"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/Nerinyan/Nerinyan-APIV2/middlewareFunc"
	"github.com/Nerinyan/Nerinyan-APIV2/route"
	"github.com/Nerinyan/Nerinyan-APIV2/src"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"runtime"
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
// TODO DOING /status에 들어갈 상태값 추가.

func init() {
	ch := make(chan struct{})
	config.LoadConfig()
	src.StartIndex()
	db.ConnectRDBMS()
	go banchoCroller.LoadBancho(ch)
	_ = <-ch

	if config.Config.Debug {
		//go banchoCroller.UpdateAllPackList()
	} else {
		go banchoCroller.RunGetBeatmapDataASBancho()
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = middlewareFunc.CustomHTTPErrorHandler

	e.Renderer = &route.Renderer

	go func() {
		for {
			<-logger.Ch
			e.Logger.SetOutput(log.Writer())
			pterm.Info.Println("UPDATED ECHO LOGGER.")
		}
	}()

	e.Pre(
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)),

		middleware.RemoveTrailingSlash(),
		middleware.Logger(),

		middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{echo.GET, echo.HEAD}}),
		//middleware.RateLimiterWithConfig(middleWareFunc.RateLimiterConfig),
		middleware.RequestID(),
		middleware.Recover(),
		//func(next echo.HandlerFunc) echo.HandlerFunc {
		//	return func(c echo.Context) (err error) {
		//		start := time.Now()
		//		if err = next(c); err != nil {
		//			c.Error(err)
		//		}
		//		stop := time.Now()
		//
		//		c.Response().Before(func() {
		//			_ = time.Now().UnixMilli()                           //	"time": "${time_custom}",
		//			_ = c.Request().Header.Get(echo.HeaderXRequestID)    //	"id": "${id}",
		//			_ = c.Response().Header().Get(echo.HeaderXRequestID) //	"id": "${id}",
		//			_ = c.RealIP()                                       //	"remote_ip": "${remote_ip}",
		//			_ = c.Request().Host                                 //	"host": "${host}",
		//			_ = c.Request().RequestURI                           //	"uri": "${uri}",
		//			_ = c.Request().Method                               //	"method": "${method}",
		//			_ = c.Request().URL.Path                             //	"uri": "${uri}",
		//			_ = c.Request().UserAgent()                          //	"user_agent": "${user_agent}",
		//			_ = c.Response().Status                              //	"status": "${status}",
		//			_ = func() string {
		//				if err != nil {
		//					// Error may contain invalid JSON e.g. `"`
		//					b, _ := json.Marshal(err.Error())
		//					return string(b[1 : len(b)-1])
		//				}
		//				return ""
		//			}()                                               //	"error": "${error}",
		//			_ = strconv.FormatInt(int64(stop.Sub(start)), 10) //	"latency": "${latency}",
		//			_ = stop.Sub(start).String()                      //	"latency_human": "${latency_human}",
		//			_ = func() string {
		//				cl := c.Request().Header.Get(echo.HeaderContentLength)
		//				if cl == "" {
		//					cl = "0"
		//				}
		//				return cl
		//			}()                                          //	"bytes_in": "${bytes_in}",
		//			_ = strconv.FormatInt(c.Response().Size, 10) //	"bytes_out": "${bytes_out}"
		//
		//			//=========================================================================
		//			c.Request()
		//
		//		})
		//		return next(c)
		//	}
		//},
	)

	// docs ============================================================================================================
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, `https://nerinyan.stoplight.io/docs/nerinyan-api`)
	})

	// 서버상태 체크용 ====================================================================================================

	e.GET("/health", route.Health)
	e.GET("/robots.txt", route.Robots)
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"CpuThreadCount":        runtime.NumCPU(),
			"RunningGoroutineCount": runtime.NumGoroutine(),
			"apiCount":              *banchoCroller.ApiCount,
		})
	})

	// 맵 파일 다운로드 ===================================================================================================
	e.GET("/d/:setId", route.DownloadBeatmapSet, route.Embed)
	e.GET("/beatmap/:mapId", route.DownloadBeatmapSet)
	e.GET("/beatmapset/:setId", route.DownloadBeatmapSet)
	//TODO 맵아이디, 맵셋아이디 지원

	// 비트맵 BG  =========================================================================================================
	e.GET("/bg/:setId", func(c echo.Context) error {
		redirectUrl := "https://subapi.nerinyan.moe/bg/" + c.Param("setId")
		return c.Redirect(http.StatusPermanentRedirect, redirectUrl)
	})

	// 비트맵 리스트 검색용 ================================================================================================
	e.GET("/search", route.Search)

	// 개발중 || 테스트중 ===================================================================================================
	e.GET("/test", func(c echo.Context) error {
		return echo.NewHTTPError(400, "test bad request")
	})

	// ====================================================================================================================
	pterm.Info.Println("ECHO STARTED AT", config.Config.Port)
	e.Logger.Fatal(e.Start(":" + config.Config.Port))

}
