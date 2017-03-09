package resources

import (
	"math/rand"
	"net/http"
	"time"

	"fmt"

	"github.com/labstack/echo"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
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

	// Payment object
	Payment struct {
		*stripe.ChargeParams
		Charge *stripe.Charge
	}
)

//Update your Stripe Test key.
const stripeKey = ""

func (c *Controller) sayHello(ec echo.Context) error {
	return ec.String(http.StatusOK, "Hello Go Echo framework!")
}

func (c *Controller) getUsers(ec echo.Context) error {
	return ec.JSON(http.StatusOK, c.users)
}

func (c *Controller) getUser(ec echo.Context) error {
	id := ec.Param("id")
	// u, err := c.users.Find(id)
	u := c.users[id]
	if u == nil {
		return ec.JSON(http.StatusOK, "Not Found")
	}
	return ec.JSON(http.StatusOK, u)
}

func generateID(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (c *Controller) createUser(ec echo.Context) error {
	u := &User{}
	stripe.Key = stripeKey
	if err := ec.Bind(u); err != nil {
		return err
	}
	if u.AccountID == "" {
		u.AccountID = generateID(8)
	}
	u.ID = generateID(8)
	ec.Logger().Info("Creating User: ", u.ID)
	stripeCustomerParams := &stripe.CustomerParams{
		Email: u.Email,
	}
	stripeCustomerParams.SetSource(u.StripeToken)
	stripeCustomer, err := customer.New(stripeCustomerParams)
	u.StripeID = stripeCustomer.ID
	if err != nil {
		return err
	}
	// c.users.Insert(u)
	if c.users == nil {
		c.users = map[string]*User{}
	}
	c.users[u.ID] = u
	fmt.Println(c.users)
	return ec.JSON(http.StatusCreated, u)
}

func (c *Controller) makePaymentModel() *Payment {
	pModel := &Payment{}
	return pModel
}

func (c *Controller) getPayment(ec echo.Context) error {
	id := ec.Param("id")
	u, err := c.users.Find(id)
	if err {
		return ec.JSON(http.StatusOK, "Not Found")
	}
	return ec.JSON(http.StatusOK, u.Payments)
}

func (c *Controller) postPayment(ec echo.Context) error {
	id := ec.Param("id")
	u := c.users[id]
	if u == nil {
		fmt.Println("User ID not found!!")
		return nil
	}
	fmt.Println("Making Payment for user: ", id)
	model := c.makePaymentModel()
	if err := ec.Bind(model); err != nil {
		fmt.Println("Bind error: ", err)
		return nil
	}
	// fmt.Println("Bill: ", model.Amount, model.Currency)
	chargeParams := &stripe.ChargeParams{
		Amount:   model.Amount,
		Currency: model.Currency,
		Customer: u.StripeID,
	}
	charge, error := charge.New(chargeParams)
	// fmt.Println("Charged user: ", charge)
	model.Charge = charge
	if error != nil {
		return ec.JSON(http.StatusUnauthorized, charge)
	}
	if u.Payments == nil {
		u.Payments = map[string]*Payment{}
	}
	u.Payments[charge.ID] = model
	return ec.JSON(http.StatusOK, charge)
}

// Routes
func (c *Controller) routes() map[string]*Handler {
	return map[string]*Handler{
		"/api/hello": &Handler{
			GET: c.sayHello,
		},
		"/users": &Handler{
			GET:  c.getUsers,
			POST: c.createUser,
		},
		"/users/:id": &Handler{
			GET: c.getUser,
		},
		"/users/:id/payments": &Handler{
			GET:  c.getPayment,
			POST: c.postPayment,
		},
	}
}
