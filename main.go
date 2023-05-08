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
	cipher_before, err2 := recrypt.Encrypt(apub, message)
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	err3 := recrypt.CheckCapsule(cipher_before.Capsule)
	if err3 != nil {
		fmt.Println("CheckCapsule() error:", err3.Error())
	} else {
		fmt.Println("CheckCapsule() success")
	}
	plain1, err4 := recrypt.Decrypt(apri, cipher_before)
	if err4 != nil {
		fmt.Println("Decrypt() error:", err4.Error())
	}
	bpub, bpri, err5 := curve.KeyGen()
	if err5 != nil {
		fmt.Println("bKeyGen() error:", err5.Error())
	}
	KF, err6 := recrypt.ReKeyGen(apub, apri, bpub, 10, 5)
	if err6 != nil {
		fmt.Println("ReKeyGen() error:", err6.Error())
	}
	cipher_after, err7 := recrypt.ReEncrypt(KF, cipher_before)
	if err7 != nil {
		fmt.Println("ReEncrypt() error:", err7.Error())
	}
	plain2, err8 := recrypt.DecryptFrags(apub, bpri, cipher_after, 5)
	if err8 != nil {
		fmt.Println("DecryptFrags() error:", err8.Error())
	}
	fmt.Println(string(plain1))
	//fmt.Println("KF:", KF)
	fmt.Println(string(plain2))
}
