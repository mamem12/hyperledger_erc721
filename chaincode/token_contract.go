package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger_erc721/model"
)

// Define key names for options
const balancePrefix = "balance"
const nftPrefix = "nft"
const approvalPrefix = "approval"

type ERC721_SmartContract struct {
	contractapi.Contract
}

func _readNFT(ctx contractapi.TransactionContextInterface, tokenId string) (*model.Nft, error) {
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})

	if err != nil {
		return nil, fmt.Errorf("failed to CreateCompositeKey %s: %v", tokenId, err)
	}

	nftBytes, err := ctx.GetStub().GetState(nftKey)

	if err != nil {
		return nil, fmt.Errorf("failed to GetState %s: %v", tokenId, err)
	}

	nft := new(model.Nft)
	err = json.Unmarshal(nftBytes, nft)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal nftBytes: %v", err)
	}

	return nft, nil
}

func _nftExists(ctx contractapi.TransactionContextInterface, tokenId string) bool {
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})
	if err != nil {
		panic(err.Error())
	}
	nftBytes, err := ctx.GetStub().GetState(nftKey)

	if err != nil {
		panic("error GetState nftBytes:" + err.Error())
	}

	return len(nftBytes) > 0
}

func (c *ERC721_SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) int {

	// 토큰 컨트랙트가 초기화 되어있는지 확인
	initialized, err := checkInitialized(ctx)

	if !initialized {
		panic(err.Error())
	}

	// BalanceOf() 함수는 balancePrefix.owner.*과 일치하는 모든 레코드를 계산함

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(balancePrefix, []string{owner})
	if err != nil {
		panic("Error creating asset chaincode:" + err.Error())
	}

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

func (c *ERC721_SmartContract) OwnerOf(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {

	initialized, err := checkInitialized(ctx)

	if !initialized {
		return "", err
	}

	nft, err := _readNFT(ctx, tokenId)

	if err != nil {
		return "", fmt.Errorf("not found token owner tokenId : %s", tokenId)
	}

	return nft.Owner, nil
}

func (c *ERC721_SmartContract) Approve(ctx contractapi.TransactionContextInterface, operator string, tokenId string) (bool, error) {

	initialized, err := checkInitialized(ctx)

	if !initialized {
		return false, err
	}

	sender64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed to GetClientIdentity: %v", err)
	}

	senderBytes, err := base64.StdEncoding.DecodeString(sender64)
	if err != nil {
		return false, fmt.Errorf("failed to DecodeString senderBytes: %v", err)
	}

	sender := string(senderBytes)

	nft, err := _readNFT(ctx, tokenId)
	if err != nil {
		return false, fmt.Errorf("failed to _readNFT: %v", err)
	}

	owner := nft.Owner
	operatorApproval, err := c.IsApprovedForAll(ctx, owner, sender)

	if err != nil {
		return false, fmt.Errorf("failed to get IsApprovedForAll: %v", err)
	}

	if owner != sender && !operatorApproval {
		return false, fmt.Errorf("the sender is not the current owner nor an authorized operator")
	}

	nft.Approved = operator
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey %s: %v", nftKey, err)
	}

	nftBytes, err := json.Marshal(nft)
	if err != nil {
		return false, fmt.Errorf("failed to marshal nftBytes: %v", err)
	}

	err = ctx.GetStub().PutState(nftKey, nftBytes)
	if err != nil {
		return false, fmt.Errorf("failed to PutState for nftKey: %v", err)
	}

	return true, nil
}

