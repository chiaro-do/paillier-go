package paillier

import (
	"crypto/rand"
	"math/big"

	"paillier-go/utils"
)

func KeyGen(bits int) (*PublicKey, *PrivateKey) {
	p, _ := rand.Prime(rand.Reader, bits/2)
	q, _ := rand.Prime(rand.Reader, bits/2)

	n := new(big.Int).Mul(p, q)
	n2 := new(big.Int).Mul(n, n)
	g := new(big.Int).Add(n, big.NewInt(1))

	p1 := new(big.Int).Sub(p, big.NewInt(1))
	q1 := new(big.Int).Sub(q, big.NewInt(1))

	lambda := utils.LCM(p1, q1)

	u := new(big.Int).Exp(g, lambda, n2)
	Lu := utils.L(u, n)

	mu := new(big.Int).ModInverse(Lu, n)

	return &PublicKey{n, g, n2}, &PrivateKey{lambda, mu}
}