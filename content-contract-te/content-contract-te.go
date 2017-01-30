package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
	//"encoding/json"
	//"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"math/rand"
	"time"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"github.com/dngroup/content-fabric/content-contract-common"
	"github.com/spf13/viper"
	"github.com/hyperledger/fabric/core/comm"
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

func main() {
	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	var chaincodeIdToSend string
	var restAddress string
	var teID string
	var percent int
	var pricePercent int
	var tlsbool bool
	var user string
	flag.StringVar(&user, "user", "admin", "id of the user (default admin)")
	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode default listen all")
	flag.StringVar(&chaincodeIdToSend, "send-to-chaincode", "", "send to given chaincode default equal as -events-from-chaincode")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server")
	flag.StringVar(&teID, "TE-ID", "", "id of the te")
	flag.IntVar(&percent, "percent", 100, "Percentage of chance of having the content default 100")
	flag.IntVar(&pricePercent, "percent-price", 100, "Percentage of chance of having a price lower than the maximum price default 100")
	flag.BoolVar(&tlsbool, "tls", false, "use tls")
	flag.Parse()
	if tlsbool {
		//fmt.Printf(strings.Trim(fmt.Sprintf(flag.Args()), "[]"))
		//fmt.Printf(flag.Args())
		fmt.Println("Use TLS")
		viper.SetDefault("peer.tls.enabled", true)
	}
	fmt.Println(comm.TLSEnabled())
	rand.Seed(time.Now().UnixNano())
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
			eventContract := content_contract_common.EventContract{}
			if analyse(ce, &eventContract, percent) {
				cPContractForTE := getCPContract(user, tlsbool, eventContract.Sha, restAddress, chaincodeID)
				price := verifyAndGetPrice(cPContractForTE, pricePercent)
				if price >= 0 {
					fmt.Printf("Price = %d\n", price)
					//userReturnID := ce.ChaincodeEvent.TxID
					createCPContract(cPContractForTE, eventContract.Sha, teID, user, price, tlsbool, restAddress, chaincodeIdToSend, )
				}
			}
		}
	}
}

func verifyAndGetPrice(contract content_contract_common.CPContractForTE, pricePercent int) int {
	random := rand.Intn(1000)
	fmt.Printf("random = %d\n", random)
	price := contract.Price + random * int(100 / pricePercent)
	if contract.PriceMax > price {
		fmt.Printf("################### Is not to expensive %d vs %d ##################\n", price, contract.PriceMax)
		return price
	}
	fmt.Printf("################### Is to expensive %d vs %d ##################\n", price, contract.PriceMax)
	return -1
}

func getCPContract(user string, tlsbool bool, CPContractSha string, restAddress string, chaincodeID string) content_contract_common.CPContractForTE {
	fmt.Println("██████████████████████████ Get-CP-contract ██████████████████████████")
	var url string
	if tlsbool {

		url = "https://" + restAddress + "/chaincode"
	} else {

		url = "http://" + restAddress + "/chaincode"
	}


	//payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"query\", \"params\": { \"type\": 1, \"chaincodeID\":{ \"name\":\"" +
	//	chaincodeID +
	//	"\" }, \"ctorMsg\": { \"function\":\"read\", \"args\":[\"Q" +
	//	userContractSha +
	//	"\"] } }, \"id\": 2}")
	payload := &content_contract_common.Request{
		Jsonrpc:"2.0",
		Method:"query",
		Params:content_contract_common.Params{
			Type:1,
			ChaincodeID:content_contract_common.ChaincodeID{
				Name:chaincodeID},
			CtorMsg:content_contract_common.CtorMsg{
				Function:"read",
				Args:[]string{CPContractSha}},
			SecureContext:user},

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
	response := content_contract_common.Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	//fmt.Println(dat)
	//result := dat["result"].(map[string]interface{})

	//rawContract := result["message"].(string)
	cPContractForTE := content_contract_common.CPContractForTE{}
	json.Unmarshal([]byte(response.Result.Message), &cPContractForTE)
	fmt.Println("-------------------------User-contract-json--------------------------------")
	fmt.Println(response.Result.Message)
	return cPContractForTE
}

//analyse what is the value as change
func analyse(event *pb.Event_ChaincodeEvent, eventContract *content_contract_common.EventContract, percent int) bool {
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

	//verify if we have the content
	if (rand.Intn(100) > percent) {
		fmt.Println("We don't have content")
		return false
	}
	fmt.Println("We have content")
	return true

}

func createCPContract(cPContractForTE content_contract_common.CPContractForTE, CPContractID string, teID string, user string, price int, tlsbool bool, restAddress string, chaincodeID string) {
	fmt.Println("██████████████████████████Creat-contract██████████████████████████")
	tEContract := content_contract_common.TEContract{
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
	//contractOnJson := strings.Replace(string(contractJson), "\"", "\\\"", -1)
	fmt.Println("----------------------------JSON-Object----------------------------")
	fmt.Println(string(contractJson))
	//create the request
	var url string
	if tlsbool {

		url = "https://" + restAddress + "/chaincode"
	} else {

		url = "http://" + restAddress + "/chaincode"
	}

	//payload := strings.NewReader("{ \"jsonrpc\": \"2.0\", \"method\": \"invoke\", \"params\": { \"type\": 1, \"chaincodeID\": { \"name\": \"" +
	//	chaincodeID +
	//	"\" }, \"ctorMsg\": { \"function\": \"" +
	//	"content-licencing-contract" +
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
				Function:"content-delevery-contract",
				Args:[]string{string(contractJson)}},
			SecureContext:user},

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
