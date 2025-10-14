package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const prefix = "Bearer"

func (h *Handler) ModeratorMiddleware(allowedRole bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error"})
			return
		}
		if blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		jwtIsModerator, ok := claims["is_moderator"].(bool)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		if allowedRole && !jwtIsModerator {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "forbidden"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		return ""
	}

	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}

	return strings.Split(bearerToken, " ")[1]
}