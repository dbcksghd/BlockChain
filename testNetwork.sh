#!/bin/bash

export PATH=${HOME}/fabric-samples/bin:$PATH
export BasicPATH=${HOME}/fabric-samples/test-network
export FABRIC_CFG_PATH=${HOME}/fabric-samples/config


export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${BasicPATH}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${BasicPATH}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

ORDERER_CA=${BasicPATH}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
PEER_CONN_PARMS="--peerAddresses localhost:7051 --tlsRootCertFiles ${BasicPATH}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${BasicPATH}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

echo "TEST1 : Register"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"Register","Args":["yoochanong", "100", "1206"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3

echo "TEST2 : Register"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"Register","Args":["yoochanhong1", "50", "1207"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3


echo "TEST3 : QueryAllUser"
set -x
peer chaincode query -C mychannel -n fabcar -c '{"Args":["QueryAllUser"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt


echo "TEST4 : QueryUser"
set -x
peer chaincode query -C mychannel -n fabcar -c '{"Args":["QueryUser", "1206"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt

echo "TEST5 : MakeBank"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"MakeBank","Args":["1000"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
sleep 3

echo "TEST6 : BorrowMoney"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"BorrowMoney","Args":["100", "1207"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txts
sleep 3

echo "TEST7 : BorrowMoney"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"BorrowMoney","Args":["100", "1207"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txts
sleep 3

echo "TEST8 : BorrowMoney"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"BorrowMoney","Args":["100", "1207"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txts
sleep 3

echo "TEST9 : TurnRoulette"
set -x
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n fabcar $PEER_CONN_PARMS -c '{"function":"TurnRoulette","Args":["60", "1206"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txts
sleep 3

echo "TEST10 : QueryAllUser"
set -x
peer chaincode query -C mychannel -n fabcar -c '{"Args":["QueryAllUser"]}' >&log.txt
{ set +x; } 2>/dev/null
cat log.txt
