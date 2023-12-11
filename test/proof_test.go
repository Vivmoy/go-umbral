package test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"goUmbral/curve"
	"goUmbral/math"
	"goUmbral/recrypt"
	"goUmbral/utils"
	"math/big"
	"testing"
)

func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func TestProof(t *testing.T) {
	apub, apri, _ := curve.KeyGen()

	cipher_before, _ := recrypt.Encrypt(apub, "gzip_test.go", "tmp.txt", big.NewInt(2))

	bpub, _, _ := curve.KeyGen()

	KF, _ := recrypt.ReKeyGen(apri, bpub, 10, 5, big.NewInt(2))

	cipher_after, _ := recrypt.ReEncrypt(KF, cipher_before)

	for i := 0; i < 10; i++ {
		kf := KF[i]
		cf := cipher_after.CF[i]

		E_1 := curve.PointScalarMul(cipher_before.Capsule.E, kf.Rk)
		V_1 := curve.PointScalarMul(cipher_before.Capsule.V, kf.Rk)

		tau, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		E_2 := curve.PointScalarMul(cipher_before.Capsule.E, tau.D)
		V_2 := curve.PointScalarMul(cipher_before.Capsule.V, tau.D)
		U_2 := curve.BigIntMulBase(tau.D)

		// get h = H(E,E_1,E_2,V,V_1,V_2,U,U_1,U_2,T)
		h1 := utils.HashToCurve(
			utils.ConcatBytes(
				utils.ConcatBytes(
					utils.ConcatBytes(
						utils.ConcatBytes(
							utils.ConcatBytes(
								utils.ConcatBytes(
									utils.ConcatBytes(
										utils.ConcatBytes(
											curve.PointToBytes(cipher_before.Capsule.E),
											curve.PointToBytes(E_1)),
										curve.PointToBytes(E_2)),
									curve.PointToBytes(cipher_before.Capsule.V)),
								curve.PointToBytes(V_1)),
							curve.PointToBytes(V_2)),
						curve.PointToBytes(kf.U_1)),
					curve.PointToBytes(U_2)),
				Int64ToBytes(kf.T.Int64())))

		// get rho = tao + h × rk
		rho := math.BigIntAdd(tau.D, math.BigIntMul(h1, kf.Rk))

		// get h = H(E,E_1,E_2,V,V_1,V_2,U,U_1,U_2,T)
		h2 := utils.HashToCurve(
			utils.ConcatBytes(
				utils.ConcatBytes(
					utils.ConcatBytes(
						utils.ConcatBytes(
							utils.ConcatBytes(
								utils.ConcatBytes(
									utils.ConcatBytes(
										utils.ConcatBytes(
											curve.PointToBytes(cipher_before.Capsule.E),
											curve.PointToBytes(cf.E_1)),
										curve.PointToBytes(E_2)),
									curve.PointToBytes(cipher_before.Capsule.V)),
								curve.PointToBytes(cf.V_1)),
							curve.PointToBytes(V_2)),
						curve.PointToBytes(kf.U_1)),
					curve.PointToBytes(U_2)),
				Int64ToBytes(cf.T.Int64())))

		if !curve.PointScalarMul(cipher_before.Capsule.E, rho).Equal(curve.PointScalarAdd(E_2, curve.PointScalarMul(cf.E_1, h2))) {
			t.Log("E验证失败")
		}

		if !curve.PointScalarMul(cipher_before.Capsule.V, rho).Equal(curve.PointScalarAdd(V_2, curve.PointScalarMul(cf.V_1, h2))) {
			t.Log("V验证失败")
		}

		if !curve.BigIntMulBase(rho).Equal(curve.PointScalarAdd(U_2, curve.PointScalarMul(kf.U_1, h2))) {
			t.Log("U验证失败")
		}
	}

	t.Log("验证成功")
}
