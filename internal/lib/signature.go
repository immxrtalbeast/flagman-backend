package lib

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

func NewSignature(user *domain.User, salt string) string {
	hash := sha256.Sum256([]byte(user.PhoneNumber + salt))
	return hex.EncodeToString(hash[:])
}
