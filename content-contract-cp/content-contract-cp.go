package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
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
	//"github.com/hyperledger/fabric/vendor/github.com/spf13/viper"
	//"strings"
	"github.com/hyperledger/fabric/core/comm"
	"strings"
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
	//viper.Set("peer.tls.enabled", true)
	//fmt.Println(viper.GetBool("peer.tls.enabled"))
	//fmt.Println(comm.TLSEnabled())

	var eventAddress string
	var listenToRejections bool
	var chaincodeID string
	var chaincodeIdToSend string
	var restAddress string
	var cpID string
	var percent int
	var user string

	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:7053", "address of events server")
	flag.BoolVar(&listenToRejections, "listen-to-rejections", false, "whether to listen to rejection events")
	flag.StringVar(&chaincodeID, "events-from-chaincode", "", "listen to events from given chaincode default listen all")
	flag.StringVar(&chaincodeIdToSend, "send-to-chaincode", "", "send to given chaincode default equal as -events-from-chaincode")
	flag.StringVar(&restAddress, "rest-address", "0.0.0.0:7050", "address of rest server")
	flag.StringVar(&cpID, "CP-ID", "", "id of the cp")
	flag.StringVar(&user, "user", "admin", "id of the user (default admin)")
	flag.IntVar(&percent, "percent", 100, "Percentage of chance of having the content default 100")
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
	if cpID == "" {
		data := make([]byte, 10)
		for i := range data {
			data[i] = byte(rand.Intn(256))
		}
		sha := sha256.Sum256(data)
		cpID = base64.StdEncoding.EncodeToString(sha[:])
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
				userContractForCP := getUserContract(eventContract.Sha, user, restAddress, chaincodeID)

				userReturnID := ce.ChaincodeEvent.TxID
				createCPContract(userContractForCP, userReturnID, eventContract.Sha, user, cpID, restAddress, chaincodeIdToSend, )
			}
		}
	}
}
func getUserContract(userContractSha string, user string, restAddress string, chaincodeID string) content_contract_common.UserContractForCP {
	fmt.Println("██████████████████████████Get-User-contract██████████████████████████")
	var urlToGEt string
	//login()
	if tlsbool {

		urlToGEt = "https://" + restAddress + "/chaincode"
	} else {

		urlToGEt = "http://" + restAddress + "/chaincode"
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
				Args:[]string{userContractSha}},
			SecureContext:user},

		ID:2}

	jsonpPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", urlToGEt, bytes.NewReader(jsonpPayload))
	//
	//proxyUrl, _ := url.Parse("http://localhost:8080")
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	Proxy: http.ProxyURL(proxyUrl)        }
	//client := &http.Client{Transport: tr}
	//
	//req.Header.Add("content-type", "application/json")
	//
	//res, _ := client.Do(req)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("-----------------------------RAW-User-contract----------------------------")
	fmt.Println(res)
	fmt.Println(string(body))
	fmt.Println("--------------------------User-contract-to-json---------------------------")
	response := content_contract_common.Response{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	//fmt.Println(dat)
	//result := dat["result"].(map[string]interface{})

	//rawContract := result["message"].(string)
	userContractForCP := content_contract_common.UserContractForCP{}
	json.Unmarshal([]byte(response.Result.Message), &userContractForCP)
	fmt.Println("-------------------------User-contract-json--------------------------------")
	fmt.Println(response.Result.Message)
	return userContractForCP
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
	if eventContract.TypeContract != "User" {
		fmt.Println("This is not for us")
		return false
	}

	//verify if we have a licence for this content
	if (rand.Intn(100) > percent) {
		fmt.Println("We don't have content")
		return false
	}
	fmt.Println("We have content")
	return true

}

func createCPContract(userContractForCP content_contract_common.UserContractForCP, userReturnID string, userContractID string, user string, cpID string, restAddress string, chaincodeID string) {
	fmt.Println("██████████████████████████Creat-contract██████████████████████████")
	price := 1000
	priceMax := 2000
	cPContract := content_contract_common.CPContract{
		CPId:cpID,
		TimestampMax:userContractForCP.TimestampMax,
		Random63:userContractForCP.Random63,
		ShaUser:userContractForCP.ShaUser,
		TimestampBrokering:userContractForCP.TimestampBrokering,
		TimestampUser:userContractForCP.TimestampUser,
		UserReturnID:userReturnID,
		UserContractID:userContractID,
		TimestampCP:time.Now().Unix(),
		ContentId:userContractForCP.ContentId,
		LicencingId: userContractForCP.ContentId + ".lic",
		Price:price,
		PriceMax:priceMax        }
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
	//contractOnJson := strings.Replace(string(contractJson), "\"", "\\\"", -1)
	fmt.Println("----------------------------JSON-Object----------------------------")
	fmt.Println(string(contractJson))
	//create the request
	var urltoSend string
	if tlsbool {

		urltoSend = "https://" + restAddress + "/chaincode"
	} else {

		urltoSend = "http://" + restAddress + "/chaincode"
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
				Function:"content-licencing-contract",
				Args:[]string{string(contractJson)}},
			SecureContext:user},
		ID:1}

	jsonpPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", urltoSend, bytes.NewReader(jsonpPayload))

	req.Header.Add("content-type", "application/json")

	//proxyUrl, _ := url.Parse("http://localhost:8080")
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	Proxy: http.ProxyURL(proxyUrl)        }
	//client := &http.Client{Transport: tr}
	//
	//res, _ := client.Do(req)
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

func login() {

	url := "https://d0ffb689045e4dfeb25fd8df4bafca84-vp0.us.blockchain.ibm.com:5002/registrar"

	payload := strings.NewReader("{\r\n  \"enrollId\": \"admin\",\r\n  \"enrollSecret\": \"6df1e6d3ac\"\r\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}