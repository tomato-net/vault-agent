package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	_getter = func(k string) interface{} {
		return k
	}

	t.Run("gets string", func(t *testing.T) {
		got := get[string]("my_key")
		assert.Equal(t, "my_key", got)
	})
}
