package coas

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type SealedDraft struct {
	Version    string `json:"version"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func SealDraft(plaintext, key []byte) (SealedDraft, error) {
	if len(plaintext) == 0 { return SealedDraft{}, errors.New("plaintext is required") }
	if len(key) != 32 { return SealedDraft{}, errors.New("key must be exactly 32 bytes") }
	block, err := aes.NewCipher(key)
	if err != nil { return SealedDraft{}, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return SealedDraft{}, err }
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return SealedDraft{}, err }
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return SealedDraft{
		Version: "coas-private-draft.v1",
		Nonce: base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

func OpenDraft(sealed SealedDraft, key []byte) ([]byte, error) {
	if sealed.Version != "coas-private-draft.v1" { return nil, fmt.Errorf("unsupported sealed draft version: %s", sealed.Version) }
	if len(key) != 32 { return nil, errors.New("key must be exactly 32 bytes") }
	nonce, err := base64.StdEncoding.DecodeString(sealed.Nonce)
	if err != nil { return nil, fmt.Errorf("decode nonce: %w", err) }
	ciphertext, err := base64.StdEncoding.DecodeString(sealed.Ciphertext)
	if err != nil { return nil, fmt.Errorf("decode ciphertext: %w", err) }
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil { return nil, errors.New("authentication failed") }
	return plaintext, nil
}
