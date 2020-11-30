package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func HSDPValidator(username, password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authErr := fmt.Errorf("unauthorized")
			req := c.Request()
			auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
			if len(auth) != 2 || auth[0] != "Basic" {
				http.Error(c.Response(), "authorization failed", http.StatusUnauthorized)
				return authErr
			}
			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 || !(pair[0] == username && pair[1] == password) {
				http.Error(c.Response(), "authorization failed", http.StatusUnauthorized)
				return authErr
			}
			return next(c)

		}
	}
}
