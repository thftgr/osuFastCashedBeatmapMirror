package Route

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Wiki(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect,"https://thftgr.stoplight.io/docs/beatmap-mirror")
}
