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

type Cipher_after_re struct {
	CF         []CFrag
	CipherText []byte
}

type KFrag struct {
	id  *ecdsa.PrivateKey
	rk  *big.Int
	X_A *ecdsa.PublicKey
	U_1 *ecdsa.PublicKey
	z_1 *big.Int
	z_2 *big.Int
}

type CFrag struct {
	E_1 *ecdsa.PublicKey
	V_1 *ecdsa.PublicKey
	id  *ecdsa.PrivateKey
	X_A *ecdsa.PublicKey
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

func ReKeyGen(aPriKey *ecdsa.PrivateKey, bPubKey *ecdsa.PublicKey, N int, t int) ([]KFrag, error) {
	if t < 2 {
		return nil, fmt.Errorf("%s", "t must bigger than 1")
	}
	X_A, x_A, err := curve.KeyGen()
	if err != nil {
		return nil, err
	}
	// get d = H3(X_A,pk_b,pk_b^(x_A))
	d := utils.HashToCurve(
		utils.ConcatBytes(
			utils.ConcatBytes(
				curve.PointToBytes(X_A),
				curve.PointToBytes(bPubKey)),
			curve.PointToBytes(curve.PointScalarMul(bPubKey, x_A.D))))
	coefficients, err := utils.GetCoefficients(aPriKey.D, d, t)
	//fmt.Println("coefficients:", coefficients)
	if err != nil {
		return nil, err
	}
	// get D = H6(pk_a,pk_b,pk_b^a)
	D := utils.HashToCurve(
		utils.ConcatBytes(
			utils.ConcatBytes(
				curve.PointToBytes(&aPriKey.PublicKey),
				curve.PointToBytes(bPubKey)),
			curve.PointToBytes(curve.PointScalarMul(bPubKey, aPriKey.D))))
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
				id.D.Bytes(),
				D.Bytes()))
		rk := utils.GetPolynomialValue(coefficients, s_x)
		U_1 := curve.BigIntMulBase(rk)
		// get z_1 = H4(Y,id,pk_a,pk_b,U_1,X_A)
		z_1 := utils.HashToCurve(
			utils.ConcatBytes(
				utils.ConcatBytes(
					utils.ConcatBytes(
						utils.ConcatBytes(
							utils.ConcatBytes(
								curve.PointToBytes(Y),
								id.D.Bytes()),
							curve.PointToBytes(&aPriKey.PublicKey)),
						curve.PointToBytes(bPubKey)),
					curve.PointToBytes(U_1)),
				curve.PointToBytes(X_A)))
		// get z_2 = y - a × z_1
		z_2 := math.BigIntSub(y.D, math.BigIntMul(aPriKey.D, z_1))
		kFrag := KFrag{
			id:  id,
			rk:  rk,
			X_A: X_A,
			U_1: U_1,
			z_1: z_1,
			z_2: z_2,
		}
		KF = append(KF, kFrag)
	}
	// KF长度为N
	return KF, nil
}

func ReEncapsulate(kFrag KFrag, capsule *Capsule) (*CFrag, error) {
	err := CheckCapsule(capsule)
	if err != nil {
		return nil, err
	}
	cFrag := CFrag{
		E_1: curve.PointScalarMul(capsule.E, kFrag.rk),
		V_1: curve.PointScalarMul(capsule.V, kFrag.rk),
		id:  kFrag.id,
		X_A: kFrag.X_A,
	}
	return &cFrag, nil
}

func ReEncrypt(KF []KFrag, cipher *Cipher_before_re) (*Cipher_after_re, error) {
	CF := []CFrag{}
	l := len(KF)
	for i := 0; i < l; i++ {
		cFrag, err := ReEncapsulate(KF[i], cipher.Capsule)
		if err != nil {
			return nil, err
		}
		CF = append(CF, *cFrag)
	}
	re_cipher := &Cipher_after_re{
		CF:         CF,
		CipherText: cipher.CipherText,
	}
	// re_cipher中CF长度为KF的长度，即默认为N
	return re_cipher, nil
}

func DecapsulateFrags(bPriKey *ecdsa.PrivateKey, aPubKey *ecdsa.PublicKey, CF []CFrag) ([]byte, error) {
	// get D = H6(pk_a,pk_b,pk_a^b)
	D := utils.HashToCurve(
		utils.ConcatBytes(
			utils.ConcatBytes(
				curve.PointToBytes(aPubKey),
				curve.PointToBytes(&bPriKey.PublicKey)),
			curve.PointToBytes(curve.PointScalarMul(aPubKey, bPriKey.D))))
	// 此处假设传入的CF切片长度为t
	t := len(CF)
	s_x := []*big.Int{}
	for i := 0; i < t; i++ {
		s_x_i := utils.HashToCurve(
			utils.ConcatBytes(
				CF[i].id.D.Bytes(),
				D.Bytes()))
		s_x = append(s_x, s_x_i)
	}
	lamda_S := []*big.Int{}
	for i := 1; i <= t; i++ {
		lamda_i_S := big.NewInt(1)
		for j := 1; j <= t; j++ {
			if j == i {
				continue
			} else {
				// bug点
				// a/b要通过a乘b的逆实现，不能通过除法
				lamda_i_S = math.BigIntMul(lamda_i_S, (math.BigIntMul(s_x[j-1], math.GetInvert(math.BigIntSub(s_x[j-1], s_x[i-1])))))
			}
		}
		lamda_S = append(lamda_S, lamda_i_S)
	}
	E := curve.PointScalarMul(CF[0].E_1, lamda_S[0])
	V := curve.PointScalarMul(CF[0].V_1, lamda_S[0])
	for i := 1; i < t; i++ {
		E = curve.PointScalarAdd(E, curve.PointScalarMul(CF[i].E_1, lamda_S[i]))
		V = curve.PointScalarAdd(V, curve.PointScalarMul(CF[i].V_1, lamda_S[i]))
	}
	// get d = H3(X_A,pk_b,X_A^b)
	d := utils.HashToCurve(
		utils.ConcatBytes(
			utils.ConcatBytes(
				curve.PointToBytes(CF[0].X_A),
				curve.PointToBytes(&bPriKey.PublicKey)),
			curve.PointToBytes(curve.PointScalarMul(CF[0].X_A, bPriKey.D))))
	point := curve.PointScalarMul(curve.PointScalarAdd(E, V), d)
	keyBytes, err := utils.Sha3Hash(curve.PointToBytes(point))
	if err != nil {
		return nil, err
	}
	return keyBytes, nil
}

func DecryptFrags(aPubKey *ecdsa.PublicKey, bPriKey *ecdsa.PrivateKey, re_cipher *Cipher_after_re, t int) (plainText []byte, err error) {
	// 此处假设传入的CF切片长度为默认的N
	CF := []CFrag{}
	// for i := 0; i < t; i++ {
	// 	CF = append(CF, re_cipher.CF[i])
	// }
	CF = append(CF, re_cipher.CF[1])
	CF = append(CF, re_cipher.CF[9])
	CF = append(CF, re_cipher.CF[7])
	CF = append(CF, re_cipher.CF[4])
	CF = append(CF, re_cipher.CF[3])
	keyBytes, err := DecapsulateFrags(bPriKey, aPubKey, CF)
	if err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)
	plainText, err = GCMDecrypt(re_cipher.CipherText, key[:32], keyBytes[:12], nil)
	if err != nil {
		return nil, err
	}
	return plainText, err
}
