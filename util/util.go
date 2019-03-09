package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"gitee.com/NotOnlyBooks/statistical_report/constant"
	"sort"
	"strings"
)

func Md5String(plain string) string {
	cipher := MD5([]byte(plain))
	return hex.EncodeToString(cipher)
}

func MD5(plain []byte) []byte {
	md5Cxt := md5.New()
	md5Cxt.Write(plain)
	cipher := md5Cxt.Sum(nil)
	return cipher[:]
}

func VerifySIGN(m map[string]string, signStr string) bool {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var writer bytes.Buffer
	for _, key := range keys {
		if m[key] == "" || m[key] == "0" {
			continue
		}

		writer.WriteString(key)
		writer.WriteString("=")
		writer.WriteString(m[key])
		writer.WriteString("&")
	}
	writer.WriteString(constant.SignKey)
	return strings.ToLower(Md5String(writer.String())) == strings.ToLower(signStr)
}
