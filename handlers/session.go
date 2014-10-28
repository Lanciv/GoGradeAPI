package handlers

import (
	// "log"

	m "github.com/Lanciv/GoGradeAPI/model"
	s "github.com/Lanciv/GoGradeAPI/store"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
)

// ErrLoginFailed ...
var ErrLoginFailed = "Login Failed! Email and/or password incorrect."

// LoginForm ...
type LoginForm struct {
	Email    string
	Password string
}

// FieldMap ...
func (lf *LoginForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&lf.Email:    binding.Field{Form: "email", Required: true},
		&lf.Password: binding.Field{Form: "password", Required: true},
	}
}

// Login ...
func Login(c *gin.Context) {
	// Get username and password
	lf := new(LoginForm)

	errs := binding.Bind(c.Request, lf)
	if errs != nil {
		writeError(c.Writer, errs, 400, nil)
		return
	}

	user, err := s.Users.FindByEmail(lf.Email)
	if err != nil {
		writeError(c.Writer, ErrLoginFailed, 401, err)
		return
	}

	if err := user.ComparePassword(lf.Password); err != nil {
		writeError(c.Writer, ErrLoginFailed, 401, nil)
		return
	}

	// Create a session for the user.
	session, err := m.NewSession(user)
	if err != nil {
		writeError(c.Writer, serverError, 500, err)
		return
	}

	s.Sessions.Store(&session)
	// Send token to the user so they can use it to to authenticate all further requests.
	c.JSON(200, &APIRes{"session": []m.Session{session}})
	return
}

// AuthRequired ...
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := jwt.ParseFromRequest(c.Request, func(t *jwt.Token) (interface{}, error) {
			return []byte("someRandomSigningKey"), nil
		})
		if err != nil {
			// c.JSON(200,http.StatusUnauthorized, "Access denied.")
			writeError(c.Writer, "Unauthorized", 401, nil)
			c.Fail(401, err)
			return
		}
		c.Set("userId", res.Claims["userId"])
	}
}
