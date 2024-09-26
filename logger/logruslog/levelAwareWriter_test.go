package logruslog

import (
	"bytes"
	"testing"
)

// TestLevelAwareWriter проверяет, что логи правильно распределяются между InfoWriter и ErrorWriter
func TestLevelAwareWriter(t *testing.T) {
	// Создаем буферы для имитации потоков вывода
	infoBuffer := new(bytes.Buffer)
	errorBuffer := new(bytes.Buffer)

	// Инициализируем наш levelAwareWriter с этими буферами
	law := &levelAwareWriter{
		InfoWriter:  infoBuffer,
		ErrorWriter: errorBuffer,
	}

	// Тестовые данные
	tests := []struct {
		input       string
		expectInfo  bool // true, если ожидаем вывод в infoBuffer
		expectError bool // true, если ожидаем вывод в errorBuffer
	}{
		{"INFO[2024-06-05T11:28:00.408] Configuration updated", true, false},
		{"WARN[2024-06-05T11:28:00.407] No handler registered", false, true},
		{"ERRO[2024-06-05T11:28:00.409] Error accessing database", false, true},
		{"FATA[2024-06-05T11:28:00.410] Critical system failure", false, true},
		{"\x1b[32mINFO[2024-06-05T11:28:00.408] Configuration updated\x1b[0m", true, false},
		{"\x1b[33mWARN[2024-06-05T11:28:00.407] No handler registered\x1b[0m", false, true},
		{"\x1b[31mERRO[2024-06-05T11:28:00.409] Error accessing database\x1b[0m", false, true},
		{"\x1b[35mFATA[2024-06-05T11:28:00.410] Critical system failure\x1b[0m", false, true},
	}

	// Проходим по каждому тестовому случаю
	for _, tt := range tests {
		// Сброс буферов
		infoBuffer.Reset()
		errorBuffer.Reset()

		// Запись данных в levelAwareWriter
		_, err := law.Write([]byte(tt.input))
		if err != nil {
			t.Errorf("Failed to write data: %s", err)
		}

		// Проверяем, что содержимое буферов соответствует ожиданиям
		if (infoBuffer.Len() > 0) != tt.expectInfo {
			t.Errorf("InfoWriter output for input %q expected %v, got %v", tt.input, tt.expectInfo, infoBuffer.Len() > 0)
		}
		if (errorBuffer.Len() > 0) != tt.expectError {
			t.Errorf("ErrorWriter output for input %q expected %v, got %v", tt.input, tt.expectError, errorBuffer.Len() > 0)
		}
	}
}
