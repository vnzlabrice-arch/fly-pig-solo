package pkg

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 加密字符串
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
