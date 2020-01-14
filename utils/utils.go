package utils

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
)

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateId() string {
	return fmt.Sprint(uuid.NewV4())
}
