export PATH=/root/go/src/github.com/hyperledger/test-simple-network/bin:$PATH
export FABRIC_CFG_PATH=/root/go/src/github.com/hyperledger/test-simple-network/config/

export CORE_PEER_TLS_ENABLED=true

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/root/go/src/github.com/hyperledger/test-simple-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/root/go/src/github.com/hyperledger/test-simple-network/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051

export ORDERER_CA_PATH=/root/go/src/github.com/hyperledger/test-simple-network/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem


#create channel
peer channel create -o orderer.example.com:7050 -c mychannel -f channel-artifacts/channel.tx --tls --cafile ${ORDERER_CA_PATH}
#peer0.org1 join channel
peer channel join -b mychannel.block

#install chaincode on peer0.org1
peer chaincode install -n mycc -v 1.0 -p chaincode/htlc

#instantiate chaincode on peer0.org1
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -v 1.0 -c '{"Args":["Init"]}' -P "AND ('Org1MSP.peer')"

sleep 10s

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile ${ORDERER_CA_PATH} -C mychannel -n mycc -c '{"Args":["InitLedger"]}'

peer chaincode query -C mychannel -n mycc -c '{"Args":["GetBalance","Alice"]}'
