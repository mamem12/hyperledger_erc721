package main

import (
	"fmt"
	"hyperledger_explorer/web"
)

var (
	orgConfig = web.OrgSetup{}
)

func init() {
	cryptoPath := "../network/organizations/peerOrganizations/org1.example.com/"

	orgConfig = web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}
}

func main() {
	orgSetup, err := web.Initialize(orgConfig)

	if err != nil {
		fmt.Printf("Error initializing setup for Org1: \n\n \t :%v", err)
	}

	web.Serve(web.OrgSetup(*orgSetup))
}
