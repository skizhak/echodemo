package resources

import (
	"github.com/labstack/echo"
)

// RunServer starts echo server
func RunServer(c ControllerInterface) {
	e := echo.New()

	// Routes
	e.Static("/", "web/dist")
	register(e, c.routes())

	e.Logger.Fatal(e.Start(":9001"))
}

func register(e *echo.Echo, routes map[string]*Handler) {
	for path, handler := range routes {
		if handler.GET != nil {
			e.GET(path, handler.GET)
		}
		if handler.POST != nil {
			e.POST(path, handler.POST)
		}
	}
}
