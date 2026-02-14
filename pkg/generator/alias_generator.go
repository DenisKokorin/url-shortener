package generator

import (
	"crypto/rand"
	"math/big"
)

var (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ0123456789_"
)

type AliasGenerator struct {
	aliasLength int
}

func NewAliasGenerator(aliasLength int) *AliasGenerator {
	return &AliasGenerator{
		aliasLength: aliasLength,
	}
}

func (g *AliasGenerator) Generate() string {
	var res string

	for range g.aliasLength {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return ""
		}
		res += string(chars[idx.Int64()])
	}

	return res
}
