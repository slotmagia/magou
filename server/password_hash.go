package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	// 生成密码哈希
	password := "123456"
	salt := "salt123"

	// 按照系统的加密方式：password + salt 的MD5
	combined := password + salt
	hash := md5.Sum([]byte(combined))
	hashString := fmt.Sprintf("%x", hash)

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("盐值: %s\n", salt)
	fmt.Printf("组合: %s\n", combined)
	fmt.Printf("MD5哈希: %s\n", hashString)

	// 验证
	fmt.Println("\n验证:")
	testHash := md5.Sum([]byte("123456admin"))
	fmt.Printf("验证哈希: %x\n", testHash)
}
