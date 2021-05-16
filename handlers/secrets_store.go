package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

type secrets struct {
	// mu guards the secrets map.
	mu      sync.Mutex
	secrets map[string]string
}

func (s *secrets) AddSecret(v string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	hasher := md5.New()
	hasher.Write([]byte(v))
	m := hex.EncodeToString(hasher.Sum(nil))
	s.secrets[m] = v
	return m
}

func (s *secrets) GetDeleteSecret(secret string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	val := secretsValues.secrets[secret]
	delete(secretsValues.secrets, secret)
	return val
}
