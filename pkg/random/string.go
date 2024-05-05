package random

import (
	"math/rand"
	"time"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func String(length int) string {
	rand.NewSource(time.Now().UnixNano())
	id := make([]byte, length)
	for i := 0; i < length; i++ {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}
