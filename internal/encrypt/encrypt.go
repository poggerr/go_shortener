package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/poggerr/go_shortener/internal/logger"
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Encrypt(pass string) string {
	src := []byte(pass)

	key, err := generateRandom(2 * aes.BlockSize)
	if err != nil {
		logger.Initialize().Error(err)
	}
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		logger.Initialize().Error(err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		logger.Initialize().Error(err)
	}

	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		logger.Initialize().Error(err)
	}

	dst := aesgcm.Seal(nil, nonce, src, nil)

	encrypt := hex.EncodeToString(dst)
	return encrypt
}
