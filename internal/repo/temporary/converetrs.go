package temporary

import (
	"crypto/md5"
	"encoding/hex"
)

func keygen(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
