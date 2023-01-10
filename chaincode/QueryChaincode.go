package chaincode

import (
	"encoding/json"
	"fmt"
	"hyperledger_erc721/model"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func _readNFT(ctx contractapi.TransactionContextInterface, tokenId string) (*model.NFT, error) {
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})

	if err != nil {
		return nil, fmt.Errorf("failed to CreateCompositeKey %s: %v", tokenId, err)
	}

	nftBytes, err := ctx.GetStub().GetState(nftKey)

	if err != nil {
		return nil, fmt.Errorf("failed to GetState %s: %v", tokenId, err)
	}

	nft := model.NewNFT("", "", "", "")
	err = json.Unmarshal(nftBytes, nft)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal nftBytes: %v", err)
	}

	return nft, nil
}

func _nftExists(ctx contractapi.TransactionContextInterface, tokenId string) bool {
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})
	if err != nil {
		panic("error creating CreateCompositeKey:" + err.Error())
	}

	nftBytes, err := ctx.GetStub().GetState(nftKey)
	if err != nil {
		panic("error GetState nftBytes:" + err.Error())
	}

	return len(nftBytes) > 0
}

func (c *TokenERC721Contract) Name(ctx contractapi.TransactionContextInterface) (string, error) {

	initialized, err := checkInitialized(ctx)

	if err != nil {
		return "", err
	}

	if !initialized {
		return "", fmt.Errorf("initialized first")
	}

	ERC721MetadataBytes, err := ctx.GetStub().GetState(InitialKey)

	if err != nil {
		return "", fmt.Errorf("failed found name")
	}

	ERC721Metadata := model.NewERC721Metadata("", "")

	err = json.Unmarshal(ERC721MetadataBytes, ERC721Metadata)

	if err != nil {
		return "", fmt.Errorf("failed unmarshal")
	}

	return *ERC721Metadata.GetName(), nil
}

func (c *TokenERC721Contract) Symbol(ctx contractapi.TransactionContextInterface) (string, error) {

	initialized, err := checkInitialized(ctx)

	if err != nil {
		return "", err
	}

	if !initialized {
		return "", fmt.Errorf("initialized first")
	}

	ERC721MetadataBytes, err := ctx.GetStub().GetState(InitialKey)

	if err != nil {
		return "", fmt.Errorf("failed found symbol")
	}

	ERC721Metadata := model.NewERC721Metadata("", "")

	err = json.Unmarshal(ERC721MetadataBytes, ERC721Metadata)

	if err != nil {
		return "", fmt.Errorf("failed unmarshal")
	}

	return *ERC721Metadata.GetSymbol(), nil
}

func (c *TokenERC721Contract) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) int {

	initialized, err := checkInitialized(ctx)
	if err != nil {
		panic(err.Error())
	}
	if !initialized {
		panic("first initialized")
	}

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(balancePrefix, []string{owner})
	if err != nil {
		panic("Error creating asset chaincode:" + err.Error())
	}

	// Count the number of returned composite keys
	balance := 0
	for iterator.HasNext() {
		_, err := iterator.Next()
		if err != nil {
			return 0
		}
		balance++

	}
	return balance
}

func checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	ERC721MetadataBytes, err := ctx.GetStub().GetState(InitialKey)

	if err != nil {
		return false, fmt.Errorf("failed to get metadata: %v", err)
	}
	if ERC721MetadataBytes == nil {
		return false, err
	}
	return true, nil
}
