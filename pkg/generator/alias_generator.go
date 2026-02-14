package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ0123456789_")
)

type AliasGenerator struct {
	aliasLength int
}

func NewAliasGenerator(aliasLength int) *AliasGenerator {
	return &AliasGenerator{
		aliasLength: aliasLength,
	}
}

func (g *AliasGenerator) Generate() (string, error) {
	var res []rune

	for range g.aliasLength {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("error while generate alias: %w", err)
		}
		res = append(res, chars[idx.Int64()])
	}

	return string(res), nil
}
