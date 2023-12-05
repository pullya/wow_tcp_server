package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type IChallenger interface {
	GenerateSolution(challenge string) string
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

func (c *Challenge) GenerateSolution(challenge string) string {
	nonce := c.mineEthash(challenge)
	return fmt.Sprint(nonce)
}

func (c *Challenge) mineEthash(challenge string) uint64 {
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
