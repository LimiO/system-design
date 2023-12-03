package jwt

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeMD5Hash(src string) string {
	hash := md5.Sum([]byte(src))
	hashed := hex.EncodeToString(hash[:])
	return hashed
}
