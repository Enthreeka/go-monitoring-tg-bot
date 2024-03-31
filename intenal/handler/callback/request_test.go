package callback

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFindTitle(t *testing.T) {
	tests := []struct {
		name     string
		caption  string
		wantName string
	}{
		{
			name:     "ok",
			caption:  "Управление каналом\nКанал:Beta test \n\nКоличество людей, которые ожидают принятия: 0",
			wantName: "Beta test",
		},
		{
			name:     "ok",
			caption:  "Управление каналом\nКанал:Beta test new channel \n\nКоличество людей, которые ожидают принятия: 0",
			wantName: "Beta test new channel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := findTitle(tt.caption)
			assert.Equal(t, tt.wantName, name)
		})
	}
}

func TestStrings(t *testing.T) {
	_, after, _ := strings.Cut("Вы действительно хотите присоединиться к каналу: china?", "к каналу:")
	t.Log(after[1 : len(after)-1])
}
