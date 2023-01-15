package model

type Approval struct {
	Owner    string `json:"owner"`
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}

func NewApproval(owner, operator string, approved bool) *Approval {
	return &Approval{
		Owner:    owner,
		Operator: operator,
		Approved: approved,
	}
}

func (a *Approval) GetOwner() *string {
	return &a.Owner
}

func (a *Approval) GetOperator() *string {
	return &a.Operator
}
func (a *Approval) GetApproved() *bool {
	return &a.Approved
}
