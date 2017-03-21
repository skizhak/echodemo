package resources

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/skizhak/echodemo/resources/payments"
	"github.com/skizhak/echodemo/resources/utils"
	"github.com/stretchr/testify/assert"
	stripe "github.com/stripe/stripe-go"
)

type MockControllerService struct {
	onCreateUser  func(ec *echo.Context) error
	onGetUser     func(ec *echo.Context) error
	onGetUsers    func(ec *echo.Context) error
	onPostPayment func(ec *echo.Context) error
	onGetPayment  func(ec *echo.Context) error
}

// Implement ControllerService interface
func (t *MockControllerService) CreateUser(ec *echo.Context) error {
	return t.onCreateUser(ec)
}

func (t *MockControllerService) GetUser(ec *echo.Context) error {
	return t.onGetUser(ec)
}

func (t *MockControllerService) GetUsers(ec *echo.Context) error {
	return t.onGetUsers(ec)
}

func (t *MockControllerService) PostPayment(ec *echo.Context) error {
	return t.onPostPayment(ec)
}

func (t *MockControllerService) GetPayment(ec *echo.Context) error {
	return t.onGetPayment(ec)
}

// End of ControllerService Implement

func TestCreateUser(t *testing.T) {
	// assert := assert.New(t)
	userJSON := `{
		"name": "Sarin",
		"email": "sarink@juniper.net",
		"payment_service": "Stripe",
		"payment_token": "tok_19zZMoCuKpXpIhNh7HF9qrrZ"
	}`
	// Setup
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/users", strings.NewReader(userJSON))
	if assert.NoError(t, err) {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ec := e.NewContext(req, rec)
		// var testController resources.ControllerService
		testController := &Controller{
			DB: ConnectDB(),
			PS: NewMockService("Stripe"),
		}

		if assert.NoError(t, testController.createUser(ec)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			// Body has the user ID. for now matching the length of UserID less leading and trailing quote
			assert.Equal(t, 8, len(rec.Body.String())-2)
		}
	}
}

func NewMockService(service string) payments.PaymentService {
	var ps payments.PaymentService
	switch service {
	case "Stripe":
		ps := CreateMockStripeService()
		return ps
	}
	return ps
}

//CreateMockStripeService mocks StripeService
func CreateMockStripeService() *MockStripeService {
	ms := &MockStripeService{
		onGetToken: func(card payments.Card) string {
			return ""
		},
		onCreateCustomer: func(email string, token string) (*payments.C3Customer, error) {
			return &payments.C3Customer{
				ID:      "cus_" + utils.GenerateID(10),
				Email:   email,
				Service: "Stripe",
			}, nil
		},
		onCharge: func(id string, amount uint64, currency string) (*payments.C3Charge, error) {
			return &payments.C3Charge{
				ID:       "ch_" + utils.GenerateID(15),
				Created:  13456576,
				Desc:     "",
				Amount:   amount,
				Currency: currency,
				Customer: &payments.C3Customer{
					ID: id,
				},
				PaymentMethod: "card",
				ReceiptNumber: "",
				Service:       "Stripe",
				Status:        "",
			}, nil
		},
	}
	return ms
}

type MockStripeService struct {
	CustomerParams   *stripe.CustomerParams
	ChargeParams     *stripe.ChargeParams
	onGetToken       func(card payments.Card) string
	onCreateCustomer func(email string, token string) (*payments.C3Customer, error)
	onCharge         func(id string, amount uint64, currency string) (*payments.C3Charge, error)
}

// Make sure MockStripeService satisfies PaymentService
var _ payments.PaymentService = (*MockStripeService)(nil)

func (ms *MockStripeService) GetToken(card payments.Card) string {
	return ms.onGetToken(card)
}

func (ms *MockStripeService) CreateCustomer(email string, token string) (*payments.C3Customer, error) {
	return ms.onCreateCustomer(email, token)
}

func (ms *MockStripeService) Charge(id string, amount uint64, currency string) (*payments.C3Charge, error) {
	return ms.onCharge(id, amount, currency)
}
