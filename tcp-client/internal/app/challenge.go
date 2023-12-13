package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockery --name=Challenger --output=mocks --case=underscore
type Challenger interface {
	GenerateSolution(ctx context.Context, challenge string) string
	SetDifficulty(diff int)
}

type Challenge struct {
	difficulty int
}

func NewChallenge() Challenge {
	return Challenge{}
}

func (c *Challenge) SetDifficulty(diff int) {
	c.difficulty = diff
}

func (c *Challenge) GenerateSolution(ctx context.Context, challenge string) string {
	nonce := c.mineEthash(ctx, challenge)
	return fmt.Sprint(nonce)
}

func (c *Challenge) mineEthash(ctx context.Context, challenge string) uint64 {
	nonce := uint64(0)
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(int64(256-c.difficulty)), nil)

	for {
		hash := crypto.Keccak256([]byte(fmt.Sprint(challenge, nonce)))
		hashInt := new(big.Int).SetBytes(hash)

		if hashInt.Cmp(target) == -1 {
			return nonce
		}

		nonce++
	}
}
