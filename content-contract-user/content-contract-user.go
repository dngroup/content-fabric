package main

import (
	"github.com/spf13/viper"
	"flag"
	"fmt"
	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"github.com/dngroup/content-fabric/content-contract-common"
	"bytes"
	"os"
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}
//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	if a.chaincodeID != "" {
		return []*pb.Interest{
			{EventType: pb.EventType_BLOCK},
			{EventType: pb.EventType_REJECTION},
			{EventType: pb.EventType_CHAINCODE,
				RegInfo: &pb.Interest_ChaincodeRegInfo{
					ChaincodeRegInfo: &pb.ChaincodeReg{
						ChaincodeID: a.chaincodeID,
						EventName:   ""}}}}, nil
	}
	return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {
	if o, e := msg.Event.(*pb.Event_Block); e {
		a.notfy <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_Rejection); e {
		if a.listenToRejections {
			a.rejected <- o
		}
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_ChaincodeEvent); e {
		a.cEvent <- o
		return true, nil
	}
	return false, fmt.Errorf("Receive unkown type event: %v", msg)
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	fmt.Printf("Disconnected...exiting\n")
	os.Exit(1)
}

func createEventClient(eventAddress string, listenToRejections bool, cid string) *adapter {
	var obcEHClient *consumer.EventsClient

	done := make(chan *pb.Event_Block)
	reject := make(chan *pb.Event_Rejection)
	adapter := &adapter{notfy: done, rejected: reject, listenToRejections: listenToRejections, chaincodeID: cid, cEvent: make(chan *pb.Event_ChaincodeEvent)}
	obcEHClient, _ = consumer.NewEventsClient(eventAddress, 10 * time.Second, adapter)
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

var tlsbool bool

func main() {
	fmt.Printf("Starting\n")
	var eventAddress string
	var userId string
	var contentId string
	var chaincodeID string
	var timeMax int
	var restAddress string
	var user string

	flag.StringVar(&user, "user", "admin", "id of the user (default admin)")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server (chaincode)")
	flag.StringVar(&chaincodeID, "chaincodeid", "", "chaincode Id to send the new contract")

	flag.StringVar(&userId, "userId", "user", "the userId of the user")
	flag.StringVar(&contentId, "contentId", "content", "the contentId of the content")
	flag.IntVar(&timeMax, "time-max", 10, "the timestamp max to get start the video allow by the user of the content default to 10s")
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&tlsbool, "tls", false, "use tls")
	flag.Parse()
	if tlsbool {
		viper.SetDefault("peer.tls.enabled", true)
	}

	contract := createContract(userId, contentId, timeMax)

	response := sendContract(user, contract, restAddress, chaincodeID)
	idToGetContract := response.Result.Message
	fmt.Println("██████████████████████████ Wait a result " + idToGetContract + " ██████████████████████████")


	//Create a listener to wait for the 1st contract
	a := createEventClient(eventAddress, false, chaincodeID)
	for time.Now().Unix() < contract.TimestampMax + 1 {
		fmt.Printf("-")
		select {

		case <-a.notfy:
			break
		case <-a.rejected:
			break
		case ce := <-a.cEvent:
			eventContract := content_contract_common.EventContract{}
			if analyse(ce, &eventContract, idToGetContract) {
				userContractForCP := getFinalContract(user, contract.TimestampMax, restAddress, chaincodeID, idToGetContract)
				fmt.Println(userContractForCP)
				return
			}
		case <-time.After(time.Second * 1):
			break

		}

	}
	// if no event receive try to get the contract
	userContractForCP := getFinalContract(user, contract.TimestampMax, restAddress, chaincodeID, idToGetContract)
	fmt.Println(userContractForCP)
	return
}
func getFinalContract(user string, timestampMax int64, restAddress string, chaincodeID string, idToGetContract string) string {
	payloadQuery := &content_contract_common.Request{
		Jsonrpc:"2.0",
		Method:"query",
		Params:content_contract_common.Params{
			Type:1,
			ChaincodeID:content_contract_common.ChaincodeID{
				Name:chaincodeID},
			CtorMsg:content_contract_common.CtorMsg{
				Function:"read",
				Args:[]string{idToGetContract}},
			SecureContext:user},

		ID:2}

	jsonpPayload, _ := json.Marshal(payloadQuery)

	var url string
	if tlsbool {

		url = "https://" + restAddress + "/chaincode"
	} else {

		url = "http://" + restAddress + "/chaincode"
	}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("content-type", "application/json")
	for time.Now().Unix() < timestampMax + 2 {
		time.Sleep(100 * time.Millisecond)
		req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonpPayload))
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		//fmt.Println("-----------------------------RAW-contract----------------------------")
		//fmt.Println(res)
		//fmt.Println(string(body))
		//fmt.Println("--------------------------contract-to-json---------------------------")
		response := content_contract_common.Response{}
		if err := json.Unmarshal(body, &response); err != nil {
			panic(err)
		}
		//fmt.Println(dat)
		//result := dat["result"].(map[string]interface{})

		//rawContract := result["message"].(string)

		if (response.Result.Message != "") {
			fmt.Println()
			fmt.Println("-------------------------te-contract-json--------------------------------")
			fmt.Println(response.Result.Message)
			finalContract := content_contract_common.FinalContract{}
			json.Unmarshal([]byte(response.Result.Message), &finalContract)
			finalContract.TimestampFinal = time.Now().Unix()
			fmt.Println("-------------------------final-contract-json--------------------------------")
			finalContractjson, _ := json.Marshal(finalContract)
			//fmt.Println(finalContractjson)
			return string(finalContractjson)

		} else {
			fmt.Printf("#")
		}

	}
	return "no contract found ..."

}

