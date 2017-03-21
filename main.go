package main

//go:generate sqlboiler mysql

import (
	"github.com/labstack/echo"
	"github.com/skizhak/echodemo/resources"
)

// runServer starts echo server
func runServer(c resources.ControllerService) {
	e := echo.New()

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Routes
	e.Static("/", "web/dist")
	register(e, c.Routes())

	e.Logger.Fatal(e.Start(":9001"))
}

func register(e *echo.Echo, routes map[string]*resources.Handler) {
	for path, handler := range routes {
		if handler.GET != nil {
			e.GET(path, handler.GET)
		}
		if handler.POST != nil {
			e.POST(path, handler.POST)
		}
	}
}

func main() {
	runServer(resources.CreateControllerService())
}
