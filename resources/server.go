package resources

import (
	"github.com/labstack/echo"
)

// RunServer starts echo server
func RunServer(c ControllerInterface) {
	e := echo.New()

	e.Static("/", "web/dist")

	for path, handler := range c.routes() {
		e.GET(path, handler.GET)
	}
	e.Logger.Fatal(e.Start(":9001"))
}

func register(e *echo.Echo, routes map[string]*Handler) {

}
