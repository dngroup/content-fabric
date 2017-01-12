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
	var valueAnalyse string
	var valueChange string
	var restAddress string
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode")
	flag.StringVar(&valueAnalyse, "value-to-analyse", "a", "listen this value and change the value to change")
	flag.StringVar(&valueChange, "value-to-change", "b", "change this value")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server")
	flag.Parse()

	fmt.Printf("Event Address: %s\n", eventAddress)

	a := createEventClient(eventAddress, listenToRejections, chaincodeID)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
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
			if analyse(ce, valueAnalyse) {
				change(valueAnalyse, valueChange, restAddress, chaincodeID)
			}
		}
	}
}

//analyse what is the value as change
func analyse(event *pb.Event_ChaincodeEvent, valueAnalyse string) bool {
	var totest,toCompar string
	value := event.ChaincodeEvent.Payload
	//TODO: may be in json is better
	//var dat map[string]interface{}
	//if err := json.Unmarshal(value, &dat); err != nil {
	//	panic(err)
	//}
	//fmt.Println(dat)
	//num := dat["Value"]
	
	//verify if the value
	totest = string(value)
	toCompar = valueAnalyse+"->1"
	fmt.Printf("\n(%s) is equal \n(%s) ?\n", totest,toCompar)
	if totest ==  toCompar{
			fmt.Printf("True :)\n", totest,toCompar)

		return true
	} else {
			fmt.Printf("False :(\n", totest,toCompar)

		return false
	}

}

func change(valueAnalyse string, valueChange string, restAddress string, chaincodeID string) {

	url := "http://" + restAddress + "/chaincode"

	payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
		chaincodeID +
		"\" }, \"ctorMsg\": { \"function\": \"" +
		"move" +
		"\", \"args\": [ \"" +
		valueAnalyse +
		"\", \"" +
		"0"+
		"\", \"" +
		valueChange+
		"\", \"" +
		"1"+
		"\" ] }, \"secureContext\": \"test_user0\" }, \"id\": 3}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
