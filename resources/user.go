package resources

type (
	// User object
	User struct {
		ID          string              `json:"id"`
		Name        string              `json:"name"`
		Email       string              `json:"email"`
		AccountID   string              `json:"account_id"`
		StripeToken string              `json:"stripe_token"`
		StripeID    string              `json:"stripe_id"`
		Payments    map[string]*Payment `json:"payments"`
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