func analyse(event *pb.Event_ChaincodeEvent, eventContract *content_contract_common.EventContract, idToGetContract string) bool {
	//fmt.Println("██████████████████████████Analyse--contract██████████████████████████")
	data := event.ChaincodeEvent.Payload
	err := json.Unmarshal([]byte(data), &eventContract)
	if err != nil {
		//fmt.Println("This is not correct ")
		return false
	}

	//verify if is a FINAL contract
	if eventContract.TypeContract != "FINAL" {
		//fmt.Println("This is not for us")
		return false
	}

	//verify if the idToGetContract is correct
	if (idToGetContract != eventContract.Id) {
		//fmt.Println("This is not for us")
		return false
	}
	fmt.Println("this is for us start to get contract")
	return true
}

func sendContract(user string, contract content_contract_common.UserContract, restAddress string, chaincodeID string) content_contract_common.Response {
	// use this format to enable the json on the payload json
	contractJson, err := json.Marshal(contract)
	if (err != nil) {
		fmt.Println(err.Error())
	}
	fmt.Println("----------------------------JSON-Object----------------------------")
	//fmt.Println(string(contractOnJson))
	//create the request
	urltosend := ""
	if tlsbool {

		urltosend = "https://" + restAddress + "/chaincode"
	} else {

		urltosend = "http://" + restAddress + "/chaincode"
	}




	//payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
	//	chaincodeID +
	//	"\" }, \"ctorMsg\": { \"function\": \"" +
	//	"content-brokering-contract" +
	//	"\", \"args\": [ \"" +
	//	contractOnJson +
	//	"\" ] } }, \"id\": 1}")

	payload := &content_contract_common.Request{
		Jsonrpc:"2.0",
		Method:"invoke",
		Params:content_contract_common.Params{
			Type:1,
			ChaincodeID:content_contract_common.ChaincodeID{
				Name:chaincodeID},
			CtorMsg:content_contract_common.CtorMsg{
				Function:"content-brokering-contract",
				Args:[]string{string(contractJson)}},
			SecureContext:user},

		ID:2}

	jsonpPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", urltosend, bytes.NewReader(jsonpPayload))
	//proxyUrl, err := url.Parse("http://localhost:8080")
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	Proxy: http.ProxyURL(proxyUrl)        }
	//client := &http.Client{Transport: tr}
	//req, _ = http.NewRequest("POST", url, bytes.NewReader(jsonpPayload))
	//res, err := client.Do(req)
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)




	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("--------------------------------SEND--------------------------------")
	fmt.Println(string(jsonpPayload))
	fmt.Println("to: ", urltosend)
	fmt.Println("-------------------------------RECIVE-------------------------------")
	fmt.Println(res)
	fmt.Println(string(body))

	response := content_contract_common.Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return response

}
func createContract(userId string, contentId string, timeMax int) content_contract_common.UserContract {
	fmt.Printf("Create a new contract for %s\n", userId)

	//create the new contract
	contract := content_contract_common.UserContract{
		UserId:userId,
		ContentId:contentId,
		TimestampMax:time.Now().Add(time.Duration(timeMax) * time.Second).Unix(),
		TimestampUser:time.Now().Unix()}
	fmt.Println("-----------------------------Raw-Object----------------------------")
	fmt.Println(contract)
	//convert to json

	return contract

}
