package helpers

import "math/rand"

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letters          = uppercaseLetters + lowercaseLetters
	numbers          = "1234567890"
	characters       = uppercaseLetters + lowercaseLetters + numbers
)

func RandomCharacters(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = characters[rand.Intn(len(characters))]
	}

	return string(b)
}

func RandomIntegers(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = numbers[rand.Intn(len(numbers))]
	}

	return string(b)
}

func RandomLowerLetters(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = lowercaseLetters[rand.Intn(len(lowercaseLetters))]
	}

	return string(b)
}

func RandomUpperLetters(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = uppercaseLetters[rand.Intn(len(uppercaseLetters))]
	}

	return string(b)
}

func Letters(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
