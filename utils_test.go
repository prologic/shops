package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseeURI(t *testing.T) {
	assert := assert.New(t)

	t.Run("local://", func(t *testing.T) {
		u := ParseURI("local://", "", "")
		assert.Equal("local", u.Type)
		assert.Equal("", u.User)
		assert.Equal("", u.Host)
		assert.Equal("", u.Port)
	})

	t.Run("ssh://foo@bar", func(t *testing.T) {
		u := ParseURI("ssh://foo@bar", "root", "22")
		assert.Equal("ssh", u.Type)
		assert.Equal("foo", u.User)
		assert.Equal("bar", u.Host)
		assert.Equal("22", u.Port)
	})

	t.Run("192.168.0.1", func(t *testing.T) {
		u := ParseURI("192.168.0.1", "root", "22")
		assert.Equal("ssh", u.Type)
		assert.Equal("root", u.User)
		assert.Equal("192.168.0.1", u.Host)
		assert.Equal("22", u.Port)
	})

	t.Run("ssh://192.168.0.1:2222", func(t *testing.T) {
		u := ParseURI("ssh://192.168.0.1:2222", "root", "22")
		assert.Equal("ssh", u.Type)
		assert.Equal("root", u.User)
		assert.Equal("192.168.0.1", u.Host)
		assert.Equal("2222", u.Port)
	})
}
