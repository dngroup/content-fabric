package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
	//"encoding/json"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"time"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}

type CPContractForTE struct {
	CPId               string    `json:"cPId"`
	ContentId          string    `json:"contentID"`
	LicencingId        string    `json:"licencingID"`

	Price              float64   `json:"price"`
	PriceMax           float64   `json:"priceMax"`

	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID`
	UserReturnID       string    `json:"userReturnID"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"timestampCP"`
	TimestampLicencing int64     `json:"timestampLicencing"`
}
type TEContract struct {
	TEId               string    `json:"tEId"`
	Price              float64   `json:"price"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID`
	CPContractID       string    `json:"CPContractID`
	UserReturnID       string    `json:"userReturnID"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"timestampCP"`
	TimestampLicencing int64     `json:"timestampLicencing"`

	TimestampTE        int64     `json:"timestampTE"`
}

type EventContract struct {
	TypeContract string    `json:"typeContract"`
	Sha          string    `json:"sha"`
	Id           string    `json:"ID"`
}
type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		} `json:"result"`
	ID      int `json:"id"`
}
type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int `json:"id"`
}

type Params struct {
	Type          int `json:"type"`
	ChaincodeID   ChaincodeID `json:"chaincodeID"`
	CtorMsg       CtorMsg `json:"ctorMsg"`
	SecureContext string `json:"secureContext"`
}
type ChaincodeID   struct {
	Name string `json:"name"`
}
type CtorMsg       struct {
	Function string `json:"function"`
	Args     []string `json:"args"`
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
	obcEHClient, _ = consumer.NewEventsClient(eventAddress, 5, adapter)
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

func main() {
	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	var chaincodeIdToSend string
	var restAddress string
	var teID string
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode default listen all")
	flag.StringVar(&chaincodeIdToSend, "send-to-chaincode", "", "send to given chaincode default equal as -events-from-chaincode")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server")
	flag.StringVar(&teID, "TE-ID", "", "id of the te")
	flag.Parse()

	fmt.Printf("Event Address: %s\n", eventAddress)

	a := createEventClient(eventAddress, listenToRejections, chaincodeID)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
	}

	//set default value to the same as events chaincode
	if chaincodeIdToSend == "" {
		if chaincodeID == "" {
			fmt.Printf("No chaincode set\n")
			return
		}
		chaincodeIdToSend = chaincodeID
	}
	//if the CP have not id set random
	if teID == "" {
		data := make([]byte, 10)
		for i := range data {
			data[i] = byte(rand.Intn(256))
		}
		sha := sha256.Sum256(data)
		teID = base64.StdEncoding.EncodeToString(sha[:])
	}

	for {
		select {
		case b := <-a.notfy:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received block\n")
			fmt.Printf("--------------\n")
			for _, r := range b.Block.Transactions {
				fmt.Printf("Transaction:\n\t[%v]\n", r)
			}
		case r := <-a.rejected:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received rejected transaction\n")
			fmt.Printf("--------------\n")
			fmt.Printf("Transaction error:\n%s\t%s\n", r.Rejection.Tx.Txid, r.Rejection.ErrorMsg)
		case ce := <-a.cEvent:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Received chaincode event\n")
			fmt.Printf("------------------------\n")
			fmt.Printf("Chaincode Event:%v\n", ce)
			eventContract := EventContract{}
			if analyse(ce, &eventContract) {
				cPContractForTE := getCPContract(eventContract.Sha, restAddress, chaincodeID)
				var price float64
				isNotToExpensive(cPContractForTE, &price)
				{
					//userReturnID := ce.ChaincodeEvent.TxID
					createCPContract(cPContractForTE, eventContract.Sha, teID, price, restAddress, chaincodeIdToSend, )
				}
			}
		}
	}
}

func isNotToExpensive(contract CPContractForTE, price float64) bool {
	price = (contract.PriceMax + float64(rand.Int31n(500) / 100))
	if contract.PriceMax > price {
		return true
	}
	fmt.Println("######################### Is to expensive ###########################")
	return false
}

