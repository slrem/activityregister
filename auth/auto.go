package auth

import (
	"activityregister/tool"
	"errors"

	"github.com/labstack/echo"
)

var sessionError = errors.New("session error")

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.FormValue("uid")
		session := c.FormValue("session")
		if session != tool.MD5("activityregister"+uid+"activityregister") {
			return sessionError
		}
		return next(c)
	}

}
