package paillier

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Structure storing encrypted counters
type EncryptedCounters struct {
	Counters map[int]string `json:"counters"` // ciphertexts as strings
	N        string         `json:"n"`        // public key n
	Nsquare  string         `json:"nsquare"`  // n^2
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, nStr string) error {

	n := new(big.Int)
	n.SetString(nStr, 10)

	nsquare := new(big.Int).Mul(n, n)

	counters := make(map[int]string)

	// Initialize all 20 counters as E(0) = 1 * r^n mod n^2
	for i := 0; i < 20; i++ {
		// For simplicity: initialize as 1 (valid Paillier zero encryption with r=1)
		counters[i] = "1"
	}

	state := EncryptedCounters{
		Counters: counters,
		N:        n.String(),
		Nsquare:  nsquare.String(),
	}

	bytes, _ := json.Marshal(state)
	return ctx.GetStub().PutState("LOG_COUNTERS", bytes)
}

func (s *SmartContract) UpdateCounter(ctx contractapi.TransactionContextInterface, k int, encryptedOne string) error {

	bytes, err := ctx.GetStub().GetState("LOG_COUNTERS")
	if err != nil || bytes == nil {
		return fmt.Errorf("state not found")
	}

	var state EncryptedCounters
	json.Unmarshal(bytes, &state)

	nsquare := new(big.Int)
	nsquare.SetString(state.Nsquare, 10)

	current := new(big.Int)
	current.SetString(state.Counters[k], 10)

	eOne := new(big.Int)
	eOne.SetString(encryptedOne, 10)

	// Homomorphic addition C_k = C_k * E(1) mod (n^2)
	updated := new(big.Int).Mul(current, eOne)
	updated.Mod(updated, nsquare)

	state.Counters[k] = updated.String()

	newBytes, _ := json.Marshal(state)
	return ctx.GetStub().PutState("LOG_COUNTERS", newBytes)
}

func (s *SmartContract) GetCounter(ctx contractapi.TransactionContextInterface, k int) (string, error) {

	bytes, err := ctx.GetStub().GetState("LOG_COUNTERS")
	if err != nil || bytes == nil {
		return "", fmt.Errorf("state not found")
	}

	var state EncryptedCounters
	json.Unmarshal(bytes, &state)

	return state.Counters[k], nil
}

func (s *SmartContract) GetAllCounters(ctx contractapi.TransactionContextInterface) (map[int]string, error) {

	bytes, err := ctx.GetStub().GetState("LOG_COUNTERS")
	if err != nil || bytes == nil {
		return nil, fmt.Errorf("state not found")
	}

	var state EncryptedCounters
	json.Unmarshal(bytes, &state)

	return state.Counters, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}