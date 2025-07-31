package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Создаем SHA-256 хеш от входной строки

func GenerateHash(text string) string {

	hasher := sha256.New()

	// Записываем данные в хешер
	hasher.Write([]byte(text))

	// Получаем финальный хеш как массив байтов
	hashBytes := hasher.Sum(nil)

	// Конвертируем байты в строку в шестнадцатеричном формате
	return hex.EncodeToString(hashBytes)
}
