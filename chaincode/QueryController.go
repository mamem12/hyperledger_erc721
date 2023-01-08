package chaincode

import (
	"encoding/json"
	"fmt"
	"hyperledger_erc721/model"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (c *TokenERC721Contract) Name(ctx contractapi.TransactionContextInterface) (string, error) {
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

func checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	ERC721MetadataBytes, err := ctx.GetStub().GetState(InitialKey)

	if err != nil {
		return false, fmt.Errorf("failed to get metadata: %v", err)
	}
	if ERC721MetadataBytes == nil {
		return false, nil
	}
	return true, nil
}
