package resources

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"fmt"

	"github.com/labstack/echo"

	"github.com/skizhak/echodemo/resources/payments"
	"github.com/skizhak/echodemo/resources/utils"
)

type (
	// ControllerService
	ControllerService interface {
		Routes() map[string]*Handler
	}

	// Controller
	Controller struct {
		DB     *sql.DB
		PS     payments.PaymentService
		Logger echo.Logger
		users  UserMap
	}

	// Handler for routes
	Handler struct {
		GET  echo.HandlerFunc
		POST echo.HandlerFunc
	}

	// Account
	Account struct {
		ID          string
		Name        string
		Description string
		Status      string
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

func (c *Controller) createAccount(name string, description string, tx *sql.Tx) (*Account, error) {
	a := &Account{
		ID:          utils.GenerateID(8),
		Name:        name,
		Description: description,
		Status:      "active",
	}
	fmt.Println("Account Creation", a)
	stmt := c.prepare("INSERT INTO accounts VALUES (?, ?, ?, ?)", tx)
	_, err := stmt.Exec(a.ID, a.Name, a.Description, a.Status)
	if err != nil {
		log.Fatal("Insert into accounts table failed.", err)
	}
	defer stmt.Close()
	return a, err
}

func (c *Controller) getAccountStatus(id string) (string, error) {
	var accountStatus string
	err := c.DB.QueryRow("SELECT status from accounts where account_id = ?", id).Scan(&accountStatus)
	if err != nil && err == sql.ErrNoRows {
		log.Fatal("Asosciated account not found.", err)
		return "", err
	}
	return accountStatus, err
}

func (c *Controller) insertUserDB(u *User, tx *sql.Tx) error {
	stmt := c.prepare("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?)", tx)
	// userSlice := []interface{}{u.ID, u.Name, u.Description, u.Email, u.Password, u.AccountID}
	_, err := stmt.Exec(u.ID, u.Name, u.Description, u.Email, u.Password, u.AccountID)
	if err != nil {
		log.Fatal("Insert into users table failed.", err)
	}
	defer stmt.Close()
	return err
}

func (c *Controller) insertPaymentMethodsDB(customer *payments.C3Customer, u *User, tx *sql.Tx) error {
	stmt := c.prepare("INSERT INTO payment_methods(id, name, description, account_id) VALUES (?, ?, ?, ?)", tx)
	_, err := stmt.Exec(customer.ID, customer.Service, customer.Desc, u.AccountID)
	if err != nil {
		fmt.Println("Insert into payment_methods table failed.")
		log.Fatal(err)
	}
	defer stmt.Close()
	return err
}

func (c *Controller) insertPaymentHistoryDB(charge *payments.C3Charge, accountID string, tx *sql.Tx) error {
	stmt := c.prepare("INSERT INTO payment_historys(id, name, description, invoice_id, payment_method, account_id, date) VALUES (?, ?, ?, ?, ?, ?, ?)", tx)
	_, err := stmt.Exec(charge.ID, charge.Service, charge.Desc, charge.Customer.ID, charge.PaymentMethod, accountID, time.Unix(charge.Created, 0).String())
	if err != nil {
		log.Fatal("Insert into payment_historys table failed.", err)
	}
	defer stmt.Close()
	return err
}

func (c *Controller) getUserAccountID(userID string) (string, error) {
	var accountID string
	// Get AccountID for the user
	err := c.DB.QueryRow("SELECT account_id FROM users WHERE id = ?", userID).Scan(&accountID)
	if err != nil {
		fmt.Println("UserID doesnt exist")
		log.Fatal(err)
	}
	return accountID, err
}

func (c *Controller) getPaymentForAccount(accountID string) (string, string, error) {
	var service string
	var paymentID string
	err := c.DB.QueryRow("SELECT id, name FROM payment_methods WHERE account_id = ?", accountID).Scan(&paymentID, &service)
	if err != nil {
		fmt.Println("Stripe customer ID not found.")
		log.Fatal(err)
	}
	return service, paymentID, err
}

func (c *Controller) getPaymentHistroyForAccount(accountID string, tx *sql.Tx) ([]*Bill, error) {
	var payments []*Bill
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
	return payments, err
}

// CreateUser
func (c *Controller) createUser(ec echo.Context) error {
	u := &User{}
	if err := ec.Bind(u); err != nil {
		return err
	}
	// Start a transaction.
	tx, err := c.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	u.ID = utils.GenerateID(8)
	u.Password = utils.GenerateID(10)
	u.Description = "User Account"
	// If no AccountID is present, create one.
	if u.AccountID == "" {
		account, err := c.createAccount("Admin", "", tx)
		if err != nil {
			return ec.JSON(http.StatusInternalServerError, "Retry after sometime.")
		}
		u.AccountID = account.ID
	} else {
		accountStatus, _ := c.getAccountStatus(u.AccountID)
		if accountStatus != "active" {
			log.Fatal("Asosciated account is not Active.", err)
			return ec.JSON(http.StatusFailedDependency, "Account ID not active.")
		}
	}
	log.Println("Creating User: ", u.ID)

	err = c.insertUserDB(u, tx)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, "Retry Later.")
	}

	// Payment API
	// ps := payments.NewService(u.PaymentService)
	// customer, err := ps.CreateCustomer(u.Email, u.PaymentToken)
	customer, err := c.PS.CreateCustomer(u.Email, u.PaymentToken)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, "Retry Later.")
	}
	// Update payment_methods table.
	err = c.insertPaymentMethodsDB(customer, u, tx)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, "Retry Later.")
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return ec.JSON(http.StatusCreated, u.ID)
}

func (c *Controller) getPayment(ec echo.Context) error {
	id := ec.Param("id")
	// Get AccountID for the user
	accountID, _ := c.getUserAccountID(id)
	tx, err := c.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	payments, _ := c.getPaymentHistroyForAccount(accountID, tx)
	tx.Commit()
	return ec.JSON(http.StatusOK, payments)
}

func (c *Controller) postPayment(ec echo.Context) error {
	id := ec.Param("id")

	accountID, _ := c.getUserAccountID(id)
	// Get Payment ID for the account
	_, paymentID, _ := c.getPaymentForAccount(accountID)

	type payment struct {
		Amount   uint64
		Currency string
	}
	pmodel := &payment{}
	if err := ec.Bind(pmodel); err != nil {
		fmt.Println("Bind error: ", err)
		return nil
	}
	// ps := payments.NewService(service)
	// charge, error := ps.Charge(paymentID, pmodel.Amount, pmodel.Currency)
	charge, error := c.PS.Charge(paymentID, pmodel.Amount, pmodel.Currency)
	if error != nil {
		return ec.JSON(http.StatusUnauthorized, charge)
	}
	tx, err := c.DB.Begin()
	err = c.insertPaymentHistoryDB(charge, accountID, tx)
	if err != nil {
		log.Fatal(err)
	}
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
func (c *Controller) Routes() map[string]*Handler {
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

// CreateControllerService create instance of controller which implements ControllerService
func CreateControllerService() *Controller {
	return &Controller{
		DB: ConnectDB(),
		PS: payments.NewService("Stripe"),
	}
}
