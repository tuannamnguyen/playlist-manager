package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	Scopes []string
}

func NewValidator(scopes ...string) *Validator {
	return &Validator{Scopes: scopes}
}

func (v *Validator) CheckTokenHasScopes(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

		claims := token.CustomClaims.(*CustomClaims)
		for _, scope := range v.Scopes {
			if !claims.HasScope(scope) {
				return c.JSON(http.StatusForbidden, []byte(`{"message":"Insufficient scope."}`))
			}
		}

		return next(c)
	}
}
