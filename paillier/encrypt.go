package paillier

import (
	"crypto/rand"
	"math/big"

	"paillier-go/utils"
)

func randomZnStar(n *big.Int) *big.Int {
	for {
		r, _ := rand.Int(rand.Reader, n)
		if utils.IsCoprime(r, n) {
			return r
		}
	}
}

func Encrypt(pk *PublicKey, m *big.Int) *big.Int {
	r := randomZnStar(pk.N)

	gm := new(big.Int).Exp(pk.G, m, pk.N2)
	rn := new(big.Int).Exp(r, pk.N, pk.N2)

	c := new(big.Int).Mul(gm, rn)
	return c.Mod(c, pk.N2)
}