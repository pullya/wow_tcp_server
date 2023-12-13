package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockery --name=Challenger --output=mocks --case=underscore
type Challenger interface {
	IsValid(challenge string, response uint64) bool
	Difficulty() int
}

type Challenge struct {
	difficulty int
}

func NewChallenge(difficulty int) Challenge {
	return Challenge{
		difficulty: difficulty,
	}
}

func (c Challenge) Difficulty() int {
	return c.difficulty
}

func (c Challenge) IsValid(challenge string, nonce uint64) bool {
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(int64(256-c.difficulty)), nil)

	hash := crypto.Keccak256([]byte(fmt.Sprint(challenge, nonce)))
	hashInt := new(big.Int).SetBytes(hash)

	return hashInt.Cmp(target) == -1
}
