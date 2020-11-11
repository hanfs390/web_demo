package v2

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func getMD5(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
func generateUuid() (uu string, err error) {
	u := uuid.NewV4()
	return u.String(), err
}
