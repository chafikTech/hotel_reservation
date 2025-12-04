package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--JWT auth")

	token := c.Get("X-API-Token")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing token")
	}

	if err := parseToken(token); err != nil {
		return err
	}

	fmt.Println("Valid token:", token)

	// continue to next handler
	return c.Next()
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("missing jwt secret")
		}

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT:", err)
		return fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("unauthorized")
	}

	// Debug prints
	fmt.Println("foo:", claims["foo"])
	fmt.Println("nbf:", claims["nbf"])

	return nil
}
