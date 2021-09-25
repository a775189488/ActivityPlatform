package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestMd5Encrypt(t *testing.T) {
	var aeskey = []byte("1234567887654321")
	pass := []byte("mccpwd")
	xpass := AesCbcEncrypt(pass, aeskey)

	pass64 := base64.StdEncoding.EncodeToString(xpass)
	fmt.Printf("加密后:%v\n", pass64)

	//bytesPass, err := base64.StdEncoding.DecodeString(pass64)
	bytesPass, err := base64.StdEncoding.DecodeString("PNlwSCDOLpqb6tW40QSbNA==")
	if err != nil {
		fmt.Println(err)
		return
	}

	tpass := AesCbcDecrypt(bytesPass, aeskey)
	fmt.Printf("解密后:%s, %d\n", tpass, len(tpass))
	fmt.Printf("MD5加密后:%s\n", strings.ToUpper(Md5Encrypt(tpass)))
}
