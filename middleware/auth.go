package middleware

import (
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromSession(c)
		if err != nil {
			// TODO: Log error
			c.Set("CurrentUser", &models.User{})
		} else {
			c.Set("CurrentUser", user)
		}
		c.Next()
	}
}

// GetUserFromSession returns the User for the current session.
func GetUserFromSession(c *gin.Context) (user *models.User, err error) {
	ctx := common.Context()
	if cookie, err := c.Cookie(ctx.Config.Cookies.SessionCookie); err == nil {
		value := ""
		if err = ctx.Config.Cookies.Secure.Decode(ctx.Config.Cookies.SessionCookie, cookie, &value); err != nil {
			return nil, common.ErrDecodeCookie
		}
		userID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, common.ErrWrongDataType
		}
		user, err = models.UserGet(userID)
	}
	return user, err
}
