package tools

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"

	"golang.org/x/crypto/bcrypt"
)

func MarshalJSON(v any) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

func UnmarshalJSON(data []byte, v any) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}

// CheckLuhn выполняет проверку последовательности
// цифр при помощи алгоритма Луны.
//
// https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0
func CheckLuhn(sequence string) bool {
	if sequence == "" {
		return false
	}

	sum := 0
	size := len(sequence)
	parity := size % 2

	for i := range size {
		digit, err := strconv.Atoi(string(sequence[i]))
		if err != nil {
			return false
		}

		if i%2 == parity {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	return sum%10 == 0
}

// HashPassword хеширует пароль.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CheckPassword сверяет пароль с хешированным паролем.
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
