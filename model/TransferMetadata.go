package model

type Transfer struct {
	From    string `json:"from"`
	To      string `json:"to"`
	TokenId string `json:"tokenId"`
}

func NewTransferMetadata(from, to, tokenId string) *Transfer {

	return &Transfer{
		From:    from,
		To:      to,
		TokenId: tokenId,
	}
}

func (t *Transfer) GetFrom() *string {
	return &t.From
}

func (t *Transfer) GetTo() *string {
	return &t.To
}

func (t *Transfer) GetTokenId() *string {
	return &t.TokenId
}
