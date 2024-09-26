package logruslog

import (
	"bytes"
	"io"
	"regexp"
)

type levelAwareWriter struct {
	InfoWriter  io.Writer
	ErrorWriter io.Writer
}

func (law *levelAwareWriter) Write(p []byte) (n int, err error) {
	// Убираем ANSI-коды для цвета перед проверкой уровня логирования
	cleanP := removeANSICodes(p)

	if bytes.Contains(cleanP, []byte("WARN[")) || bytes.Contains(cleanP, []byte("ERRO[")) || bytes.Contains(cleanP, []byte("FATA[")) {
		return law.ErrorWriter.Write(p)
	}
	return law.InfoWriter.Write(p) // По умолчанию все остальное идет в stdout
}

// precompile the ANSI escape sequence regular expression
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func removeANSICodes(data []byte) []byte {
	// Убираем ANSI reset
	cleanData := bytes.ReplaceAll(data, []byte("\x1b[0m"), []byte(""))
	// Убираем все ANSI escape-последовательности
	return ansiRegex.ReplaceAll(cleanData, []byte(""))
}
