package main

import (
	"fmt"
	"goUmbral/curve"
	"goUmbral/recrypt"
	"goUmbral/utils"
	"math/big"
	"os"
)

func main() {
	expected, _ := os.ReadFile("./test_txt/expected.txt")
	data1 := utils.GzipCompress(expected)
	os.WriteFile("./test_txt/compress.txt", data1, 0644)

	apub, apri, err1 := curve.KeyGen()
	if err1 != nil {
		fmt.Println("aKeyGen() error:", err1.Error())
	}
	cipher_before, err2 := recrypt.Encrypt(apub, "./test_txt/compress.txt", "./test_txt/encrypt.txt", big.NewInt(2))
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	err3 := recrypt.CheckCapsule(cipher_before.Capsule, big.NewInt(2))
	if err3 != nil {
		fmt.Println("CheckCapsule() error:", err3.Error())
	} else {
		fmt.Println("CheckCapsule() success")
	}
	err4 := recrypt.Decrypt(apri, cipher_before, "./test_txt/encrypt.txt", "./test_txt/a_decrypt.txt", big.NewInt(2))
	if err4 != nil {
		fmt.Println("Decrypt() error:", err4.Error())
	}

	data2, _ := os.ReadFile("./test_txt/a_decrypt.txt")
	actual1 := utils.GzipUnCompress(data2)
	os.WriteFile("./test_txt/a_actual.txt", actual1, 0644)

	bpub, bpri, err5 := curve.KeyGen()
	if err5 != nil {
		fmt.Println("bKeyGen() error:", err5.Error())
	}
	KF, err6 := recrypt.ReKeyGen(apri, bpub, 10, 5, big.NewInt(2))
	if err6 != nil {
		fmt.Println("ReKeyGen() error:", err6.Error())
	}
	cipher_after, err7 := recrypt.ReEncrypt(KF, cipher_before)
	if err7 != nil {
		fmt.Println("ReEncrypt() error:", err7.Error())
	}
	err8 := recrypt.DecryptFrags(apub, bpri, cipher_after, 5, "./test_txt/encrypt.txt", "./test_txt/b_decrypt.txt")
	if err8 != nil {
		fmt.Println("DecryptFrags() error:", err8.Error())
	}

	data3, _ := os.ReadFile("./test_txt/b_decrypt.txt")
	actual2 := utils.GzipUnCompress(data3)
	os.WriteFile("./test_txt/b_actual.txt", actual2, 0644)

	fmt.Println(utils.FileToMd5("./test_txt/expected.txt"))
	fmt.Println(utils.FileToMd5("./test_txt/a_actual.txt"))
	fmt.Println(utils.FileToMd5("./test_txt/b_actual.txt"))
}
