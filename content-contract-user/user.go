package main

import (
	"flag"
	"fmt"
	"strings"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
)

type UserContract struct {
	userId        string    `json:"userID"`
	contentId     string    `json:"contentID"`
	//time max after the request is deleted
	timestampMax  int64     `json:"timestampMax"`
	//use for stat
	timestampUser int64     `json:"timestampUser"`
}

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
	flag.StringVar(&timeMax, "time-max", 10, "the timestamp max to get start the video allow by the user of the content default to 10s")
	flag.Parse()

	fmt.Printf("Create a new contract for %s\n", userId)

	//creat the new contract
	contract := &UserContract{
		userId:userId,
		contentId:contentId,
		timestampMax: time.Now().Add(time.Duration(timeMax) * time.Second).Unix(),
		timestampUser: time.Now().Unix()}
	//convert to json
	contractJson, _ := json.Marshal(contract)
	// use this format to enable the json on the payload json
	contractOnJson := strings.Replace(string(contractJson), "\"", "\\\"", -1)

	//create the request
	url := "http://" + restAddress + "/chaincode"

	payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
		chaincodeID +
		"\" }, \"ctorMsg\": { \"function\": \"" +
		"content-brockering-contract" +
		"\", \"args\": [ \"" +
		contractOnJson +
		"\" ] } }, \"id\": 1}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))

}
