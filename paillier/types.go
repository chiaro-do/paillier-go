package paillier

import "math/big"

type PublicKey struct {
	N  *big.Int
	G  *big.Int
	N2 *big.Int
}

type PrivateKey struct {
	Lambda *big.Int
	Mu     *big.Int
}