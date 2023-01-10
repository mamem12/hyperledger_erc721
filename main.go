/*
 */

package main

import (
	"hyperledger_erc721/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	nftContract := new(chaincode.TokenERC721Contract)
	nftContract.Info.Version = "0.0.1"
	nftContract.Info.Description = "ERC-721 fabric develop"
	nftContract.Info.License = new(metadata.LicenseMetadata)
	nftContract.Info.License.Name = "None"
	nftContract.Info.Contact = new(metadata.ContactMetadata)
	nftContract.Info.Contact.Name = "None"

	chaincode, err := contractapi.NewChaincode(nftContract)
	chaincode.Info.Title = "ERC-721 chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from TokenERC721Contract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
