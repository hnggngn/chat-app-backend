package utils

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"
)

func GetPrivateKey() ed25519.PrivateKey {
	seed, _ := hex.DecodeString(os.Getenv("PRIVATE_KEY"))
	return ed25519.NewKeyFromSeed(seed)
}
