package sdkcm

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// BeforeFunc defines a function which is executed just before the middleware.
		BeforeFunc BeforeFunc

		// SuccessHandler defines a function which is executed for a valid token.
		SuccessHandler JWTSuccessHandler

		// ErrorHandler defines a function which is executed for an invalid token.
		// It may be used to define a custom JWT error.
		ErrorHandler JWTErrorHandler

		// ErrorHandlerWithContext is almost identical to ErrorHandler, but it's passed the current context.
		ErrorHandlerWithContext JWTErrorHandlerWithContext

		// Signing key to validate token. Used as fallback if SigningKeys has length 0.
		// Required. This or SigningKeys.
		SigningKey interface{}

		// Map of signing keys to validate token with kid field usage.
		// Required. This or SigningKey.
		SigningKeys map[string]interface{}

		// Signing method, used to check token signing method.
		// Optional. Default value HS256.
		SigningMethod string

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey string

		// Claims are extendable claims data defining token content.
		// Optional. Default value jwt.MapClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "param:<name>"
		// - "cookie:<name>"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		KeyFunc jwt.Keyfunc
	}

	// JWTSuccessHandler defines a function which is executed for a valid token.
	JWTSuccessHandler func(*gin.Context)

	// JWTErrorHandler defines a function which is executed for an invalid token.
	JWTErrorHandler func(error) error

	// JWTErrorHandlerWithContext is almost identical to JWTErrorHandler, but it's passed the current context.
	JWTErrorHandlerWithContext func(error, *gin.Context) error

	jwtExtractor func(*gin.Context) (string, error)

	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(*gin.Context) bool

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(*gin.Context)
)

// Algorithms
const (
	AlgorithmHS256      = "HS256"
	HeaderAuthorization = "Authorization"
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(*gin.Context) bool {
	return false
}

// Errors
var (
	ErrJWTMissing = ErrInvalidRequestWithMessage(errors.New("missing or malformed jwt"), "missing or malformed jwt")
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:       DefaultSkipper,
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		TokenLookup:   "header:" + HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
	}
)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
// func JWT(key interface{}) echo.MiddlewareFunc {
// 	c := DefaultJWTConfig
// 	c.SigningKey = key
// 	return JWTWithConfig(c)
// }
// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func JwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		auth := c.Request.Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func JwtFromQuery(param string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

// jwtFromParam returns a `jwtExtractor` that extracts token from the url param string.
func JwtFromParam(param string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		token := c.Param(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from the named cookie.
func JwtFromCookie(name string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", ErrJWTMissing
		}
		return cookie, nil
	}
}
