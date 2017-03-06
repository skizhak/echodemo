package resources

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (
	// ControllerInterface
	ControllerInterface interface {
		routes() map[string]*Handler
	}

	// Controller
	Controller struct {
		users UserMap
	}

	// Handler for routes
	Handler struct {
		GET  echo.HandlerFunc
		POST echo.HandlerFunc
	}
)

func (c *Controller) sayHello(ec echo.Context) error {
	return ec.String(http.StatusOK, "Hello Go Echo framework!")
}

func (c *Controller) getUsers(ec echo.Context) error {
	return ec.JSON(http.StatusOK, c.users)
}

func (c *Controller) getUser(ec echo.Context) error {
	id, _ := strconv.Atoi(ec.Param("id"))
	u, err := c.users.Find(id)
	if err {
		return ec.JSON(http.StatusOK, "Not Found")
	}
	return ec.JSON(http.StatusOK, u)
}

func (c *Controller) createUser(ec echo.Context) error {
	u := &User{}
	if err := ec.Bind(u); err != nil {
		return err
	}
	return nil
}

// Routes
func (c *Controller) routes() map[string]*Handler {
	return map[string]*Handler{
		"/api/hello": &Handler{
			GET: c.sayHello,
		},
		"/users": &Handler{
			GET: c.getUsers,
		},
		"/users/:id": &Handler{
			GET:  c.getUser,
			POST: c.createUser,
		},
	}
}
