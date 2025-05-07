package lastfm

// https://www.last.fm/api/show/auth.getSession

type SessionParams struct {
	// The session key returned by the AuthGetToken method.
	Token string `url:"token"`
}

// https://www.last.fm/api/show/auth.getMobileSession

type MobileSessionParams struct {
	// The username or email address of the user to authenticate.
	Username string `url:"username"`
	// The plain text password of the user to authenticate.
	Password string `url:"password"`
}

type Session struct {
	Name       string  `xml:"name"`
	Key        string  `xml:"key"`
	Subscriber IntBool `xml:"subscriber"`
}
