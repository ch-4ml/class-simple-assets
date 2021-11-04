#!/bin/bash

# 설치
docker exec cli peer chaincode install -n simpleasset -v 1.1 -p github.com/simpleasset/v1.1

# 업그래이드 init 
docker exec cli peer chaincode upgrade -n simpleasset -v 1.1 -c '{"Args":[]}' -C mychannel -P 'AND ("Org1MSP.member")'
sleep 3

docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","c","300"]}'
sleep 3

docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","d","400"]}'
sleep 3
# invoke - transfer "c" "d" "20"
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["transfer","c","d","20"]}'
sleep 3

# 테스트 query - get "c"
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","c"]}'
# query - get "d"
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","d"]}'
# query - history "d"
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["history","d"]}'