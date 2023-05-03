package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnShorten(t *testing.T) {
	tests := []struct {
		name   string
		oldUrl string
	}{
		{
			name:   "Normal test",
			oldUrl: "https://practicum.yandex.ru/",
		},
		{
			name:   "Empty test",
			oldUrl: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newUrl := Shorting(tt.oldUrl)
			ans := UnShorting(newUrl)
			assert.Equal(t, tt.oldUrl, ans)
		})
	}
}
