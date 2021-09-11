package Route

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

//Health CollectHost godoc
//@Summary Host health checker.
//@Description 서버 상태 체크.
//@Success 200
//@Router /health [get]
func Health(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
