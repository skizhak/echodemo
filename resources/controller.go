package resources

import (
	"net/http"

	"github.com/labstack/echo"
)

// ControllerInterface
type ControllerInterface interface {
	routes() map[string]*Handler
}

// Controller
type Controller struct{}

// Handler for routes
type Handler struct {
	GET echo.HandlerFunc
}

func (c *Controller) sayHello(context echo.Context) error {
	return context.String(http.StatusOK, "Hello Go Echo framework!")
}

// Routes
func (c *Controller) routes() map[string]*Handler {
	return map[string]*Handler{
		"/api/hello": &Handler{
			GET: c.sayHello,
		},
	}
}
