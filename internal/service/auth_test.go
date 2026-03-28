package service

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	token := "test-token-12345"

	hash1 := HashToken(token)
	hash2 := HashToken(token)

	if hash1 != hash2 {
		t.Error("same token should produce same hash")
	}

	if hash1 == token {
		t.Error("hash should differ from original token")
	}

	if len(hash1) != 64 {
		t.Errorf("SHA256 hex hash should be 64 chars, got %d", len(hash1))
	}
}

func TestHashTokenDifferentInputs(t *testing.T) {
	hash1 := HashToken("token-a")
	hash2 := HashToken("token-b")

	if hash1 == hash2 {
		t.Error("different tokens should produce different hashes")
	}
}
