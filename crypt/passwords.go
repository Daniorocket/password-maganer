package crypt

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func GeneratePassword(lenPass int64, wordsAlplabet string) string {
	var generatedPassword string
	for i := 0; i < int(lenPass); i++ {
		bigInt := big.NewInt(int64(len(wordsAlplabet)))
		index, err := rand.Int(rand.Reader, bigInt)
		Index := index.Int64()
		if err != nil {
			return ""
		}
		generatedPassword = generatedPassword + string(wordsAlplabet[Index])
	}
	return generatedPassword
}
func CalculateEntropy(password string, alphabet string) string {
	entrophy := math.Log2(float64(len(alphabet))) * float64(len(password))
	str := fmt.Sprintf("Entropia:%.2f bitów", entrophy)
	return str
}
func PrepareAlphabet(status bool, alphabetName string, alphabetValue string) string {
	if status == true {
		alphabetName = alphabetValue
	} else {
		alphabetName = ""
	}
	return alphabetName
}
