package resources

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
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
	userJSON := `{"name": "Sarin", "email": "sarink@juniper.net", "stripe_token": "tok_19xNPoCuKpXpIhNhPcO4WFs9"}`
	// Setup
	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/users", strings.NewReader(userJSON))
	if assert.NoError(t, err) {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ec := e.NewContext(req, rec)
		controller := &Controller{DB: ConnectDB()}

		if assert.NoError(t, controller.createUser(ec)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// // Body has the user ID. for now matching the length of UserID less leading and trailing quote
			// assert.Equal(t, 8, len(rec.Body.String())-2)
		}
	}
}
