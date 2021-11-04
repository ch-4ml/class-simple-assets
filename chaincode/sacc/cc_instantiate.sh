#!/bin/bash

# 설치
docker exec cli peer chaincode install -n simpleasset -v 1.0 -p github.com/simpleasset/v1.0
# 배포 init "a" "100"
docker exec cli peer chaincode instantiate -n simpleasset -v 1.0 -c '{"Args":["a","100"]}' -C mychannel -P 'AND ("Org1MSP.member")'
sleep 3

# 테스트 query - get "a"
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","a"]}'
# invoke - set "b" "200"
docker exec cli peer chaincode invoke -n simpleasset -C mychannel -c '{"Args":["set","b","200"]}'
sleep 3
# query - get "b"
docker exec cli peer chaincode query -n simpleasset -C mychannel -c '{"Args":["get","b"]}'