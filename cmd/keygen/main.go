package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateSecret(bytes int) string {
	b := make([]byte, bytes)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func main() {
	fmt.Println("SESSION_SECRET=" + generateSecret(32))
	fmt.Println("CSRF_SECRET=" + generateSecret(32))
	fmt.Println("APP_KEY=" + generateSecret(32))
}
