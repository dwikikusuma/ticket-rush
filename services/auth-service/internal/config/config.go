package config

import "os"

var JWTSecrete []byte

func LoadConfig() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-123"
	}
	JWTSecrete = []byte(secret)
}
