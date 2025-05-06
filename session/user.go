package session

type User struct {
	session *Session
}

// NewUser creates and returns a new User API route.
func NewUser(session *Session) *User {
	return &User{session: session}
}
