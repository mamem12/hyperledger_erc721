package chaincode

// Define structs to be used by chaincode

type Approval struct {
	Owner    string `json:"owner"`
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}
