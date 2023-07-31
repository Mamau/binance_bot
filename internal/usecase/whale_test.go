package usecase

import (
	"binance_bot/internal/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupWhales(t *testing.T) {
	cases := []struct {
		Name     string
		Data     []entity.WhaleAction
		Expected []entity.WhaleAction
	}{
		{
			Name: "проверка успешногй группировки",
			Expected: []entity.WhaleAction{
				{
					ValueETH: 3,
					Hash:     "whale1",
				},
				{
					ValueETH: 1,
					Hash:     "whale2",
				},
			},
			Data: []entity.WhaleAction{
				{
					ValueETH: 1,
					Hash:     "whale1",
				},
				{
					ValueETH: 1,
					Hash:     "whale1",
				},
				{
					ValueETH: 1,
					Hash:     "whale1",
				},
				{
					ValueETH: 1,
					Hash:     "whale2",
				},
			},
		},
	}
	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			result := groupWhales(v.Data)
			assert.Equal(t, v.Expected, result)
		})
	}
}
