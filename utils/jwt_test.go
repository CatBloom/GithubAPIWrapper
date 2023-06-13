package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTFunctions(t *testing.T) {
	t.Logf("utils/jwt test")

	InitEnv()
	t.Run("Test GenerateJWT", func(t *testing.T) {
		key := "key"
		str := "str"
		jwtStr, err := GenerateJWT(key, str)
		if err != nil {
			t.Errorf("GenerateJWT returned an error: %v", err)
		}

		expectedResponse := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJzdHIifQ.Apo7im0Nynwy7mI7iByGZmnArjNy3LKlRfn8WicIjzY"
		assert.Equal(t, expectedResponse, jwtStr)
	})

	t.Run("Test RefreshJWT", func(t *testing.T) {
		key := "key"
		jwtStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJzdHIifQ.Apo7im0Nynwy7mI7iByGZmnArjNy3LKlRfn8WicIjzY"

		str, err := RefreshJWT(key, jwtStr)
		if err != nil {
			t.Errorf("RefreshJWT returned an error: %v", err)
		}

		expectedResponse := "str"
		assert.Equal(t, expectedResponse, str)
	})

	t.Run("Test GenerateRandomSecretKey", func(t *testing.T) {
		secretKey, err := GenerateRandomSecretKey()
		if err != nil {
			t.Errorf("GenerateRandomSecretKey returned an error: %v", err)
		}
		secretKey2, err := GenerateRandomSecretKey()
		if err != nil {
			t.Errorf("GenerateRandomSecretKey returned an error: %v", err)
		}

		assert.NotEqual(t, secretKey, secretKey2)
	})
}
