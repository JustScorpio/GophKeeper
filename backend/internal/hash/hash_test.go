// hash_test.go
package hash_test

import (
	"testing"

	"github.com/JustScorpio/GophKeeper/backend/internal/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("Хэширование пустого пароля", func(t *testing.T) {
		_, err := hash.HashPassword("")
		require.NoError(t, err)
	})

	t.Run("Очень длинный пароль", func(t *testing.T) {
		longPassword := "a" + string(make([]byte, 71)) // Длинный пароль (лимит bcrypt 72 бита)
		_, err := hash.HashPassword(longPassword)
		require.NoError(t, err)
	})

	t.Run("Специальные символы в пароле", func(t *testing.T) {
		specialPassword := "!@#$%^&*()_+-=[]{}|;:,.<>?~`"
		hashed, err := hash.HashPassword(specialPassword)
		require.NoError(t, err)

		assert.True(t, hash.CheckPasswordHash(specialPassword, hashed))
	})

	t.Run("Unicode символы в пароле", func(t *testing.T) {
		unicodePassword := "пароль123🔐🎉"
		hashed, err := hash.HashPassword(unicodePassword)
		require.NoError(t, err)

		assert.True(t, hash.CheckPasswordHash(unicodePassword, hashed))
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("Проверка с пустым хэшем", func(t *testing.T) {
		assert.False(t, hash.CheckPasswordHash("password", ""))
	})

	t.Run("Проверка с невалидным хэшем", func(t *testing.T) {
		assert.False(t, hash.CheckPasswordHash("password", "invalid-hash"))
	})

	t.Run("Проверка с правильным хэшем", func(t *testing.T) {
		password := "testpassword"
		hashed, err := hash.HashPassword(password)
		require.NoError(t, err)

		assert.True(t, hash.CheckPasswordHash(password, hashed))
		assert.False(t, hash.CheckPasswordHash("wrongpassword", hashed))
	})
}
