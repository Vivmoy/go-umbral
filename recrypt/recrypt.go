package recrypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

type KFrag struct {
	id  *big.Int
	rk  *big.Int
	X_A *ecdsa.PublicKey
	U_1 *ecdsa.PublicKey
	z_1 *big.Int
	z_2 *big.Int
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
func Decapsulate(aPriKey *ecdsa.PrivateKey, capsule *Capsule) (keyBytes []byte, err error) {
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
	keyBytes, err := Decapsulate(aPriKey, cipher.Capsule)
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

func ReKeyGen(aPubKey *ecdsa.PublicKey, aPriKey *ecdsa.PrivateKey, bPubKey *ecdsa.PublicKey, N int, t int) ([]KFrag, error) {
	if t < 2 {
		return nil, fmt.Errorf("%s", "t must bigger than 1")
	}
	X_A, x_A, err := curve.KeyGen()
	if err != nil {
		return nil, err
	}
	// get d = H3(X_A,pk_b,pk_b^(X_A))
	d := utils.HashToCurve(
		utils.ConcatBytes(
			curve.PointToBytes(X_A),
			utils.ConcatBytes(
				curve.PointToBytes(bPubKey),
				curve.PointToBytes(curve.PointScalarMul(bPubKey, x_A.D)))))
	coefficients, err := utils.GetCoefficients(aPriKey.D, d, t)
	fmt.Println("coefficients:", coefficients)
	if err != nil {
		return nil, err
	}
	// get D = H6(pk_a,pk_b,pk_b^a)
	D := utils.HashToCurve(
		utils.ConcatBytes(
			curve.PointToBytes(aPubKey),
			utils.ConcatBytes(
				curve.PointToBytes(bPubKey),
				curve.PointToBytes(curve.PointScalarMul(bPubKey, aPriKey.D)))))
	KF := []KFrag{}
	for i := 0; i < N; i++ {
		Y, y, err := curve.KeyGen()
		if err != nil {
			return nil, err
		}
		id, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		// get s_x = H5(id,D)
		s_x := utils.HashToCurve(
			utils.ConcatBytes(
				id.X.Bytes(),
				D.Bytes()))
		rk := utils.GetPolynomialValue(coefficients, s_x)
		U_1 := curve.BigIntMulBase(rk)
		// get z_1 = H4(Y,id,pk_a,pk_b,U_1,X_A)
		z_1 := utils.HashToCurve(
			utils.ConcatBytes(
				curve.PointToBytes(Y),
				utils.ConcatBytes(
					id.D.Bytes(),
					utils.ConcatBytes(
						curve.PointToBytes(aPubKey),
						utils.ConcatBytes(
							curve.PointToBytes(bPubKey),
							utils.ConcatBytes(
								curve.PointToBytes(U_1),
								curve.PointToBytes(X_A)))))))
		// get z_2 = y - a × z_1
		z_2 := math.BigIntSub(y.D, math.BigIntMul(aPriKey.D, z_1))
		kFrag := KFrag{
			id:  id.D,
			rk:  rk,
			X_A: X_A,
			U_1: U_1,
			z_1: z_1,
			z_2: z_2,
		}
		KF = append(KF, kFrag)
	}
	return KF, nil
}
