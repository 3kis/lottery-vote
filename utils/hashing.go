package utils

import (
	"crypto/sha1"
	"encoding/binary"
)

func GetUsernameHash(username string) uint32 {
	hash := sha1.Sum([]byte(username))
	return binary.BigEndian.Uint32(hash[:])
}