func (c *ERC721_SmartContract) GetApproved(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {

	initialized, err := checkInitialized(ctx)
	if !initialized {
		return "false", err
	}

	nft, err := _readNFT(ctx, tokenId)
	if err != nil {
		return "false", fmt.Errorf("failed GetApproved for tokenId : %v", err)
	}
	return nft.Approved, nil
}

func (c *ERC721_SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, from string, to string, tokenId string) (bool, error) {
	initialized, err := checkInitialized(ctx)

	if !initialized {
		return false, err
	}

	sender64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed to GetClientIdentity: %v", err)
	}

	senderBytes, err := base64.StdEncoding.DecodeString(sender64)
	if err != nil {
		return false, fmt.Errorf("failed to DecodeString sender: %v", err)
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

	// owner와 'from'이 같은지 확인합니다.
	if owner != from {
		return false, fmt.Errorf("the from is not the current owner")
	}

	nft.Approved = ""

	// 소유권 변경
	nft.Owner = to
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey: %v", err)
	}

	nftBytes, err := json.Marshal(nft)
	if err != nil {
		return false, fmt.Errorf("failed to marshal approval: %v", err)
	}

	// 소유권 변경 기록
	err = ctx.GetStub().PutState(nftKey, nftBytes)

	if err != nil {
		return false, fmt.Errorf("failed to PutState : %v", err)
	}

	// 소유자에게서 컴포짓 키 삭제
	balanceKeyFrom, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{from, tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey from: %v", err)
	}

	err = ctx.GetStub().DelState(balanceKeyFrom)
	if err != nil {
		return false, fmt.Errorf("failed to DelState balanceKeyFrom %s: %v", nftBytes, err)
	}

	// 새 소유자에게 컴포짓 키 생성
	balanceKeyTo, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{to, tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey to: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKeyTo, []byte{0})
	if err != nil {
		return false, fmt.Errorf("failed to PutState balanceKeyTo %s: %v", balanceKeyTo, err)
	}

	// 교환 요청 저장
	transferEvent := new(model.Transfer)
	transferEvent.From = from
	transferEvent.To = to
	transferEvent.TokenId = tokenId

	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return false, fmt.Errorf("failed to marshal transferEventBytes: %v", err)
	}

	err = ctx.GetStub().SetEvent("Transfer", transferEventBytes)
	if err != nil {
		return false, fmt.Errorf("failed to SetEvent transferEventBytes %s: %v", transferEventBytes, err)
	}

	return true, nil
}

func (c *ERC721_SmartContract) Name(ctx contractapi.TransactionContextInterface) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)

	if !initialized {
		return "", err
	}

	bytes, err := ctx.GetStub().GetState(initialKey)
	erc721 := model.NewERC721Metadata("", "")
	json.Unmarshal(bytes, erc721)
	if err != nil {
		return "", fmt.Errorf("failed to get Name bytes: %s", err)
	}

	return erc721.Name, nil
}

func (c *ERC721_SmartContract) Symbol(ctx contractapi.TransactionContextInterface) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if !initialized {
		return "", err
	}

	bytes, err := ctx.GetStub().GetState(initialKey)
	erc721 := model.NewERC721Metadata("", "")
	json.Unmarshal(bytes, erc721)

	if err != nil {
		return "", fmt.Errorf("failed to get Symbol: %v", err)
	}

	return erc721.Symbol, nil
}

func (c *ERC721_SmartContract) TokenURI(ctx contractapi.TransactionContextInterface, tokenId string) (string, error) {

	// Check if contract has been intilized first
	initialized, err := checkInitialized(ctx)
	if !initialized {
		return "", err
	}

	nft, err := _readNFT(ctx, tokenId)
	if err != nil {
		return "", fmt.Errorf("failed to get TokenURI: %v", err)
	}
	return nft.TokenURI, nil
}

