package paillier

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type PublicKey struct {
	N        *big.Int
	NSquare  *big.Int
	G        *big.Int
}

type PrivateKey struct {
	Lambda *big.Int
	Mu     *big.Int
	Pub    *PublicKey
}

func KeyGen(bits int) (*PublicKey, *PrivateKey, error) {

	p, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, nil, err
	}

	q, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, nil, err
	}

	n := new(big.Int).Mul(p, q)
	nsquare := new(big.Int).Mul(n, n)

	g := new(big.Int).Add(n, big.NewInt(1))

	// lambda = lcm(p-1, q-1)
	pm1 := new(big.Int).Sub(p, big.NewInt(1))
	qm1 := new(big.Int).Sub(q, big.NewInt(1))

	lambda := lcm(pm1, qm1)

	// mu
	u := new(big.Int).Exp(g, lambda, nsquare)
	Lu := L(u, n)

	mu := new(big.Int).ModInverse(Lu, n)
	if mu == nil {
		return nil, nil, errors.New("failed to compute mu")
	}

	pub := &PublicKey{
		N:       n,
		NSquare: nsquare,
		G:       g,
	}

	priv := &PrivateKey{
		Lambda: lambda,
		Mu:     mu,
		Pub:    pub,
	}

	return pub, priv, nil
}

func Encrypt(pub *PublicKey, m *big.Int) (*big.Int, error) {

	r, err := rand.Int(rand.Reader, pub.N)
	if err != nil {
		return nil, err
	}

	for r.Cmp(big.NewInt(0)) == 0 {
		r, _ = rand.Int(rand.Reader, pub.N)
	}

	gm := new(big.Int).Exp(pub.G, m, pub.NSquare)
	rn := new(big.Int).Exp(r, pub.N, pub.NSquare)

	c := new(big.Int).Mul(gm, rn)
	c.Mod(c, pub.NSquare)

	return c, nil
}

func Decrypt(priv *PrivateKey, c *big.Int) (*big.Int, error) {

	u := new(big.Int).Exp(c, priv.Lambda, priv.Pub.NSquare)
	Lu := L(u, priv.Pub.N)

	m := new(big.Int).Mul(Lu, priv.Mu)
	m.Mod(m, priv.Pub.N)

	return m, nil
}

func Add(pub *PublicKey, c1, c2 *big.Int) *big.Int {
	result := new(big.Int).Mul(c1, c2)
	result.Mod(result, pub.NSquare)
	return result
}

func Mul(pub *PublicKey, c *big.Int, k *big.Int) *big.Int {
	result := new(big.Int).Exp(c, k, pub.NSquare)
	return result
}

func L(u, n *big.Int) *big.Int {
	u.Sub(u, big.NewInt(1))
	return u.Div(u, n)
}

func lcm(a, b *big.Int) *big.Int {
	gcd := new(big.Int).GCD(nil, nil, a, b)
	abs := new(big.Int).Mul(a, b)
	return abs.Div(abs, gcd)
}