package model

type ERC721Metadata struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

func NewERC721Metadata(name, symbol string) *ERC721Metadata {
	return &ERC721Metadata{Name: name, Symbol: symbol}
}

func (e *ERC721Metadata) GetName() *string {
	return &e.Name
}

func (e *ERC721Metadata) GetSymbol() *string {
	return &e.Symbol
}
