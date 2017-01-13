/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"crypto/sha256"
	"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/dngroup/content-fabric/content-contract-common"
	"encoding/base64"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type UserContract struct {
	userId             string    `json:"userID"`
	contentId          string    `json:"contentID"`
	//time max after the request is deleted
	timestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	shaUser            string    `json:"sha_user"`
	// random int
	random63           int64     `json:"random63"`
	//use for state
	timestampUser      int64     `json:"timestampUser"`
	timestampBrokering int64     `json:"timestampBrokering"`
}
type EventContract struct {
	typeContract string    `json:"typeContract"`
	sha          string    `json:"sha"`
	Id           string    `json:"ID"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {

	fmt.Println("Starting Content-Contract")
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Content Contract chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "content-brockering-contract" {
		return t.write(stub, args)
	}
	//else if function == "content-licensing-contract" {
	//		return t.licensing(stub, args)
	//	} else if function == "content-delivery-contract" {
	//		return t.delivery(stub, args)
	//	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	//if function == "dummy_query" {
	//	//read a variable
	//	fmt.Println("hi there " + function) //error
	//	return nil, nil
	//} else
	if function == "read" {
		//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

// move - invoke function to change the value of key1 to 0 if is 1 and the value of key2 to 1
func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key1, value1, key2, value2 string
	var err error
	fmt.Println("Runing move")
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. name of the key and value to set both")
	}

	key1 = args[0] //rename for fun
	value1 = args[1]
	key2 = args[2] //rename for fun
	value2 = args[3]
	//verfifaction if value1 is set to 1
	agr2 := args[0:1]
	valAsbytes, err := t.read(stub, agr2)
	if err != nil {
		return nil, err
	}
	if string(valAsbytes) == "1" {
		err = stub.PutState(key1, []byte(value1)) //write the variable into the chaincode state
		if err != nil {
			return nil, err
		}
		err = stub.PutState(key2, []byte(value2)) //write the variable into the chaincode state
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("Exepect value of " + key1 + " equal 0", )
	}
	tosend := key2 + "->" + value2
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	}

	//tosend = "Change " + key1 + " to " + value1
	//err = stub.SetEvent("evtsender", []byte(tosend))
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var dat map[string]interface{}
	var err error
	fmt.Println("Runing write")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}

	if err := json.Unmarshal([]byte(args[0]), &dat); err != nil {
		panic(err)
	}
	shaByte := sha256.Sum256([]byte(args[0]))
	shaUser := base64.StdEncoding.EncodeToString(shaByte[:])
	userId := dat["userId"].(string)
	contentId := dat["contentId"].(string)
	timestampMax := dat["timestampMax"].(int64)
	timestampUser := dat["timestampUser"].(int64)
	timestamp := time.Now().Unix()
	//verify if the contract is correct
	if timestamp >= timestampMax {
		return nil, errors.New("Contract to old " + time.Unix(timestampMax, 0).String() + ". " +
			"Curent chainecode date " + time.Unix(timestamp, 0).String())
	}

	contract := &UserContract{
		userId: userId,
		contentId: contentId,
		shaUser: shaUser,
		random63:rand.Int63(),
		timestampMax:   timestampMax,
		timestampUser: timestampUser,
		timestampBrokering: timestampUser}
	contractJson, _ := json.Marshal(contract)
	shaByte = sha256.Sum256(contractJson)
	shaContract := base64.StdEncoding.EncodeToString(shaByte[:])
	err = stub.PutState(shaContract, contractJson) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	event := &EventContract{
		typeContract:"User",
		Id:contentId,
		sha:shaContract}
	tosend, _ := json.Marshal(event)
	err = stub.SetEvent("evtsender", tosend)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Read value
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	fmt.Println("Runing read")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	//fmt.Println("write retrun " + string(valAsbytes))
	return valAsbytes, nil
}
