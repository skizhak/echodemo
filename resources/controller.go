package resources

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/labstack/echo"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"

	"github.com/skizhak/echodemo/resources/utils"
)

type (
	// ControllerInterface
	ControllerInterface interface {
		routes() map[string]*Handler
	}

	// Controller
	Controller struct {
		DB     *sql.DB
		Logger echo.Logger
		users  UserMap
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

	// Bill
	Bill struct {
		ID            string
		Name          string
		Description   string
		InvoiceID     string
		PaymentMethod string
		AccountID     string
		Date          string
	}
)

//Update your Stripe Test key.
const stripeKey = "sk_test_LAp101ps8FLj5NsDb1dULGny"

func (c *Controller) sayHello(ec echo.Context) error {
	return ec.String(http.StatusOK, "Hello Go Echo framework!")
}

func (c *Controller) getUsers(ec echo.Context) error {
	var users []*User
	tx, err := c.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Query("SELECT id, name, description, email, account_id FROM users")
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Description, &user.Email, &user.AccountID)
		users = append(users, &user)
		if err != nil {
			log.Fatal(err)
		}
	}
	if rows.Err() != nil {
		log.Fatal("Error iterating rows on users table", rows.Err())
	}
	tx.Commit()
	return ec.JSON(http.StatusOK, users)
}

