package httpErrorHandler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"runtime"
)

func HttpErrorHandler(err error, c echo.Context) {

	for i := -3; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Printf("%s:%d | %s \r\n", file, line, err)
	}

}
