package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"goUmbral/curve"
	"goUmbral/math"
	"io"
	"log"
	"math/big"
	"os"

	"golang.org/x/crypto/sha3"
)

// concat bytes
func ConcatBytes(a, b []byte) []byte {
	var buf bytes.Buffer
	buf.Write(a)
	buf.Write(b)
	return buf.Bytes()
}

// convert message to hash value
func Sha3Hash(message []byte) ([]byte, error) {
	sha := sha3.New256()
	_, err := sha.Write(message)
	if err != nil {
		return nil, err
	}
	return sha.Sum(nil), nil
}

// map hash value to curve
func HashToCurve(hash []byte) *big.Int {
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt.Mod(hashInt, curve.N)
}

func GetCoefficients(a *big.Int, d *big.Int, t int) ([]*big.Int, error) {
	coefficients := []*big.Int{math.BigIntMul(a, math.GetInvert(d))}
	for i := 1; i < t; i++ {
		f, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		coefficients = append(coefficients, f.D)
	}
	return coefficients, nil
}

func GetPolynomialValue(coefficients []*big.Int, x *big.Int) *big.Int {
	t := len(coefficients)
	result := coefficients[t-1]
	for i := 1; i < t; i++ {
		result = math.BigIntAdd(math.BigIntMul(result, x), coefficients[t-i-1])
	}
	return result
}

// Gzip
func GzipCompress(bs []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(bs)
	zw.Flush()
	zw.Close()

	return buf.Bytes()
}

func GzipUnCompress(bs []byte) []byte {
	var buf bytes.Buffer
	bsBuf := bytes.NewBuffer(bs)
	zr, _ := gzip.NewReader(bsBuf)
	io.Copy(&buf, zr)
	zr.Close()

	return buf.Bytes()
}

func FileToMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
