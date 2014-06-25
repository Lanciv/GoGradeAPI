package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	// 	"log"
	"strings"
	"time"
)

var (
	ErrInvalidPassword = errors.New("Username or password incorrect.")
)

type User struct {
	Id             int64
	Email          string
	EmailLower     string `db:"email_lower"`
	HashedPassword []byte `db:"hashed_password"`
	Role           string
	Disabled       bool
	TimeStamp
}

// Create a token for the user after we verified their password.
// TODO: Store this in a db? This would be helpful if we would like to invalidate a login.
func (a *User) CreateToken() (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	token.Claims["Id"] = a.Id
	token.Claims["Email"] = a.Email
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// TODO: Move this to a config file.
	tokenString, err := token.SignedString([]byte("someRandomSigningKey"))
	if err != nil {
		return "", err
	}
	// Sould we just return a Token instead of a string???
	return tokenString, err
}

func NewUser(email, role string) *User {

	emailLower := strings.ToLower(email)

	user := User{
		Email:      email,
		EmailLower: emailLower,
		Role:       role,
	}

	return &user
}

func (u *User) SetPassword(password string) error {
	// Password validation.
	switch {
	case len(password) < 6:
		return ErrInvalidPassword
	case len(password) > 265:
		return ErrInvalidPassword
	}
	// Hash password
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//
	u.HashedPassword = b

	return nil
}
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}