package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Sum(str string) (md5str string) {
	m := md5.New()
	m.Write([]byte(str))
	md5str = hex.EncodeToString(m.Sum(nil))
	return
}
