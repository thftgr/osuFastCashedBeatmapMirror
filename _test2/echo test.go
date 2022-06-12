package main

import (
	_ "embed"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//X-RateLimit-Limit            //간격당 허용되는 호출 수입니다.
//X-RateLimit-Remaining        //제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.
//X-Burst-RateLimit-Remaining  //버스트 제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.

//X-Burst-RateLimit-Reset      //다음 간격이 시작될 때까지 남은 시간(초)입니다.
//X-Retry-After                //와 동일합니다 X-RateLimit-Reset.

const (
	RATELIMIT_LIMIT           = "X-RateLimit-Limit"           //간격당 허용되는 호출 수입니다.
	RATELIMIT_REMAINING       = "X-RateLimit-Remaining"       //제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.
	BURST_RATELIMIT_REMAINING = "X-Burst-RateLimit-Remaining" //버스트 제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.
	BURST_RATELIMIT_RESET     = "X-Burst-RateLimit-Reset"     //다음 버스트 리셋이 시작될 때까지 남은 시간(초)입니다.
	RETRY_AFTER               = "X-Retry-After"               //와 동일합니다 X-RateLimit-Reset.
)

var mutex = sync.Mutex{}
var store = struct {
	history map[string][]time.Time
	block   map[string]time.Time
	burst   map[string]time.Time
}{
	history: map[string][]time.Time{},
	block:   map[string]time.Time{},
	burst:   map[string]time.Time{},
}

var limitCount = 5
var burstCount = 10

func threadSafe(f func()) {
	mutex.Lock()
	defer mutex.Unlock()
	f()

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	e := echo.New()
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			now := time.Now()
			//버스트 리셋 10분
			if store.burst[ip].After(now.Add(-time.Minute * 10)) {
				threadSafe(func() {
					delete(store.burst, ip)
				})
			}

			// 일반 요청 제한
			for i, t := range store.history[ip] {
				if t.After(now.Add(-time.Minute)) { // 지금으로부터 1분 이전것 모두 삭제
					threadSafe(func() {
						store.history[ip] = store.history[ip][i:]
					})
					break
				}
			}

			var (
				requestCount   = len(store.history[ip]) + 1
				limitRemaining = limitCount - requestCount
				burstRemaining = burstCount - requestCount
			)

			c.Response().Header().Add(RATELIMIT_LIMIT, strconv.Itoa(limitCount))               //간격당 허용되는 호출 수입니다.
			c.Response().Header().Add(RATELIMIT_REMAINING, strconv.Itoa(limitRemaining))       //제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.
			c.Response().Header().Add(BURST_RATELIMIT_REMAINING, strconv.Itoa(burstRemaining)) //버스트 제한에 도달하기 전에 간격에 남아 있는 호출 수입니다.

			c.Response().Header().Add(BURST_RATELIMIT_RESET, "") //다음 버스트 리셋이 시작될 때까지 남은 시간(초)입니다.

			//           차단여부                           버스트 후 시간이 10분 경과했는지
			if store.block[ip].After(now) || (store.burst[ip].After(now) && requestCount > burstCount) {
				threadSafe(func() {
					store.block[ip] = now.Add(time.Minute * 10)
				})

				c.Response().Header().Add(RETRY_AFTER, "600") //와 동일합니다 X-RateLimit-Reset.
				return c.String(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)+". do not request at "+store.block[ip].Format(time.RFC3339))
			}

			threadSafe(func() {
				delete(store.block, ip) // 차단되지 않았다면 메모리에서 해제
				store.history[ip] = append(store.history[ip], now)
			})

			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	log.Fatalln(e.Start(":80"))
}
