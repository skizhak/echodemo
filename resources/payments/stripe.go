package payments

import (
	"log"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

// StripeService implements PaymentService
type StripeService struct {
	Key            string
	CustomerParams *stripe.CustomerParams
	ChargeParams   *stripe.ChargeParams
}

// Stripe API Key
func (s *StripeService) SetAPIKey(key string) {
	s.Key = key
	stripe.Key = key
}

// GetToken
func (s *StripeService) GetToken(card Card) string {
	return ""
}

// CreateCustomer returns a stripe Customer
func (s *StripeService) CreateCustomer(email string, token string) (*C3Customer, error) {
	s.CustomerParams = &stripe.CustomerParams{
		Email: email,
	}
	s.CustomerParams.SetSource(token)
	sc, err := customer.New(s.CustomerParams)
	if err != nil {
		log.Fatal("Stripe Customer creation failed", err)
	}
	return &C3Customer{
		ID:      sc.ID,
		Service: "Stripe",
	}, err
}

// Charge a stipe Customer
func (s *StripeService) Charge(id string, amount uint64, currency string) (*C3Charge, error) {
	s.ChargeParams = &stripe.ChargeParams{
		Customer: id,
		Amount:   amount,
		Currency: stripe.Currency(currency),
	}
	ch, err := charge.New(s.ChargeParams)
	if err != nil {
		log.Fatal("Stripe charge error.", err)
	}

	return &C3Charge{
		ID:       ch.ID,
		Created:  ch.Created,
		Desc:     ch.Desc,
		Amount:   ch.Amount,
		Currency: string(ch.Currency),
		Customer: &C3Customer{
			ID: ch.Customer.ID,
		},
		PaymentMethod: string(ch.Source.Type),
		ReceiptNumber: ch.ReceiptNumber,
		Service:       "Stripe",
		Status:        ch.Status,
	}, err
}
