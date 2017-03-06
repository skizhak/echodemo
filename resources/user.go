package resources

type (
	// User object
	User struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AccountID string `json:"account_id"`
	}
)

// UserMap type points to list of User object
type UserMap map[int]*User

// Find user from users
func (um UserMap) Find(id int) (*User, bool) {
	user, found := um[id]
	return user, found
}

// Insert new user if doesn't exist
func (um UserMap) Insert(u *User) bool {
	_, found := um[u.ID]
	if !found {
		um[u.ID] = u
	}
	return !found
}
