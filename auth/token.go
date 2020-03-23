package auth

import "math/rand"

var (
	tokenCharacters   = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	randomTokenLength = 14
	applicationPrefix = "A"
	clientPrefix      = "C"
	pluginPrefix      = "P"
)

//GenerateApplicationToken generates an application token.
func GenerateApplicationToken() string {
	return generateRandomToken(applicationPrefix)
}

//GenerateClientToken generates a client token.
func GenerateClientToken() string {
	return generateRandomToken(clientPrefix)
}

//GenerateImageName generates an image name
func GenerateImageName() string {
	return generateRandomString(25)
}

func generateRandomToken(prefix string) string {
	return prefix + generateRandomString(randomTokenLength)
}

func generateRandomString(length int) string {
	res := make([]byte, length)
	for i := range res {
		res[i] = tokenCharacters[rand.Intn(len(tokenCharacters))]
	}
	return string(res)
}
