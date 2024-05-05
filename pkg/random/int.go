package random

import (
	"math/rand"
	"time"
)

const (
	intArray = "0123456789"
)

func Int(length int) string {
	rand.NewSource(time.Now().UnixNano())
	id := make([]byte, length)
	for i := 0; i < length; i++ {
		id[i] = intArray[rand.Intn(len(intArray))]
	}
	return string(id)
}
