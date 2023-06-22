package main

import (
	"bufio"
	"fmt"
	"goUmbral/curve"
	"goUmbral/recrypt"
	"goUmbral/utils"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func use_gzip() {
	expected, _ := os.ReadFile("./test_gzip_txt/expected.txt")
	data1 := utils.GzipCompress(expected)
	os.WriteFile("./test_gzip_txt/compress.txt", data1, 0644)

	apub, apri, err1 := curve.KeyGen()
	if err1 != nil {
		fmt.Println("aKeyGen() error:", err1.Error())
	}
	cipher_before, err2 := recrypt.Encrypt(apub, "./test_gzip_txt/compress.txt", "./test_gzip_txt/encrypt.txt", big.NewInt(2))
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	err3 := recrypt.CheckCapsule(cipher_before.Capsule, big.NewInt(2))
	if err3 != nil {
		fmt.Println("CheckCapsule() error:", err3.Error())
	} else {
		fmt.Println("CheckCapsule() success")
	}
	err4 := recrypt.Decrypt(apri, cipher_before, "./test_gzip_txt/encrypt.txt", "./test_gzip_txt/a_decrypt.txt", big.NewInt(2))
	if err4 != nil {
		fmt.Println("Decrypt() error:", err4.Error())
	}

	data2, _ := os.ReadFile("./test_gzip_txt/a_decrypt.txt")
	actual1 := utils.GzipUnCompress(data2)
	os.WriteFile("./test_gzip_txt/a_actual.txt", actual1, 0644)

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
	err8 := recrypt.DecryptFrags(apub, bpri, cipher_after, 5, "./test_gzip_txt/encrypt.txt", "./test_gzip_txt/b_decrypt.txt")
	if err8 != nil {
		fmt.Println("DecryptFrags() error:", err8.Error())
	}

	data3, _ := os.ReadFile("./test_gzip_txt/b_decrypt.txt")
	actual2 := utils.GzipUnCompress(data3)
	os.WriteFile("./test_gzip_txt/b_actual.txt", actual2, 0644)

	fmt.Println(utils.FileToMd5("./test_gzip_txt/expected.txt"))
	fmt.Println(utils.FileToMd5("./test_gzip_txt/a_actual.txt"))
	fmt.Println(utils.FileToMd5("./test_gzip_txt/b_actual.txt"))
}

func use_gorilla() {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("./test_gorilla_txt/expected.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("./test_gorilla_txt/compress.txt", data1, 0644)

	apub, apri, err1 := curve.KeyGen()
	if err1 != nil {
		fmt.Println("aKeyGen() error:", err1.Error())
	}
	cipher_before, err2 := recrypt.Encrypt(apub, "./test_gorilla_txt/compress.txt", "./test_gorilla_txt/encrypt.txt", big.NewInt(2))
	if err2 != nil {
		fmt.Println("Encrypt() error:", err2.Error())
	}
	err3 := recrypt.CheckCapsule(cipher_before.Capsule, big.NewInt(2))
	if err3 != nil {
		fmt.Println("CheckCapsule() error:", err3.Error())
	} else {
		fmt.Println("CheckCapsule() success")
	}
	err4 := recrypt.Decrypt(apri, cipher_before, "./test_gorilla_txt/encrypt.txt", "./test_gorilla_txt/a_decrypt.txt", big.NewInt(2))
	if err4 != nil {
		fmt.Println("Decrypt() error:", err4.Error())
	}

	data2, _ := os.ReadFile("./test_gorilla_txt/a_decrypt.txt")
	fWrite1, _ := os.Create("./test_gorilla_txt/a_actual.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	it1 := s2.Iter()
	for it1.Next() {
		t, v := it1.Values()
		strResult := fmt.Sprintf("%v %v\n", v, t)
		fWrite1.WriteString(strResult)
	}

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
	err8 := recrypt.DecryptFrags(apub, bpri, cipher_after, 5, "./test_gorilla_txt/encrypt.txt", "./test_gorilla_txt/b_decrypt.txt")
	if err8 != nil {
		fmt.Println("DecryptFrags() error:", err8.Error())
	}

	data3, _ := os.ReadFile("./test_gorilla_txt/b_decrypt.txt")
	fWrite2, _ := os.Create("./test_gorilla_txt/b_actual.txt")
	s3 := utils.New(uint32(1))
	s3.UnmarshalBinary(data3)
	it2 := s3.Iter()
	for it2.Next() {
		t, v := it2.Values()
		strResult := fmt.Sprintf("%v %v\n", v, t)
		fWrite2.WriteString(strResult)
	}

	fmt.Println(utils.FileToMd5("./test_gorilla_txt/expected.txt"))
	fmt.Println(utils.FileToMd5("./test_gorilla_txt/a_actual.txt"))
	fmt.Println(utils.FileToMd5("./test_gorilla_txt/b_actual.txt"))
}

func main() {
	use_gzip()
	fmt.Println("-----------------")
	use_gorilla()
}