func getCPContract(CPContractSha string, restAddress string, chaincodeID string) CPContractForTE {
	fmt.Println("██████████████████████████Get-User-contract██████████████████████████")
	url := "http://" + restAddress + "/chaincode"


	//payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"query\", \"params\": { \"type\": 1, \"chaincodeID\":{ \"name\":\"" +
	//	chaincodeID +
	//	"\" }, \"ctorMsg\": { \"function\":\"read\", \"args\":[\"Q" +
	//	userContractSha +
	//	"\"] } }, \"id\": 2}")
	payload := &Request{
		Jsonrpc:"2.0",
		Method:"query",
		Params:Params{
			Type:1,
			ChaincodeID:ChaincodeID{
				Name:chaincodeID},
			CtorMsg:CtorMsg{
				Function:"read",
				Args:[]string{CPContractSha}}},
		ID:2}

	jsonpPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonpPayload))

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("-----------------------------RAW-CP-contract----------------------------")
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println("--------------------------CP-contract-to-json---------------------------")
	response := Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	//fmt.Println(dat)
	//result := dat["result"].(map[string]interface{})

	//rawContract := result["message"].(string)
	cPContractForTE := CPContractForTE{}
	json.Unmarshal([]byte(response.Result.Message), &cPContractForTE)
	fmt.Println("-------------------------User-contract-json--------------------------------")
	fmt.Println(response.Result.Message)
	return cPContractForTE
}

//analyse what is the value as change
func analyse(event *pb.Event_ChaincodeEvent, eventContract *EventContract) bool {
	fmt.Println("██████████████████████████Analyse--contract██████████████████████████")
	data := event.ChaincodeEvent.Payload
	err := json.Unmarshal([]byte(data), &eventContract)
	if err != nil {
		fmt.Println("This is not correct ")
		return false
	}

	//verify if is a user contract
	if eventContract.TypeContract != "Licencing" {
		fmt.Println("This is not for us")
		return false
	}

	//verify if we have a licence for this content
	//TODO: edit this value to have a real random
	if (rand.Intn(10) < 0) {
		fmt.Println("We don't have content")
		return false
	}
	fmt.Println("We have content")
	return true

}

func createCPContract(cPContractForTE CPContractForTE, CPContractID string, teID string, price float64, restAddress string, chaincodeID string) {
	fmt.Println("██████████████████████████Creat-contract██████████████████████████")
	tEContract := TEContract{
		TEId:teID,
		TimestampMax:cPContractForTE.TimestampMax,
		Random63:cPContractForTE.Random63,
		ShaUser:cPContractForTE.ShaUser,

		UserReturnID:cPContractForTE.UserReturnID,
		UserContractID:cPContractForTE.UserContractID,

		//ContentId:cPContractForTE.ContentId,
		//LicencingId: cPContractForTE.ContentId + ".lic",
		Price:price,
		CPContractID:CPContractID,
		TimestampBrokering:cPContractForTE.TimestampBrokering,
		TimestampUser:cPContractForTE.TimestampUser,
		TimestampCP:cPContractForTE.TimestampCP,
		TimestampLicencing:cPContractForTE.TimestampLicencing,
		TimestampTE:time.Now().Unix()        }
	fmt.Println("-----------------------------Raw-Object----------------------------")
	fmt.Println(tEContract)
	//convert to json
	contractJson, err := json.Marshal(tEContract)
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

	//payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
	//	chaincodeID +
	//	"\" }, \"ctorMsg\": { \"function\": \"" +
	//	"content-licencing-contract" +
	//	"\", \"args\": [ \"" +
	//	contractOnJson +
	//	"\" ] } }, \"id\": 1}")
	payload := &Request{
		Jsonrpc:"2.0",
		Method:"invoke",
		Params:Params{
			Type:1,
			ChaincodeID:ChaincodeID{
				Name:chaincodeID},
			CtorMsg:CtorMsg{
				Function:"content-licencing-contract",
				Args:[]string{contractOnJson}}},
		ID:1}

	jsonpPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonpPayload))

	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("--------------------------------SEND--------------------------------")
	fmt.Println(string(jsonpPayload))
	fmt.Println("-------------------------------RECIVE-------------------------------")
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println("██████████████████████████Contract-send██████████████████████████")
}