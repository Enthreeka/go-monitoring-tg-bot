package callback

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindTitle(t *testing.T) {
	tests := []struct {
		name     string
		caption  string
		wantName string
	}{
		{
			name: "ok",
			caption: "Управление каналом" +
				"" +
				"Канал:chinazez " +
				"Количество людей, которые ожидают принятия: 2",
			wantName: "chinazez",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := findTitle(tt.caption)
			assert.Equal(t, tt.wantName, name)
		})
	}
}
