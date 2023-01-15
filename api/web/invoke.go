package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type InvokeBody struct {
	Channelid   string   `json:"channelid"`
	Chaincodeid string   `json:"chaincodeid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}

	// defer r.Body.Close()

	// jsonDecode := json.NewDecoder(r.Body)
	// jsonDecode.DisallowUnknownFields()
	// var invokeBody InvokeBody
	// err := jsonDecode.Decode(&invokeBody)
	// if err != nil {
	// 	fmt.Printf("json decode error : %v", err)
	// }

	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}