func (c *ERC721_SmartContract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string) (bool, error) {

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get clientMSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return false, fmt.Errorf("client is not authorized to set the name and symbol of the token")
	}

	bytes, err := ctx.GetStub().GetState(initialKey)
	if err != nil {
		return false, fmt.Errorf("failed to get Name: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	erc721 := model.NewERC721Metadata(name, symbol)
	erc721Bytes, err := json.Marshal(erc721)

	if err != nil {
		return false, fmt.Errorf("failed erc721 marshal")
	}

	err = ctx.GetStub().PutState("", erc721Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to putstate %v", erc721)
	}

	return true, nil
}

func (c *ERC721_SmartContract) MintWithTokenURI(ctx contractapi.TransactionContextInterface, tokenId string, tokenURI string) (*model.Nft, error) {

	initialized, err := checkInitialized(ctx)

	if !initialized {
		return nil, err
	}

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get clientMSPID: %v", err)
	}

	if clientMSPID != "Org1MSP" {
		return nil, fmt.Errorf("client is not authorized to set the name and symbol of the token")
	}

	minter64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get minter id: %v", err)
	}

	minterBytes, err := base64.StdEncoding.DecodeString(minter64)
	if err != nil {
		return nil, fmt.Errorf("failed to DecodeString minter64: %v", err)
	}
	minter := string(minterBytes)

	exists := _nftExists(ctx, tokenId)

	if exists {
		return nil, fmt.Errorf("exists token uri : %v", tokenURI)
	}

	nft := new(model.Nft)
	nft.Owner = minter
	nft.TokenURI = tokenURI
	nft.TokenId = tokenId

	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})

	if err != nil {
		return nil, err
	}

	nftBytes, err := json.Marshal(nft)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v", nft)
	}

	err = ctx.GetStub().PutState(nftKey, nftBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to putstate %v", nft)
	}

	balanceKey, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{minter, tokenId})
	if err != nil {
		return nil, fmt.Errorf("failed to CreateCompositeKey to balanceKey: %v", err)
	}

	err = ctx.GetStub().PutState(balanceKey, []byte{'\u0000'})
	if err != nil {
		return nil, fmt.Errorf("failed to PutState balanceKey %s: %v", nftBytes, err)
	}

	transferEvent := new(model.Transfer)
	transferEvent.From = "0x0"
	transferEvent.To = minter
	transferEvent.TokenId = tokenId

	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transferEventBytes: %v", err)
	}

	err = ctx.GetStub().SetEvent("Transfer", transferEventBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to SetEvent transferEventBytes %s: %v", transferEventBytes, err)
	}

	return &model.Nft{}, nil
}

func (c *ERC721_SmartContract) Burn(ctx contractapi.TransactionContextInterface, tokenId string) (bool, error) {
	initialized, err := checkInitialized(ctx)

	if !initialized {
		return false, err
	}

	owner64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed to GetClientIdentity owner64: %v", err)
	}

	ownerBytes, err := base64.StdEncoding.DecodeString(owner64)
	if err != nil {
		return false, fmt.Errorf("failed to DecodeString owner64: %v", err)
	}

	owner := string(ownerBytes)

	nft, err := _readNFT(ctx, tokenId)

	if err != nil {
		return false, fmt.Errorf("failed to _readNFT nft : %v", err)
	}
	if nft.Owner != owner {
		return false, fmt.Errorf("non-fungible token %s is not owned by %s", tokenId, owner)
	}

	// Delete the token
	nftKey, err := ctx.GetStub().CreateCompositeKey(nftPrefix, []string{tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey tokenId: %v", err)
	}

	err = ctx.GetStub().DelState(nftKey)
	if err != nil {
		return false, fmt.Errorf("failed to DelState nftKey: %v", err)
	}

	// Remove a composite key from the balance of the owner
	balanceKey, err := ctx.GetStub().CreateCompositeKey(balancePrefix, []string{owner, tokenId})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey balanceKey %s: %v", balanceKey, err)
	}

	err = ctx.GetStub().DelState(balanceKey)
	if err != nil {
		return false, fmt.Errorf("failed to DelState balanceKey %s: %v", balanceKey, err)
	}

	// Emit the Transfer event
	transferEvent := new(model.Transfer)
	transferEvent.From = owner
	transferEvent.To = "0x0"
	transferEvent.TokenId = tokenId

	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return false, fmt.Errorf("failed to marshal transferEventBytes: %v", err)
	}

	err = ctx.GetStub().SetEvent("Transfer", transferEventBytes)
	if err != nil {
		return false, fmt.Errorf("failed to SetEvent transferEventBytes: %v", err)
	}

	return false, errors.New("")
}

