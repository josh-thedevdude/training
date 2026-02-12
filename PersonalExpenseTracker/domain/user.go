package domain

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    string
	UpdatedAt    string
}

func NewUser(email, password, role string) *User {
	return &User{
		Email:        email,
		PasswordHash: password,
		Role:         role,
	}

}
