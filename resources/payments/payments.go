package payments

// PaymentService
type PaymentService interface {
	GetToken(card Card) string
	CreateCustomer(email string, token string) (*C3Customer, error)
	Charge(id string, amount uint64, currency string) (*C3Charge, error)
}

type C3Customer struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Balance int64  `json:"account_balance"`
	Desc    string `json:"description"`
	Email   string `json:"email"`
	Deleted bool   `json:"deleted"`
	Service string `json:"service"`
}

type C3Charge struct {
	Amount        uint64      `json:"amount"`
	Created       int64       `json:"created"`
	Currency      string      `json:"currency"`
	Customer      *C3Customer `json:"customer"`
	Desc          string      `json:"description"`
	PaymentMethod string      `json:"payment_method"`
	ReceiptNumber string      `json:"receipt_number"`
	ID            string      `json:"id"`
	Paid          bool        `json:"paid"`
	Service       string      `json:"service"`
	Status        string      `json:"status"`
}

type Card struct {
	ID       string      `json:"id"`
	Month    uint8       `json:"exp_month"`
	Year     uint16      `json:"exp_year"`
	LastFour string      `json:"last4"`
	City     string      `json:"address_city"`
	Country  string      `json:"address_country"`
	Address1 string      `json:"address_line1"`
	Address2 string      `json:"address_line2"`
	State    string      `json:"address_state"`
	Zip      string      `json:"address_zip"`
	Customer *C3Customer `json:"customer"`
	Name     string      `json:"name"`
	Deleted  bool        `json:"deleted"`
}

//Update your Stripe Test key.
const stripeKey = "sk_test_LAp101ps8FLj5NsDb1dULGny"

// NewService returns a PaymentService of type service
func NewService(service string) PaymentService {
	var ps PaymentService
	switch service {
	case "Stripe":
		ps := &StripeService{}
		ps.SetAPIKey(stripeKey)
		return ps
	}
	return ps
}

