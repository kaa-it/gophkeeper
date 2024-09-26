package domain

// Credentials represents a user's login credentials including their login, password, and additional metadata.
type Credentials struct {
	Login    string
	Password string
	Metadata string
}

// NewCredentials creates a new instance of the Credentials type with the provided login, password, and metadata.
func NewCredentials(login, password, metadata string) *Credentials {
	return &Credentials{
		Login:    login,
		Password: password,
		Metadata: metadata,
	}
}

// Clone creates and returns a deep copy of the Credentials instance.
func (c *Credentials) Clone() *Credentials {
	return &Credentials{
		Login:    c.Login,
		Password: c.Password,
		Metadata: c.Metadata,
	}
}