func (c *Controller) getUser(ec echo.Context) error {
	id := ec.Param("id")
	var user User
	err := c.DB.QueryRow("SELECT id, name, description, email, account_id FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Description, &user.Email, &user.AccountID)
	if err != nil && err == sql.ErrNoRows {
		log.Fatal("User not found.", err)
		return ec.JSON(http.StatusFailedDependency, "User does not exist.")
	}
	return ec.JSON(http.StatusOK, user)
}

func (c *Controller) createUser(ec echo.Context) error {
	u := &User{}
	stripe.Key = stripeKey
	if err := ec.Bind(u); err != nil {
		return err
	}
	// Start a transaction.
	tx, err := c.DB.Begin()
	var stmt *sql.Stmt
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	u.ID = utils.GenerateID(8)
	u.Password = utils.GenerateID(10)
	u.Description = "User Account"
	// If no AccountID is present, create one.
	if u.AccountID == "" {
		u.AccountID = utils.GenerateID(8)
		u.Description = "Admin Account"
		stmt = c.prepare("INSERT INTO accounts VALUES (?, ?, ?, ?)", tx)
		_, err = stmt.Exec(u.AccountID, u.Name, u.Description, "active")
		if err != nil {
			log.Fatal("Insert into accounts table failed.", err)
		}
	} else {
		var accountStatus string
		err := c.DB.QueryRow("SELECT status from accounts where account_id = ?", u.AccountID).Scan(&accountStatus)
		if err != nil && err == sql.ErrNoRows {
			log.Fatal("Asosciated account not found.", err)
			return ec.JSON(http.StatusFailedDependency, "Account ID does not exist.")
		}
		if accountStatus != "active" {
			log.Fatal("Asosciated account is not Active.", err)
			return ec.JSON(http.StatusFailedDependency, "Account ID not active.")
		}
	}
	log.Println("Creating User: ", u.ID)
	stmt = c.prepare("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?)", tx)
	// userSlice := []interface{}{u.ID, u.Name, u.Description, u.Email, u.Password, u.AccountID}
	_, err = stmt.Exec(u.ID, u.Name, u.Description, u.Email, u.Password, u.AccountID)
	if err != nil {
		log.Fatal("Insert into users table failed.", err)
		return ec.JSON(http.StatusInternalServerError, "Retry Later.")
	}

	// Stripe API
	stripeCustomerParams := &stripe.CustomerParams{
		Email: u.Email,
	}
	stripeCustomerParams.SetSource(u.StripeToken)
	stripeCustomer, err := customer.New(stripeCustomerParams)
	if err != nil {
		log.Fatal("Stripe Customer creation failed", err)
	}

	stmt = c.prepare("INSERT INTO payment_methods(id, name, description, token, account_id) VALUES (?, ?, ?, ?, ?)", tx)
	_, err = stmt.Exec(stripeCustomer.ID, "Stripe", stripeCustomer.Desc, u.StripeToken, u.AccountID)
	if err != nil {
		fmt.Println("Insert into payment_methods table failed.")
		log.Fatal(err)
	}
	stmt.Close()

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return ec.JSON(http.StatusCreated, u.ID)
}

func (c *Controller) makePaymentModel() *Payment {
	pModel := &Payment{}
	return pModel
}

func (c *Controller) getPayment(ec echo.Context) error {
	id := ec.Param("id")
	var accountID string
	var payments []*Bill
	// Get AccountID for the user
	err := c.DB.QueryRow("SELECT account_id FROM users WHERE id = ?", id).Scan(&accountID)
	if err != nil {
		fmt.Println("UserID doesnt exist")
		log.Fatal(err)
	}
	tx, err := c.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := tx.Query("SELECT id, name, description, account_id, date FROM payment_historys WHERE account_id = ?", accountID)
	defer rows.Close()
	for rows.Next() {
		var payment Bill
		err := rows.Scan(&payment.ID, &payment.Name, &payment.Description, &payment.AccountID, &payment.Date)
		payments = append(payments, &payment)
		if err != nil {
			log.Fatal(err)
		}
	}
	if rows.Err() != nil {
		log.Fatal("Error iterating rows on payment_historys table", rows.Err())
	}
	tx.Commit()
	return ec.JSON(http.StatusOK, payments)
}

func (c *Controller) postPayment(ec echo.Context) error {
	var accountID string
	var stripeID string

	id := ec.Param("id")

	// Get AccountID for the user
	err := c.DB.QueryRow("SELECT account_id FROM users WHERE id = ?", id).Scan(&accountID)
	if err != nil {
		fmt.Println("UserID doesnt exist")
		log.Fatal(err)
	}

	// Get Stripe ID for the account
	err = c.DB.QueryRow("SELECT id FROM payment_methods WHERE name = 'Stripe' AND account_id = ?", accountID).Scan(&stripeID)
	if err != nil {
		fmt.Println("Stripe customer ID not found.")
		log.Fatal(err)
	}

	model := c.makePaymentModel()
	if err := ec.Bind(model); err != nil {
		fmt.Println("Bind error: ", err)
		return nil
	}
	fmt.Println("Making Payment for user: ", id, model.Amount, model.Currency)
	chargeParams := &stripe.ChargeParams{
		Amount:   model.Amount,
		Currency: model.Currency,
		Customer: stripeID,
	}
	charge, error := charge.New(chargeParams)
	if error != nil {
		return ec.JSON(http.StatusUnauthorized, charge)
	}
	tx, err := c.DB.Begin()
	stmt := c.prepare("INSERT INTO payment_historys(id, name, description, account_id, date) VALUES (?, ?, ?, ?, ?)", tx)
	_, err = stmt.Exec(charge.Customer.ID, "Stripe", charge.Desc, accountID, time.Unix(charge.Created, 0).String())
	if err != nil {
		log.Fatal("Insert into payment_historys table failed.", err)
	}
	stmt.Close()
	tx.Commit()

	return ec.JSON(http.StatusOK, charge)
}

func (c *Controller) prepare(query string, tx *sql.Tx) *sql.Stmt {
	var (
		stmt *sql.Stmt
		err  error
	)
	if tx != nil {
		stmt, err = tx.Prepare(query)
	} else {
		stmt, err = c.DB.Prepare(query)
	}
	if err != nil {
		log.Fatal(err)
	}
	return stmt
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

// ControllerService instance of controller.
func ControllerService() *Controller {
	return &Controller{
		DB: ConnectDB(),
	}
}
