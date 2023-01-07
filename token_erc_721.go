/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
	"github.com/hyperledger_erc721/chaincode"
)

// https://unsplash.com/photos/nzyzAUsbV0M
func main() {
	nftContract := new(chaincode.ERC721_SmartContract)
	nftContract.Info.Version = "0.0.1"
	nftContract.Info.Description = "My Smart Contract"
	nftContract.Info.License = new(metadata.LicenseMetadata)
	nftContract.Info.License.Name = "Apache-2.0"
	nftContract.Info.Contact = new(metadata.ContactMetadata)
	nftContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(nftContract)
	chaincode.Info.Title = "erc721 chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from ERC721 Contract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}

}
