package config

import "os"

var SecretKey = []byte(os.Getenv("JWT_SECRET"))
