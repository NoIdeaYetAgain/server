package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)

}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomFullName() string {
	return RandomString(6)
}

func RandomPasswordHash() string {
	return RandomString(8)
}

func RandomRole() string {
	roles := []string{"student", "agent", "admin"}
	n := len(roles)
	return roles[rand.Intn(n)]
}

func RandomPhoneNumber() string {
	// Generates a pseudo-random 10-digit phone number, starting with a non-zero digit
	return fmt.Sprintf("0%d%d", rand.Intn(9)+1, rand.Intn(1e8))
}

func RandomGmail() string {
	// Generates something like "xabcde@gmail.com"
	username := RandomString(8)
	return fmt.Sprintf("%s@gmail.com", username)
}
