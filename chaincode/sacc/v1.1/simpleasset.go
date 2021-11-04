// 1. package
package main

// 2. 외부모듈 포함
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 3. SimpleAsset 구조체 정의
type SimpleAsset struct {
}

type Account struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 4. Init 함수 구현
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// 5. Invoke 함수 구현
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "set" {
		return t.Set(stub, args)
	} else if fn == "get" {
		return t.Get(stub, args)
	} else if fn == "transfer" {
		return t.Transfer(stub, args)
	} else if fn == "history" {
		return t.History(stub, args)
	} else {
		return shim.Error("Not supported function name.")
	}
}

// 6. set 함수 구현 params key, value (2개)
func (t *SimpleAsset) Set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and value")
	}
	var account = Account{Key: args[0], Value: args[1]}
	accAsBytes, _ := json.Marshal(account)

	err := stub.PutState(args[0], accAsBytes)
	if err != nil {
		return shim.Error("Failed to set asset:" + args[0])
	}
	return shim.Success(accAsBytes)
}

// 6. transfer 함수 구현 params key from, key to, value (3개)
func (t *SimpleAsset) Transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect arguments. Expecting from_key, to_key and a value")
	}
	// 1. from key GetState
	fromAsBytes, _ := stub.GetState(args[0])

	// 2. to key GetState
	toAsBytes, _ := stub.GetState(args[1])

	// (TODO) 3. error 체크 from, to가 유효한지
	if fromAsBytes == nil || toAsBytes == nil {
		return shim.Error("Asset not found")
	}

	// unmarshall from, to
	from := Account{}
	to := Account{}

	json.Unmarshal(fromAsBytes, &from)
	json.Unmarshal(toAsBytes, &to)

	// 정보 변경 -
	amount, _ := strconv.Atoi(args[2])
	from_value, _ := strconv.Atoi(from.Value)
	to_value, _ := strconv.Atoi(to.Value)

	// (TODO) 유효성검증 -> from -> to value만큼 이동가능한지
	if from_value < amount {
		return shim.Error("Not enough amount in from account")
	}

	from_value -= amount
	to_value += amount

	from.Value = strconv.Itoa(from_value)
	to.Value = strconv.Itoa(to_value)

	// marshal from, to
	fromAsBytes, _ = json.Marshal(from)
	toAsBytes, _ = json.Marshal(to)

	// from PutState
	stub.PutState(from.Key, fromAsBytes)

	// to PutState
	stub.PutState(to.Key, toAsBytes)

	return shim.Success([]byte("Transfer TX excuted"))
}

// 7. get 함수 구현 params key (1개)
func (t *SimpleAsset) Get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting a key ")
	}
	value, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get asset: " + args[0] + " with error: " + err.Error())
	}
	if value == nil {
		return shim.Error("Asset not found: " + args[0])
	}
	return shim.Success([]byte(value))
}

func (t *SimpleAsset) History(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carName := args[0]

	fmt.Printf("- start History: %s\n", carName)

	resultsIterator, err := stub.GetHistoryForKey(carName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the car
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON car)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- History returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// 8. main 함수 구현
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
