package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/chafikTeck/hotel_reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	// Extract header properly
	return func(c *fiber.Ctx) error {
		values, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok || len(values) == 0 {
			return fmt.Errorf("unauthorized")
		}

		tokenStr := values[0]

		claims, err := validateToken(tokenStr)
		if err != nil {
			return err
		}

		// Check token expiration (custom "expires" claim)
		expiresFloat, ok := claims["expires"].(float64)
		if !ok {
			return fmt.Errorf("invalid token format")
		}

		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		userID := claims["id"].(primitive.ObjectID)
		user, err := userStore.GetUserByID(c.Context(), userID.String())
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		// Load secret
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("missing JWT secret")
		}

		return []byte(secret), nil
	})

	// Parse error
	if err != nil {
		fmt.Println("Failed to parse JWT:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
