package utils

import "math/big"

func LCM(a, b *big.Int) *big.Int {
	gcd := new(big.Int).GCD(nil, nil, a, b)
	abs := new(big.Int).Mul(a, b)
	return abs.Div(abs, gcd)
}

func L(u, n *big.Int) *big.Int {
	tmp := new(big.Int).Sub(u, big.NewInt(1))
	return tmp.Div(tmp, n)
}

func IsCoprime(a, b *big.Int) bool {
	g := new(big.Int).GCD(nil, nil, a, b)
	return g.Cmp(big.NewInt(1)) == 0
}