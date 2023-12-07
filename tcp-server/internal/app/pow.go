package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockery --name=IChallenger --output=mocks --case=underscore
type IChallenger interface {
	IsValidPoW(challenge string, response uint64) bool
	GetPowDifficulty() int
}

type Challenge struct {
	PowDifficulty int
}

func NewChallenge(powDifficulty int) Challenge {
	return Challenge{
		PowDifficulty: powDifficulty,
	}
}

func (c Challenge) GetPowDifficulty() int {
	return c.PowDifficulty
}

func (c Challenge) IsValidPoW(challenge string, nonce uint64) bool {
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(int64(256-c.PowDifficulty)), nil)

	hash := crypto.Keccak256([]byte(fmt.Sprint(challenge, nonce)))
	hashInt := new(big.Int).SetBytes(hash)

	return hashInt.Cmp(target) == -1
}
