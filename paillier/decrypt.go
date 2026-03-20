package paillier

import (
	"math/big"

	"paillier-go/utils"
)

func Decrypt(pk *PublicKey, sk *PrivateKey, c *big.Int) *big.Int {
	u := new(big.Int).Exp(c, sk.Lambda, pk.N2)
	Lu := utils.L(u, pk.N)

	m := new(big.Int).Mul(Lu, sk.Mu)
	return m.Mod(m, pk.N)
}