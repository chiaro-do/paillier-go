package main

import (
	"fmt"
	"math/big"
	"paillier-go/utils"
	"paillier-go/paillier"
)

const cats = 20

func main() {
	pk, sk := paillier.KeyGen(512)

	C := make([]*big.Int, cats)
	for i := 0; i < cats; i++ {
		C[i] = paillier.Encrypt(pk, big.NewInt(0))
	}
	filePath := "crud.csv"
	records := utils.ReadCSV(filePath)
	plaintextCounts := make([]int, cats)

	for _, record := range records {
		encOne := paillier.Encrypt(pk, big.NewInt(1))
		k := record[2].(int)
		C[k] = new(big.Int).Mul(C[k], encOne)
		C[k].Mod(C[k], pk.N2)

		plaintextCounts[k]++
	}
	
	// Decrypt and compare
	fmt.Println("Category | Expected | Decrypted")
	sum := 0
	for i := 0; i < cats; i++ {
		dec := paillier.Decrypt(pk, sk, C[i])
		sum += plaintextCounts[i]
		fmt.Printf("%8d | %8d | %s\n", i, plaintextCounts[i], dec.String())
	}
	fmt.Printf("sum: %8d\n", sum)
}