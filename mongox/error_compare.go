package mongox

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func NotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func IsDuplicateKey(err error) bool {
	if mongo.IsDuplicateKeyError(err) {
		return true
	}

	// Дополнительная проверка для разных версий драйвера
	if err != nil {
		errMsg := err.Error()
		return strings.Contains(errMsg, "duplicate key") ||
			strings.Contains(errMsg, "E11000")
	}

	return false
}

func IsTimeout(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "context deadline exceeded")
}

func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	return strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "network") ||
		strings.Contains(errMsg, "server selection timeout")
}

func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	// Проверяем ошибки валидации схемы
	errMsg := err.Error()
	return strings.Contains(errMsg, "Document failed validation") ||
		strings.Contains(errMsg, "ValidationError")
}

func ExtractDuplicateKeyInfo(err error) (field string, value interface{}) {
	if !IsDuplicateKey(err) {
		return "", nil
	}

	errMsg := err.Error()

	// Парсим сообщение об ошибке E11000
	// Пример: E11000 duplicate key error collection: db.campaigns index: name_1 dup key: { name: "Test Campaign" }
	if strings.Contains(errMsg, "dup key:") {
		// Простой парсинг - в продакшене лучше использовать regex
		parts := strings.Split(errMsg, "dup key:")
		if len(parts) > 1 {
			keyPart := strings.TrimSpace(parts[1])
			// Здесь можно добавить более сложный парсинг JSON-подобной структуры
			return "name", keyPart // Упрощенная версия
		}
	}

	return "unknown", nil
}
