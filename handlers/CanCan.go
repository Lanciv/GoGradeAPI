package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	m "github.com/Lanciv/GoGradeAPI/model"
	"github.com/Lanciv/GoGradeAPI/store"
)

// RoleIn ...
func RoleIn(role string, roles []string) bool {
	for _, r := range roles {
		if role == r {
			return true
		}
	}
	return false
}

// Can ...
func Can(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := c.Get("userId")
		if err != nil && userID == nil {
			c.Fail(401, err)
			return
		}

		id := userID.(string)

		user := m.User{}
		err = store.UserH.One(&user, id)
		if err != nil {
			c.Fail(401, err)
			return
		}

		if !RoleIn(user.Role, roles) {

			c.Fail(401, errors.New("Unauthorized"))
			return
		}

	}
}