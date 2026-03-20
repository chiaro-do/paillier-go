package test

import (
	"math/big"
	"testing"

	"paillier-go/paillier"
)

func TestHomomorphicAddition(t *testing.T) {
	pk, sk := paillier.KeyGen(512)

	c1 := paillier.Encrypt(pk, big.NewInt(5))
	c2 := paillier.Encrypt(pk, big.NewInt(3))

	// Homomorphic addition
	c := new(big.Int).Mul(c1, c2)
	c.Mod(c, pk.N2)

	result := paillier.Decrypt(pk, sk, c)

	if result.Cmp(big.NewInt(8)) != 0 {
		t.Errorf("Expected 8, got %s", result.String())
	}
}

func TestRepeatedAddition(t *testing.T) {
	pk, sk := paillier.KeyGen(512)

	c := paillier.Encrypt(pk, big.NewInt(0))

	for i := 0; i < 10; i++ {
		encOne := paillier.Encrypt(pk, big.NewInt(1))
		c.Mul(c, encOne)
		c.Mod(c, pk.N2)
	}

	result := paillier.Decrypt(pk, sk, c)

	if result.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected 10, got %s", result.String())
	}
}