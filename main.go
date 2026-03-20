package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PublicKey struct {
	N  *big.Int
	G  *big.Int
	N2 *big.Int
}

type PrivateKey struct {
	Lambda *big.Int
	Mu     *big.Int
}

func lcm(a, b *big.Int) *big.Int {
	gcd := new(big.Int).GCD(nil, nil, a, b)
	abs := new(big.Int).Mul(a, b)
	return abs.Div(abs, gcd)
}

func L(u, n *big.Int) *big.Int {
	tmp := new(big.Int).Sub(u, big.NewInt(1))
	return tmp.Div(tmp, n)
}

func KeyGen(bits int) (*PublicKey, *PrivateKey) {
	p, _ := rand.Prime(rand.Reader, bits/2)
	q, _ := rand.Prime(rand.Reader, bits/2)

	n := new(big.Int).Mul(p, q)
	n2 := new(big.Int).Mul(n, n)

	g := new(big.Int).Add(n, big.NewInt(1))

	p1 := new(big.Int).Sub(p, big.NewInt(1))
	q1 := new(big.Int).Sub(q, big.NewInt(1))

	lambda := lcm(p1, q1)

	u := new(big.Int).Exp(g, lambda, n2)
	Lu := L(u, n)

	mu := new(big.Int).ModInverse(Lu, n)

	return &PublicKey{n, g, n2}, &PrivateKey{lambda, mu}
}

func Encrypt(pk *PublicKey, m *big.Int) *big.Int {
	r, _ := rand.Int(rand.Reader, pk.N)

	gm := new(big.Int).Exp(pk.G, m, pk.N2)
	rn := new(big.Int).Exp(r, pk.N, pk.N2)

	c := new(big.Int).Mul(gm, rn)
	return c.Mod(c, pk.N2)
}

func Decrypt(pk *PublicKey, sk *PrivateKey, c *big.Int) *big.Int {
	u := new(big.Int).Exp(c, sk.Lambda, pk.N2)
	Lu := L(u, pk.N)

	m := new(big.Int).Mul(Lu, sk.Mu)
	return m.Mod(m, pk.N)
}

func main() {
	pk, sk := KeyGen(512)

	// 20 categories initialized C_i = E(0)
	C := make([]*big.Int, 20)
	for i := 0; i < 20; i++ {
		C[i] = Encrypt(pk, big.NewInt(0))
	}

	plaintextCounts := make([]int, 20)

	updates := []int{1, 5, 5, 2, 1, 5, 0, 19, 5, 1, 1}

	for _, k := range updates {
		encOne := Encrypt(pk, big.NewInt(1))

		C[k] = new(big.Int).Mul(C[k], encOne)
		C[k].Mod(C[k], pk.N2)

		plaintextCounts[k]++
	}

	// Decrypt and compare
	fmt.Println("Category | Expected | Decrypted")

	for i := 0; i < 20; i++ {
		dec := Decrypt(pk, sk, C[i])
		fmt.Printf("%8d | %8d | %s\n", i, plaintextCounts[i], dec.String())
	}
}