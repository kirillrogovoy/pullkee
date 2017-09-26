package progress

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBar(t *testing.T) {
	t.Run(`0% progress`, func(t *testing.T) {
		var result string
		bar := Bar{10, func(s string) {
			result = s
		}}

		bar.Set(0)
		require.Equal(t, "[          ]", result)
	})

	t.Run(`50% progress`, func(t *testing.T) {
		var result string
		bar := Bar{10, func(s string) {
			result = s
		}}

		bar.Set(0.5)
		require.Equal(t, "[#####     ]", result)
	})

	t.Run(`100% progress`, func(t *testing.T) {
		var result string
		bar := Bar{10, func(s string) {
			result = s
		}}

		bar.Set(1)
		require.Equal(t, "[##########]", result)
	})

	t.Run("Panics on wrong input", func(t *testing.T) {
		bar := Bar{10, func(s string) {}}

		require.Panics(t, func() {
			bar.Set(2)
		})
	})
}
