package payments

import stripe "github.com/stripe/stripe-go"

type MockStripeService struct {
	CustomerParams   *stripe.CustomerParams
	ChargeParams     *stripe.ChargeParams
	onGetToken       func() string
	onCreateCustomer func() *stripe.Customer
	onCharge         func() *stripe.Charge
}

func (ms *MockStripeService) GetToken() string {
	return ms.onGetToken()
}

func (ms *MockStripeService) CreateCustomer() *stripe.Customer {
	return ms.onCreateCustomer()
}

func (ms *MockStripeService) Charge() *stripe.Charge {
	return ms.onCharge()
}

func CreateMockStripeService(token string, customerParams *stripe.CustomerParams, chargeParams *stripe.ChargeParams) MockStripeService {
	ms := MockStripeService{
		onGetToken: func() string {
			return token
		},
		onCreateCustomer: func() *stripe.Customer {
			return &stripe.Customer{}
		},
		onCharge: func() *stripe.Charge {
			return &stripe.Charge{}
		},
	}
	return ms
}
