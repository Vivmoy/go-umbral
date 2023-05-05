package main

import (
	"fmt"
	"goUmbral/curve"
	"goUmbral/recrypt"
)

func main() {
	message := "hello proxy-reencryption"
	apub, apri, err1 := curve.KeyGen()
	if err1 != nil {
		fmt.Println("aKeyGen() error:", err1.Error())
	}
	cipher, err2 := recrypt.Encrypt(apub, message)
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	err3 := recrypt.CheckCapsule(cipher.Capsule)
	if err3 != nil {
		fmt.Println("CheckCapsule() error:", err3.Error())
	} else {
		fmt.Println("CheckCapsule() success")
	}
	plain, err4 := recrypt.Decrypt(apri, cipher)
	if err4 != nil {
		fmt.Println("Decrypt() error:", err4.Error())
	}
	bpub, _, err5 := curve.KeyGen()
	if err5 != nil {
		fmt.Println("bKeyGen() error:", err5.Error())
	}
	KF, err6 := recrypt.ReKeyGen(apub, apri, bpub, 10, 5)
	if err6 != nil {
		fmt.Println("ReKeyGen() error:", err6.Error())
	}
	fmt.Println(string(plain))
	fmt.Println("KF:", KF)
}
