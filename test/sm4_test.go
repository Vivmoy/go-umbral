package test

import (
	"goUmbral/curve"
	"goUmbral/recrypt"
	"goUmbral/utils"
	"math/big"
	"os"
	"testing"

	"github.com/tjfoc/gmsm/sm4"
)

func TestSm4(t *testing.T) {
	pub, _, _ := curve.KeyGen()
	keyBytes, _, _ := recrypt.Encapsulate(pub, big.NewInt(2))
	// t.Log("len(keyBytes):", len(keyBytes[:16]))
	// key := []byte("1234567890abcdef")
	// t.Log("len(key):", len(key))
	data := "你好"
	iv := []byte("0000000000000000")
	sm4.SetIV(iv)
	ecbMsg, _ := sm4.Sm4Ecb(keyBytes[:16], []byte(data), true)
	ecbDec, _ := sm4.Sm4Ecb(keyBytes[:16], ecbMsg, false)
	t.Log("ecbDec:", string(ecbDec))
}

func TestSm4File(t *testing.T) {
	pub, _, _ := curve.KeyGen()
	keyBytes, _, _ := recrypt.Encapsulate(pub, big.NewInt(2))
	expected, _ := os.ReadFile("CA_fameleBirth.txt")
	ofbMsg, _ := sm4.Sm4OFB(keyBytes[:16], expected, true)
	ofbDec, _ := sm4.Sm4OFB(keyBytes[:16], ofbMsg, false)
	os.WriteFile("Sm4FileTest.txt", ofbDec, 0644)
	t.Log(utils.FileToMd5("CA_fameleBirth.txt"))
	t.Log(utils.FileToMd5("Sm4FileTest.txt"))
}
