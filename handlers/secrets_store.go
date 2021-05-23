package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/borkod/secrets-app/fileHandler"
)

type secretsStore struct {
	// mu guards the secrets map.
	mu sync.Mutex
	// stores secret values
	secrets      map[string]string
	dataFileName string
}

// Adds a secret to the map data structure
func (s *secretsStore) AddSecret(v string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := fileHandler.ReadFromFile(s.dataFileName)
	if err != nil {
		return "", err
	}
	if len(jsonData) != 0 {
		json.Unmarshal(jsonData, &s.secrets)
	}
	hasher := md5.New()
	hasher.Write([]byte(v))
	m := hex.EncodeToString(hasher.Sum(nil))
	s.secrets[m] = v
	jsonData, err = json.Marshal(s.secrets)
	if err != nil {
		return "", err
	}
	return m, fileHandler.WriteToFile(jsonData, s.dataFileName)
}

// Reads the requested secret value; deletes the secret in the map data structure; returns the secret value
func (s *secretsStore) GetDeleteSecret(secret string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := fileHandler.ReadFromFile(s.dataFileName)
	if err != nil {
		return "", err
	}
	if len(jsonData) != 0 {
		json.Unmarshal(jsonData, &s.secrets)
	}

	val := ""
	_, ok := secretsValues.secrets[secret]
	if ok {
		val = secretsValues.secrets[secret]
	}
	delete(secretsValues.secrets, secret)
	jsonData, err = json.Marshal(s.secrets)
	if err != nil {
		return "", err
	}
	return val, fileHandler.WriteToFile(jsonData, s.dataFileName)
}
