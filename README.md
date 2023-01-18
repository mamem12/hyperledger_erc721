hyperledger fabric/samples erc721 refactoring

source is :
https://github.com/hyperledger/fabric-samples/tree/main/token-erc-721/chaincode-go

하이퍼레저 패브릭 샘플 erc721 기반으로 만들어져 있습니다.

fabric-samples/test-network를 이용하여 실행할 수 있습니다.

아래 명령어로 실행과 동시에 채널을 만들 수 있습니다.
```
cd test-network/
./network.sh up createChannel
```
아래 명령어로 체인코드를 네트워크에 배포할 수 있습니다.
```
./network.sh deployCC -ccn token_erc721 -ccp ../chaincode -ccl go
```
fabric-samples의 api를 실행하여 네트워크와 통신할 수 있습니다.
```
cd ../api
go run main.go
```
