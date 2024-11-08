package auth

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for i := 0; i < 6; i++ {
		randomIndex := rng.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
