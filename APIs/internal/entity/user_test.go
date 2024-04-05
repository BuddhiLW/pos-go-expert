package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("name", "email", "password")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("password"))
	assert.False(t, user.ValidatePassword("password!"))
	assert.NotEqual(t, "password", user.Password)
}
