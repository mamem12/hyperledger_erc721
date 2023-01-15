package chaincode

import (
	"encoding/json"
	"fmt"
	"hyperledger_erc721/chaincode/model"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define objectType names for prefix
const balancePrefix = "balance"
const nftPrefix = "nft"
const approvalPrefix = "approval"

// SetEvent() key
const (
	TransferEventKey       = "Transfer"
	ApprovalForAllEventKey = "ApprovalForAll"
)

// Define key names for options
const InitialKey = "initial"

// TokenERC721Contract contract for managing CRUD operations
type TokenERC721Contract struct {
	contractapi.Contract
}

// ============== ERC721 enumeration extension ===============
//
// param {String} name The name of the token
// param {String} symbol The symbol of the token
/*
`Initialize` is set information for a token and intialize contract.
*/
func (c *TokenERC721Contract) Initialize(ctx contractapi.TransactionContextInterface, name string, symbol string) (bool, error) {
	// Check minter authorization - this sample assumes Org1 is the issuer with privilege to set the name and symbol
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get clientMSPID: %v", err)
	}
	if clientMSPID != "Org1MSP" {
		return false, fmt.Errorf("client is not authorized to set the name and symbol of the token")
	}

	bytes, err := ctx.GetStub().GetState(InitialKey)
	if err != nil {
		return false, fmt.Errorf("failed to get Metadata: %v", err)
	}
	if bytes != nil {
		return false, fmt.Errorf("contract options are already set, client is not authorized to change them")
	}

	ERC721Metadata := model.NewERC721Metadata(name, symbol)

	ERC721MetadataBytes, err := json.Marshal(ERC721Metadata)

	if err != nil {
		return false, fmt.Errorf("failed marshal name : %s, symbol : %s", name, symbol)
	}

	err = ctx.GetStub().PutState(InitialKey, ERC721MetadataBytes)

	if err != nil {
		return false, fmt.Errorf("failed putstate : %v", ERC721Metadata)
	}

	// err = ctx.GetStub().PutState(nameKey, []byte(name))
	// if err != nil {
	// 	return false, fmt.Errorf("failed to PutState nameKey %s: %v", nameKey, err)
	// }

	// err = ctx.GetStub().PutState(symbolKey, []byte(symbol))
	// if err != nil {
	// 	return false, fmt.Errorf("failed to PutState symbolKey %s: %v", symbolKey, err)
	// }

	return true, nil
}

/*
`Name` is returns a descriptive name for a collection of non-fungible tokens in this contract
*/
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

/*
`Symbol` is returns an abbreviated name for non-fungible tokens in this contract.
*/
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

func (c *TokenERC721Contract) Test(ctx contractapi.TransactionContextInterface) (string, error) {
	// x509::CN=org1admin,OU=admin,O=Hyperledger,ST=North Carolina,C=US::CN=fabric-ca-server,OU=Fabric,O=Hyperledger,ST=North Carolina,C=US
	ctx.GetClientIdentity().GetMSPID()
	ctx.GetClientIdentity().GetX509Certificate()

	fmt.Println(ctx.GetStub().GetChannelID())

	return "", nil
}

/*
Checks that contract options have been already initialized
*/
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
