package main

import (
	"embed"
	_ "embed"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var store = map[string][]time.Time{}
var block = map[string]time.Time{}
var limit = struct {
	Count int
	Time  time.Duration
}{
	Count: 100,
	Time:  time.Minute,
}

//go:embed * embed.html
var file embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	e := echo.New()
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			for i, t := range store[ip] {
				if t.After(time.Now().Add(-time.Minute)) { // 지금으로부터 1분 이전것 모두 삭제
					store[ip] = store[ip][i:]
					break
				}
			}

			if block[ip].After(time.Now()) || len(store[ip]) > 100 { // 차단해제시간 > 지금 || len(store[ip]) > 100
				block[ip] = time.Now().Add(time.Minute * 10)
				return c.String(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)+" do not request at "+block[ip].Format(time.RFC3339))
			}

			store[ip] = append(store[ip], time.Now())
			c.Response().Header().Add("X-Request-Limit", strconv.Itoa(100-len(store[ip])))
			return next(c)
		}
	})

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFS(file, "*.html")),
	}
	e.Renderer = renderer

	e.GET("/2", func(c echo.Context) error {
		return c.Render(http.StatusOK, "embed.html", map[string]interface{}{
			"setId": 123456,
			"title": "title",
			"content": `Ranked osu! beatmap by Kyuukai.
⏫ Ranked · 📚 9 Difficulties · 🎵 190 · ❤️ 601
(osu!) Normal - ⭐ 2.07 · ⏳ 2:17 | CS 3.2 · AR 4.5
(osu!) Shogun&#39;s Hard - ⭐ 3.26 · ⏳ 2:17 | CS 3.5 · AR 8
(taiko) qoot8123&#39;s Muzukashii - ⭐ 3.48 · ⏳ 2:19 | CS 3.7 · AR 9
(taiko) 29&#39;s Oni - ⭐ 4.59 · ⏳ 2:19 | CS 6 · AR 4.5
(osu!) Meftly&#39;s Insane - ⭐ 4.62 · ⏳ 2:19 | CS 4 · AR 9
(osu!) Xenokai&#39;s Insane - ⭐ 5.18 · ⏳ 2:19 | CS 3.8 · AR 9.1
(taiko) Charlotte&#39;s Inner Oni - ⭐ 5.45 · ⏳ 2:17 | CS 5 · AR 4.5
(osu!) den0saur&#39;s Extra - ⭐ 5.69 · ⏳ 2:19 | CS 4 · AR 9.2
(osu!) Comet - ⭐ 6.18 · ⏳ 2:19 | CS 3.7 · AR 9.4`,
		})
		//return c.HTML(http.StatusOK, file)
	})
	go func() {
		time.Sleep(1000 * time.Millisecond)
		res, err := http.Get("http://localhost/error")
		if err != nil {
			log.Println(err)
			return
		}
		defer res.Body.Close()
		log.Println("response", res.StatusCode)
	}()
	log.Fatalln(e.Start(":80"))
}
