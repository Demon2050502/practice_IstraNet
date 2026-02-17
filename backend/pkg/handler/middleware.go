package handler

import (
	"net/http"
	"os"
	"strings"

	"practice_IstraNet/pkg/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ctxUserIDKey = "userID"
	ctxRoleKey   = "role"
	ctxNameKey   = "name"
)

func (h *Handler) userIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			writeErr(c, http.StatusUnauthorized, "unauthorized", "нет токена")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeErr(c, http.StatusUnauthorized, "unauthorized", "неверный формат токена")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			writeErr(c, http.StatusInternalServerError, "internal_error", "JWT_SECRET не задан")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			writeErr(c, http.StatusUnauthorized, "unauthorized", "токен недействителен")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeErr(c, http.StatusUnauthorized, "unauthorized", "неверные claims")
			c.Abort()
			return
		}

		// sub может прийти как float64
		sub, ok := claims["sub"]
		if !ok {
			writeErr(c, http.StatusUnauthorized, "unauthorized", "нет sub")
			c.Abort()
			return
		}

		var userID int64
		switch v := sub.(type) {
		case float64:
			userID = int64(v)
		case int64:
			userID = v
		default:
			writeErr(c, http.StatusUnauthorized, "unauthorized", "неверный sub")
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)
		name, _ := claims["name"].(string)

		c.Set(ctxUserIDKey, userID)
		c.Set(ctxRoleKey, role)
		c.Set(ctxNameKey, name)

		c.Next()
	}
}

func getUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ctxUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

func getRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxRoleKey)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
var _ = dto.ErrorResponse{}