func (c *ERC721_SmartContract) SetApprovalForAll(ctx contractapi.TransactionContextInterface, operator string, approved bool) (bool, error) {

	initialized, err := checkInitialized(ctx)
	if !initialized {
		return false, err
	}

	sender64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return false, fmt.Errorf("failed to GetClientIdentity: %v", err)
	}

	senderBytes, err := base64.StdEncoding.DecodeString(sender64)
	if err != nil {
		return false, fmt.Errorf("failed to DecodeString sender: %v", err)
	}
	sender := string(senderBytes)

	nftApproval := new(model.Approval)
	nftApproval.Owner = sender
	nftApproval.Operator = operator
	nftApproval.Approved = approved

	approvalKey, err := ctx.GetStub().CreateCompositeKey(approvalPrefix, []string{sender, operator})
	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey: %v", err)
	}

	approvalBytes, err := json.Marshal(nftApproval)
	if err != nil {
		return false, fmt.Errorf("failed to marshal approvalBytes: %v", err)
	}

	err = ctx.GetStub().PutState(approvalKey, approvalBytes)
	if err != nil {
		return false, fmt.Errorf("failed to PutState approvalBytes: %v", err)
	}

	// Emit the ApprovalForAll event
	err = ctx.GetStub().SetEvent("ApprovalForAll", approvalBytes)
	if err != nil {
		return false, fmt.Errorf("failed to SetEvent ApprovalForAll: %v", err)
	}

	return true, nil
}

func (c *ERC721_SmartContract) IsApprovedForAll(ctx contractapi.TransactionContextInterface, owner, operator string) (bool, error) {
	initialized, err := checkInitialized(ctx)

	if !initialized {
		return false, err
	}

	approvalKey, err := ctx.GetStub().CreateCompositeKey(approvalPrefix, []string{owner, operator})

	if err != nil {
		return false, fmt.Errorf("failed to CreateCompositeKey: %v", err)
	}
	approvalBytes, err := ctx.GetStub().GetState(approvalKey)
	if err != nil {
		return false, fmt.Errorf("failed to GetState approvalBytes %s: %v", approvalBytes, err)
	}

	if len(approvalBytes) < 1 {
		return false, nil
	}

	approval := new(model.Approval)
	err = json.Unmarshal(approvalBytes, approval)
	if err != nil {
		return false, fmt.Errorf("failed to Unmarshal: %v, string %s", err, string(approvalBytes))
	}

	return approval.Approved, nil
}

func (c *ERC721_SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	initialized, err := checkInitialized(ctx)
	if !initialized {
		return 0, err
	}

	clientAccountID64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to GetClientIdentity minter: %v", err)
	}

	clientAccountIDBytes, err := base64.StdEncoding.DecodeString(clientAccountID64)
	if err != nil {
		return 0, fmt.Errorf("failed to DecodeString sender: %v", err)
	}

	clientAccountID := string(clientAccountIDBytes)

	return c.BalanceOf(ctx, clientAccountID), nil
}

func (c *ERC721_SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	initialized, err := checkInitialized(ctx)

	if !initialized {
		return "", err
	}

	clientAccountID64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to GetClientIdentity minter: %v", err)
	}

	clientAccountBytes, err := base64.StdEncoding.DecodeString(clientAccountID64)
	if err != nil {
		return "", fmt.Errorf("failed to DecodeString clientAccount64: %v", err)
	}
	clientAccount := string(clientAccountBytes)

	return clientAccount, nil
}

func (c *ERC721_SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface) int {

	initialized, err := checkInitialized(ctx)

	if !initialized {
		panic(err.Error())
	}

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(nftPrefix, []string{})

	if err != nil {
		panic("Error creating GetStateByPartialCompositeKey:" + err.Error())
	}

	totalSupply := 0
	for iterator.HasNext() {
		_, err := iterator.Next()
		if err != nil {
			return 0
		}
		totalSupply++

	}
	return totalSupply
}
