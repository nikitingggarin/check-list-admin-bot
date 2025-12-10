package utils

import (
	"strconv"
	"strings"
)

// IsPencilNumberButton проверяет, является ли текст кнопкой с карандашом и номером
func IsPencilNumberButton(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}

	// Проверяем, начинается ли с карандаша
	if !strings.HasPrefix(text, "✏️") {
		return false
	}

	// Убираем карандаш и пробелы
	cleanText := strings.TrimPrefix(text, "✏️")
	cleanText = strings.TrimSpace(cleanText)

	if cleanText == "" {
		return false
	}

	// Используем strconv.Atoi для проверки (более надежно)
	_, err := strconv.Atoi(cleanText)
	return err == nil
}

// ExtractNumberFromPencilButton извлекает номер из кнопки с карандашом
func ExtractNumberFromPencilButton(text string) (int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, strconv.ErrSyntax
	}

	// Убираем карандаш и пробелы
	cleanText := strings.TrimPrefix(text, "✏️")
	cleanText = strings.TrimSpace(cleanText)

	if cleanText == "" {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(cleanText)
}

// ExtractBlockIndexFromButton - извлекает индекс блока из текста кнопки
func ExtractBlockIndexFromButton(buttonText string) (int, bool) {
	// Формат: "🧱 1. Название (3)" или "📭 1. Название (0)"
	parts := strings.SplitN(buttonText, ".", 2)
	if len(parts) < 2 {
		return -1, false
	}

	// Извлекаем номер из первой части (убираем эмодзи и пробелы)
	numPart := strings.TrimSpace(parts[0])
	// Убираем эмодзи
	for i, r := range numPart {
		if r >= '0' && r <= '9' {
			numPart = numPart[i:]
			break
		}
	}

	index, err := strconv.Atoi(numPart)
	if err != nil {
		return -1, false
	}

	return index - 1, true // Конвертируем в 0-based индекс
}
