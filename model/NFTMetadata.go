package model

type NFT struct {
	TokenId  string `json:"tokenId"`
	Owner    string `json:"owner"`
	TokenURI string `json:"tokenURI"`
	Approved string `json:"approved"`
}

func NewNFT(tokenId, owner, tokenURI, approved string) *NFT {
	return &NFT{
		TokenId:  tokenId,
		Owner:    owner,
		TokenURI: tokenURI,
		Approved: approved,
	}
}

func (n *NFT) GetTokenId() *string {
	return &n.TokenId
}

func (n *NFT) GetOwner() *string {
	return &n.Owner
}

func (n *NFT) GetTokenURI() *string {
	return &n.TokenURI
}

func (n *NFT) GetApproved() *string {
	return &n.Approved
}
