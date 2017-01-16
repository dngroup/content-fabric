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
)

type adapter struct {
	notfy              chan *pb.Event_Block
	rejected           chan *pb.Event_Rejection
	cEvent             chan *pb.Event_ChaincodeEvent
	listenToRejections bool
	chaincodeID        string
}

type UserContractForCP struct {
	UserId             string    `json:"userID"`
	ContentId          string    `json:"contentID"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
}
type CPContract struct {
	UserReturnID       string    `json:"userReturnID"`
	ContentId          string    `json:"contentID"`
	LicencingId        string    `json:"licencingID"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"timestampLicencing"`
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
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode default listen all")
	flag.StringVar(&chaincodeIdToSend, "send-to-chaincode", "", "send to given chaincode default equal as -events-from-chaincode")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server")
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
			userContractForCP := UserContractForCP{}
			if analyse(ce, &userContractForCP) {
				userReturnID := ce.ChaincodeEvent.TxID
				createCPContract(userContractForCP, userReturnID, restAddress, chaincodeIdToSend)
			}
		}
	}
}

//analyse what is the value as change
func analyse(event *pb.Event_ChaincodeEvent, userContractForCP *UserContractForCP) bool {

	data := event.ChaincodeEvent.Payload
	err := json.Unmarshal([]byte(data), &userContractForCP)
	if err != nil {
		fmt.Println("This is not for us")
		return false
	}

	//verify if we have a licence for this content
	//TODO: edit this value to have a real random
	if (rand.Intn(10) < 0) {
		fmt.Println("We don't have content")
		return false
	}

	return true

}

func createCPContract(userContractForCP UserContractForCP, userReturnID string, restAddress string, chaincodeID string) {

	cPContract := CPContract{
		TimestampMax:userContractForCP.TimestampMax,
		Random63:userContractForCP.Random63,
		ShaUser:userContractForCP.ShaUser,
		TimestampBrokering:userContractForCP.TimestampBrokering,
		TimestampUser:userContractForCP.TimestampUser,
		UserReturnID:userReturnID,

		ContentId:userContractForCP.ContentId,
		LicencingId: userContractForCP.ContentId + ".lic",

		TimestampCP:time.Now().Unix(),
	}
	fmt.Println("-----------------------------Raw-Object----------------------------")
	fmt.Println(cPContract)
	//convert to json
	contractJson, err := json.Marshal(cPContract)
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
		"content-licencing-contract" +
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
}
