package server

import (
	"github.com/labstack/echo/v4"
	"github.com/philips-software/go-hsdp-signer"
	"net/http"
)

func HSDPValidator(key, secret string) echo.MiddlewareFunc {
	s, err := signer.New(key, secret)
	if err != nil {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return echo.ErrUnauthorized
			}
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			valid, err := s.ValidateRequest(req)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  err.Error(),
					Internal: err,
				}
			} else if valid {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}
