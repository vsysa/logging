package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopyMapContext(t *testing.T) {
	tests := []struct {
		name     string
		original map[string]interface{}
		want     map[string]interface{}
	}{
		{
			name:     "Copy non-empty map",
			original: map[string]interface{}{"key1": "value1", "key2": 2},
			want:     map[string]interface{}{"key1": "value1", "key2": 2},
		},
		{
			name:     "Copy empty map",
			original: map[string]interface{}{},
			want:     map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Копируем исходную карту
			got := CopyMapContext(tt.original)

			// Проверяем, что изменения в оригинальной мапе не затронули копию
			tt.original["key1"] = "changedValue"
			assert.Equal(t, tt.want, got)

			// Проверяем, что копия осталась такой же, как до изменения оригинала
			if len(tt.original) > 0 {
				assert.NotEqual(t, tt.original, got, "CopyMapContext() создал неглубокую копию")
			}
		})
	}
}
