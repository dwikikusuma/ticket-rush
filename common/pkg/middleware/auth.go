package middleware

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var RequestIDKet = "request-id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		parts := strings.Split(header, " ")
		if len(parts) != 1 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev_123"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userId, ok := claims["sub"].(float64)
			if !ok {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			email, ok := claims["email"].(string)
			if !ok {
				c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
				return
			}

			c.Set("userId", int(userId))
			c.Set("email", email)

			c.Next()
		}
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.NewString()

		c.Set(RequestIDKet, reqID)
		c.Header("X-Request-Id", reqID)

		ctx := context.WithValue(c.Request.Context(), RequestIDKet, reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func TimeOut(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})

		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(504, gin.H{"error": "Request timed out"})
		case <-finished:
		}
	}
}
