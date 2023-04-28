package main

import (
	"fmt"
	"goUmbral/curve"
	"goUmbral/recrypt"
)

func main() {
	message := "hello proxy-reencryption"
	pub, pri, err1 := curve.KeyGen()
	if err1 != nil {
		fmt.Println("KeyGen() error:", err1.Error())
	}
	cipher, err2 := recrypt.Encrypt(pub, message)
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	plain, err3 := recrypt.Decrypt(pri, cipher)
	if err3 != nil {
		fmt.Println("Decrypt() error:", err3.Error())
	}
	fmt.Println(string(plain))
}
