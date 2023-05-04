package recrypt

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"goUmbral/curve"
	"goUmbral/math"
	"goUmbral/utils"
	"math/big"
)

type Capsule struct {
	E *ecdsa.PublicKey
	V *ecdsa.PublicKey
	S *big.Int
}

type Cipher_before_re struct {
	CipherText []byte
	Capsule    *Capsule
}

func Encapsulate(pubKey *ecdsa.PublicKey) (keyBytes []byte, capsule *Capsule, err error) {
	s := new(big.Int)
	// generate E,V key-pairs
	pubE, priE, err := curve.KeyGen()
	pubV, priV, err := curve.KeyGen()
	if err != nil {
		return nil, nil, err
	}
	// get H2(E || V)
	h := utils.HashToCurve(
		utils.ConcatBytes(
			curve.PointToBytes(pubE),
			curve.PointToBytes(pubV)))
	// get s = v + e * H2(E || V)
	s = math.BigIntAdd(priV.D, math.BigIntMul(priE.D, h))
	// get (pk_A)^{e+v}
	point := curve.PointScalarMul(pubKey, math.BigIntAdd(priE.D, priV.D))
	// generate aes key
	keyBytes, err = utils.Sha3Hash(curve.PointToBytes(point))
	if err != nil {
		return nil, nil, err
	}
	capsule = &Capsule{
		E: pubE,
		V: pubV,
		S: s,
	}
	return keyBytes, capsule, nil
}

func Encrypt(pubKey *ecdsa.PublicKey, message string) (cipher *Cipher_before_re, err error) {
	keyBytes, capsule, err := Encapsulate(pubKey)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	cipherText, err := GCMEncrypt([]byte(message), key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, err
	}
	cipher = &Cipher_before_re{
		CipherText: cipherText,
		Capsule:    capsule,
	}
	return cipher, nil
}

// Recreate aes key
func Decapsulate(capsule *Capsule, aPriKey *ecdsa.PrivateKey) (keyBytes []byte, err error) {
	point1 := curve.PointScalarAdd(capsule.E, capsule.V)
	point := curve.PointScalarMul(point1, aPriKey.D)
	// generate aes key
	keyBytes, err = utils.Sha3Hash(curve.PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

func Decrypt(aPriKey *ecdsa.PrivateKey, cipher *Cipher_before_re) (plainText []byte, err error) {
	keyBytes, err := Decapsulate(cipher.Capsule, aPriKey)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)
	// use aes gcm algorithm to encrypt
	// mark keyBytes[:12] as nonce
	plainText, err = GCMDecrypt(cipher.CipherText, key[:32], keyBytes[:12], nil)
	return plainText, err
}

func CheckCapsule(capsule *Capsule) (err error) {
	left := curve.BigIntMulBase(capsule.S)
	h1 := utils.HashToCurve(
		utils.ConcatBytes(
			curve.PointToBytes(capsule.E),
			curve.PointToBytes(capsule.V)))
	h2 := curve.PointScalarMul(capsule.E, h1)
	right := curve.PointScalarAdd(capsule.V, h2)
	if left.Equal(right) {
		return nil
	}
	return fmt.Errorf("%s", "Capsule not match")
}
