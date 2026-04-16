// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

type PasswordRecord struct {
	Algorithm  string `json:"algorithm"`
	Salt       string `json:"salt"`
	Hash       string `json:"hash"`
	Iterations uint32 `json:"iterations"`
	MemoryKB   uint32 `json:"memory_kb"`
	Parallel   uint8  `json:"parallel"`
	KeyLen     uint32 `json:"key_len"`
}

func passwordFile(secretsDir string) string {
	return filepath.Join(secretsDir, "escalation_password.json")
}

func SetPassword(secretsDir string, password string) error {
	if len(strings.TrimSpace(password)) < 10 {
		return errors.New("password must be at least 10 chars")
	}
	algo := strings.ToLower(os.Getenv("GLIDECLAW_PASSWORD_ALGO"))
	if algo == "" {
		algo = "argon2id"
	}
	rec, err := hashPassword(password, algo)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(secretsDir, 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(passwordFile(secretsDir), data, 0o600)
}

func VerifyPassword(secretsDir string, password string) (bool, error) {
	data, err := os.ReadFile(passwordFile(secretsDir))
	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.New("escalation password is not configured")
		}
		return false, err
	}
	var rec PasswordRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(rec.Salt)
	if err != nil {
		return false, err
	}
	expected, err := base64.RawStdEncoding.DecodeString(rec.Hash)
	if err != nil {
		return false, err
	}
	actual, err := derive(rec, password, salt)
	if err != nil {
		return false, err
	}
	ok := subtle.ConstantTimeCompare(actual, expected) == 1
	return ok, nil
}

func hashPassword(password string, algo string) (PasswordRecord, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return PasswordRecord{}, err
	}
	rec := PasswordRecord{KeyLen: 32}
	switch algo {
	case "argon2id":
		rec.Algorithm = "argon2id"
		rec.Iterations = 3
		rec.MemoryKB = 64 * 1024
		rec.Parallel = 1
	case "pbkdf2":
		rec.Algorithm = "pbkdf2"
		rec.Iterations = 600000
	default:
		return PasswordRecord{}, fmt.Errorf("unsupported password algorithm %q", algo)
	}
	derived, err := derive(rec, password, salt)
	if err != nil {
		return PasswordRecord{}, err
	}
	rec.Salt = base64.RawStdEncoding.EncodeToString(salt)
	rec.Hash = base64.RawStdEncoding.EncodeToString(derived)
	return rec, nil
}

func derive(rec PasswordRecord, password string, salt []byte) ([]byte, error) {
	switch rec.Algorithm {
	case "argon2id":
		return argon2.IDKey([]byte(password), salt, rec.Iterations, rec.MemoryKB, rec.Parallel, rec.KeyLen), nil
	case "pbkdf2":
		return pbkdf2.Key([]byte(password), salt, int(rec.Iterations), int(rec.KeyLen), sha256.New), nil
	default:
		return nil, fmt.Errorf("unsupported password algorithm %q", rec.Algorithm)
	}
}

func PasswordConfigured(secretsDir string) bool {
	_, err := os.Stat(passwordFile(secretsDir))
	return err == nil
}
