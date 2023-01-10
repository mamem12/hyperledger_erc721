package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hyperledger_erc721/model"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
TransferFrom is invoke fnc that moves token
from is the owner's address, to is reciepient's address
*/
func (c *TokenERC721Contract) TransferFrom(ctx contractapi.TransactionContextInterface, from, to, tokenId string) (bool, error) {

	initialized, err := checkInitialized(ctx)

	if err != nil {
		return false, err
	}

	if !initialized {
		return false, fmt.Errorf("initialized first")
	}

	sender64, err := ctx.GetClientIdentity().GetID()

	if err != nil {
		return false, fmt.Errorf("failed to GetClientIdentity : %v", err)
	}

	senderBytes, err := base64.StdEncoding.DecodeString(sender64)

	if err != nil {
		return false, fmt.Errorf("failed to DecodeString sender ID : %v", err)
	}

	sender := string(senderBytes)

	nft, err := _readNFT(ctx, tokenId)

	if err != nil {
		return false, fmt.Errorf("failed to _readNFT : %v", err)
	}

	owner := nft.Owner
	operator := nft.Approved
	operatorApproval, err := c.IsApprovedForAll(ctx, owner, sender)

	if err != nil {
		return false, fmt.Errorf("failed to get IsApprovedForAll : %v", err)
	}

	if owner != sender && operator != sender && !operatorApproval {
		return false, fmt.Errorf("the sender is not the current owner nor an authorized operator")
	}

	// Check if `from` is the current owner
	if owner != from {
		return false, fmt.Errorf("the from is not the current owner")
	}

	// Clear the approved client for this non-fungible token
	nft.Approved = ""

	// Overwrite a non-fungible token to assign a new owner.
	nft.Owner = to
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})

	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey: %v", err)
	}

	nftBytes, err := json.Marshal(nft)
	if err != nil {
		return false, fmt.Errorf("failed to marshal approval: %v", err)
	}

	err = ctx.GetStub().PutState(nftKey, nftBytes)
	if err != nil {
		return false, fmt.Errorf("failed to PutState nftBytes %s: %v", nftBytes, err)
	}

	// Remove a composite key from the balance of the current owner
	balanceKeyFrom, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{from, tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey from: %v", err)
	}

	err = ctx.GetStub().DelState(balanceKeyFrom)
	if err != nil {
		return false, fmt.Errorf("failed to DelState balanceKeyFrom %s: %v", nftBytes, err)
	}

	// Save a composite key to count the balance of a new owner
	balanceKeyTo, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{to, tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey to: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKeyTo, []byte{0})
	if err != nil {
		return false, fmt.Errorf("failed to PutState balanceKeyTo %s: %v", balanceKeyTo, err)
	}

	// Emit the Transfer event
	transferEvent := model.NewTransferMetadata(from, to, tokenId)

	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return false, fmt.Errorf("failed to marshal transferEventBytes: %v", err)
	}

	err = ctx.GetStub().SetEvent(TransferEventKey, transferEventBytes)
	if err != nil {
		return false, fmt.Errorf("failed to SetEvent transferEventBytes %s: %v", transferEventBytes, err)
	}

	return true, nil
}
