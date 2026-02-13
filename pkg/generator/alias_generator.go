package generator

type AliasGenerator struct {
	aliasLength int
}

func NewAliasGenerator(aliasLength int) *AliasGenerator {
	return &AliasGenerator{
		aliasLength: aliasLength,
	}
}

func (g *AliasGenerator) Generate(s string) string {
	return ""
}
