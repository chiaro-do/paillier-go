package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"encoding/csv"
	"os"
	"log"
	"strconv"
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

func ReadCSV(filePath string) [][]interface{} {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	records = records[1:]
	var result [][]interface{}

	for _, record := range records {
		a, err := strconv.Atoi(record[5])
		l, err := strconv.Atoi(record[6])
		e, err := strconv.Atoi(record[7])
		if err != nil {
			log.Fatalf("Error parsing age: %v\n", err)
		}
		newRecord := []interface{}{record[0], record[1], mapping(a, l, e)}
		result = append(result, newRecord)
	}
	
	return result
}

func mapping(a int, l int, e int) int {
	return 4 * a + 2 * l + e
}

func main() {
	pk, sk := KeyGen(512)

	// 20 categories initialized C_i = E(0)
	C := make([]*big.Int, 20)
	for i := 0; i < 20; i++ {
		C[i] = Encrypt(pk, big.NewInt(0))
	}
	filePath := "crud.csv"
	records := ReadCSV(filePath)
	plaintextCounts := make([]int, 20)

	for _, record := range records {
		encOne := Encrypt(pk, big.NewInt(1))
		k := record[2].(int)
		C[k] = new(big.Int).Mul(C[k], encOne)
		C[k].Mod(C[k], pk.N2)

		plaintextCounts[k]++
	}
	
	// Decrypt and compare
	fmt.Println("Category | Expected | Decrypted")
	sum := 0
	for i := 0; i < 20; i++ {
		dec := Decrypt(pk, sk, C[i])
		sum += plaintextCounts[i]
		fmt.Printf("%8d | %8d | %s\n", i, plaintextCounts[i], dec.String())
	}
	fmt.Printf("sum: %8d\n", sum)
}