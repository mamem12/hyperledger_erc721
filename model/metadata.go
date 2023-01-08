package model

type ERC721 struct {
	Name   string
	Symbol string
}

func NewERC721Metadata(name, symbol string) *ERC721 {

	return &ERC721{
		Name:   name,
		Symbol: symbol,
	}
}
