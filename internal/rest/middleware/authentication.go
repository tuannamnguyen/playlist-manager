package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

type ScopeValidator struct {
	Scopes []string
}

func NewScopeValidator(scopes ...string) *ScopeValidator {
	return &ScopeValidator{Scopes: scopes}
}

func (v *ScopeValidator) CheckTokenHasScopes(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

		claims := token.CustomClaims.(*CustomClaims)
		for _, scope := range v.Scopes {
			if !claims.HasScope(scope) {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "Insufficient scope.",
				})
			}
		}

		return next(c)
	}
}
