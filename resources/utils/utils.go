package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

func GenerateID(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Concat strings in slice with suffix
func Concat(suffix string, strSlice ...string) string {
	var outBuffer bytes.Buffer
	for _, str := range strSlice {
		outBuffer.WriteString(str)
		outBuffer.WriteString(suffix)
	}
	return strings.TrimSuffix(outBuffer.String(), ",")
}
