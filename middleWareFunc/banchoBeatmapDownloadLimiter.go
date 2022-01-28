package middleWareFunc

import "github.com/labstack/echo/v4"

func BanchoBeatmapDownloadLimiter(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(c)
	}
}
