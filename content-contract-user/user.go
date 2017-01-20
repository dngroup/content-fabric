package main

import (
	"flag"
	"fmt"
	"strings"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"github.com/dngroup/content-fabric/content-contract-common"
	"bytes"
)

func main() {
	fmt.Printf("Starting\n")
	var userId string
	var contentId string
	var chaincodeID string
	var timeMax int
	var restAddress string
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server (chaincode)")
	flag.StringVar(&chaincodeID, "chaincodeid", "", "chaincode Id to send the new contract")

	flag.StringVar(&userId, "userId", "user", "the userId of the user")
	flag.StringVar(&contentId, "contentId", "content", "the contentId of the content")
	flag.IntVar(&timeMax, "time-max", 10, "the timestamp max to get start the video allow by the user of the content default to 10s")
	flag.Parse()

	fmt.Printf("Create a new contract for %s\n", userId)

	//creat the new contract
	contract := content_contract_common.UserContract{
		userId,
		contentId,
		time.Now().Add(time.Duration(timeMax) * time.Second).Unix(),
		time.Now().Unix()}
	fmt.Println("-----------------------------Raw-Object----------------------------")
	fmt.Println(contract)
	//convert to json
	contractJson, err := json.Marshal(contract)
	if (err != nil) {
		return
	}
	fmt.Println("----------------------------JSON-Object----------------------------")
	fmt.Println("len:", len(contractJson))
	// use this format to enable the json on the payload json
	contractOnJson := strings.Replace(string(contractJson), "\"", "\\\"", -1)
	fmt.Println("----------------------------JSON-Object----------------------------")
	fmt.Println(string(contractOnJson))
	//create the request
	url := "http://" + restAddress + "/chaincode"

	payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
		chaincodeID +
		"\" }, \"ctorMsg\": { \"function\": \"" +
		"content-brokering-contract" +
		"\", \"args\": [ \"" +
		contractOnJson +
		"\" ] } }, \"id\": 1}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("--------------------------------SEND--------------------------------")
	fmt.Println(payload)
	fmt.Println("-------------------------------RECIVE-------------------------------")
	fmt.Println(res)
	fmt.Println(string(body))

	response := content_contract_common.Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	fmt.Println("██████████████████████████ Wait a result " + response.Result.Message + " ██████████████████████████")

	payloadQuery := &content_contract_common.Request{
		Jsonrpc:"2.0",
		Method:"query",
		Params:content_contract_common.Params{
			Type:1,
			ChaincodeID:content_contract_common.ChaincodeID{
				Name:chaincodeID},
			CtorMsg:content_contract_common.CtorMsg{
				Function:"read",
				Args:[]string{response.Result.Message}}},
		ID:2}

	jsonpPayload, _ := json.Marshal(payloadQuery)


	req.Header.Add("content-type", "application/json")

	for time.Now().Unix() < contract.TimestampMax {
		time.Sleep(1000 * time.Millisecond)
		req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonpPayload))
		res, _ = http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)
		//fmt.Println("-----------------------------RAW-contract----------------------------")
		//fmt.Println(res)
		//fmt.Println(string(body))
		//fmt.Println("--------------------------contract-to-json---------------------------")
		response = content_contract_common.Response{}
		if err := json.Unmarshal(body, &response); err != nil {
			panic(err)
		}
		//fmt.Println(dat)
		//result := dat["result"].(map[string]interface{})

		//rawContract := result["message"].(string)
		teContractForTE := content_contract_common.TEContract{}
		json.Unmarshal([]byte(response.Result.Message), &teContractForTE)
		fmt.Println("-------------------------te-contract-json--------------------------------")
		fmt.Println(response.Result.Message)

	}
}
