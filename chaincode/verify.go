package chaincode

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const nameKey = "name"
const symbolKey = "symbol"

func checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	tokenName, err := ctx.GetStub().GetState(nameKey)
	if err != nil {
		return false, fmt.Errorf("토큰의 이름을 가져오는데 실패했습니다.: %v", err)
	}
	if tokenName == nil {
		return false, errors.New("Initialize() 함수를 이용하여 초기화를 진행해주세요")
	}

	return true, nil
}
