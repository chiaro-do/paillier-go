package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	cats          = 20
	publicKeyKey  = "PAILLIER_PUBLIC_KEY"
	counterPrefix = "COUNTER_"
)

type SmartContract struct {
	contractapi.Contract
}

type PublicKey struct {
	N       string `json:"n"`
	NSquare string `json:"nsquare"`
}

func counterKey(k int) string {
	return fmt.Sprintf("%s%d", counterPrefix, k)
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, nStr string) error {
	existing, err := ctx.GetStub().GetState(publicKeyKey)
	if err != nil {
		return fmt.Errorf("failed to read public key state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("ledger already initialized")
	}

	n := new(big.Int)
	if _, ok := n.SetString(nStr, 10); !ok {
		return fmt.Errorf("invalid n")
	}

	nsquare := new(big.Int).Mul(n, n)

	pubKey := PublicKey{
		N:       n.String(),
		NSquare: nsquare.String(),
	}
	pubBytes, err := json.Marshal(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}
	if err := ctx.GetStub().PutState(publicKeyKey, pubBytes); err != nil {
		return fmt.Errorf("failed to store public key: %v", err)
	}

	for i := 0; i < cats; i++ {
		counterState := map[string]string{
			"value": "1",
		}
		bytes, err := json.Marshal(counterState)
		if err != nil {
			return fmt.Errorf("failed to marshal counter %d: %v", i, err)
		}
		if err := ctx.GetStub().PutState(counterKey(i), bytes); err != nil {
			return fmt.Errorf("failed to store counter %d: %v", i, err)
		}
	}

	return nil
}

func (s *SmartContract) UpdateCounter(ctx contractapi.TransactionContextInterface, k int) error {
	if k < 0 || k >= cats {
		return fmt.Errorf("invalid counter index")
	}

	bytes, err := ctx.GetStub().GetState(counterKey(k))
	if err != nil || bytes == nil {
		return fmt.Errorf("counter not found")
	}
	var state map[string]string
	if err := json.Unmarshal(bytes, &state); err != nil {
		return fmt.Errorf("failed to unmarshal counter: %v", err)
	}

	pubBytes, err := ctx.GetStub().GetState(publicKeyKey)
	if err != nil || pubBytes == nil {
		return fmt.Errorf("public key not found")
	}
	var pubKey PublicKey
	if err := json.Unmarshal(pubBytes, &pubKey); err != nil {
		return fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	current := new(big.Int)
	if _, ok := current.SetString(state["value"], 10); !ok {
		return fmt.Errorf("invalid counter value")
	}

	n := new(big.Int)
	if _, ok := n.SetString(pubKey.N, 10); !ok {
		return fmt.Errorf("invalid public key n")
	}
	nsquare := new(big.Int)
	if _, ok := nsquare.SetString(pubKey.NSquare, 10); !ok {
		return fmt.Errorf("invalid public key nsquare")
	}

	eOne := new(big.Int).Add(n, big.NewInt(1))

	// Homomorphic addition: C_k = C_k * E(1) mod n^2
	updated := new(big.Int).Mul(current, eOne)
	updated.Mod(updated, nsquare)

	state["value"] = updated.String()
	newBytes, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal updated counter: %v", err)
	}
	return ctx.GetStub().PutState(counterKey(k), newBytes)
}

func (s *SmartContract) GetCounter(ctx contractapi.TransactionContextInterface, k int) (string, error) {
	if k < 0 || k >= cats {
		return "", fmt.Errorf("invalid counter index")
	}

	bytes, err := ctx.GetStub().GetState(counterKey(k))
	if err != nil || bytes == nil {
		return "", fmt.Errorf("counter not found")
	}

	var state map[string]string
	if err := json.Unmarshal(bytes, &state); err != nil {
		return "", fmt.Errorf("failed to unmarshal counter: %v", err)
	}

	return state["value"], nil
}

func (s *SmartContract) GetAllCounters(ctx contractapi.TransactionContextInterface) (map[int]string, error) {
	counters := make(map[int]string)

	for i := 0; i < cats; i++ {
		val, err := s.GetCounter(ctx, i)
		if err != nil {
			return nil, fmt.Errorf("failed to read counter %d: %v", i, err)
		}
		counters[i] = val
	}

	return counters, nil
}
