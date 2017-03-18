package resources

import _ "github.com/go-sql-driver/mysql"
import "github.com/skizhak/echodemo/resources/payments"

type (
	// User object
	User struct {
		ID             string                        `json:"id"`
		Name           string                        `json:"name"`
		Description    string                        `json:"description"`
		Email          string                        `json:"email"`
		Password       string                        `json:"password"`
		AccountID      string                        `json:"account_id"`
		PaymentService string                        `json:"payment_service"`
		PaymentToken   string                        `json:"payment_token"`
		Payments       map[string]*payments.C3Charge `json:"payments"`
	}
)

// UserMap type points to list of User object
type UserMap map[string]*User

// Find user from users
func (um UserMap) Find(id string) (*User, bool) {
	user, found := um[id]
	return user, found
}

// Insert new user if doesn't exist
func (um UserMap) Insert(u *User) bool {
	found := false
	if um == nil {
		um = map[string]*User{}
	} else {
		_, found = um[u.ID]
	}
	if !found {
		um[u.ID] = u
	}
	return !found
}
