package utils

import "math/rand"

func GenerateCode(max int) string {
	chars := "0123456789"

	code := ""
	for i := 0; i <= max; i++ {
		code += string(chars[rand.Int()%len(chars)])
	}
	return code
}
