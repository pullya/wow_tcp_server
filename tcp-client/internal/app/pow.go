package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockery --name=IChallenger --output=mocks --case=underscore
type IChallenger interface {
	GenerateSolution(ctx context.Context, challenge string) string
	SetPowDifficulty(diff int)
}

type Challenge struct {
	PowDifficulty int
}

func NewChallenge() Challenge {
	return Challenge{}
}

func (c *Challenge) SetPowDifficulty(diff int) {
	c.PowDifficulty = diff
}

func (c *Challenge) GenerateSolution(ctx context.Context, challenge string) string {
	nonce := c.mineEthash(ctx, challenge)
	return fmt.Sprint(nonce)
}

func (c *Challenge) mineEthash(ctx context.Context, challenge string) uint64 {
	nonce := uint64(0)
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(int64(256-c.PowDifficulty)), nil)

	for {
		hash := crypto.Keccak256([]byte(fmt.Sprint(challenge, nonce)))
		hashInt := new(big.Int).SetBytes(hash)

		if hashInt.Cmp(target) == -1 {
			return nonce
		}

		nonce++
	}
}
