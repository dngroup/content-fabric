package main

type TEContract struct {
	TEId                   string    `json:"tEId"`
	Price                  int   `json:"price"`
	//time max after the request is deleted
	TimestampMax           int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser                string    `json:"sha_user"`
	UserContractID         string    `json:"userContractID"`
	CPContractID           string    `json:"CPContractID"`
	UserReturnID           string    `json:"userReturnID"`
	// random int
	Random63               int64     `json:"random63"`
	//use for state
	TimestampUser          int64     `json:"timestampUser"`
	TimestampUserNano      int64     `json:"timestampUserNano"`
	TimestampBrokering     int64     `json:"timestampBrokering"`
	TimestampBrokeringNano int64     `json:"timestampBrokeringNano"`
	TimestampCP            int64     `json:"timestampCP"`
	TimestampCPNano        int64     `json:"timestampCPNano"`
	TimestampLicencing     int64     `json:"timestampLicencing"`
	TimestampLicencingNano int64     `json:"timestampLicencingNano"`

	TimestampTE            int64     `json:"timestampTE"`
	TimestampTENano        int64     `json:"timestampTENano"`
}

type FinalContract struct {
	TEId                   string    `json:"tEId"`
	Price                  int   `json:"price"`
	//time max after the request is deleted
	TimestampMax           int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser                string    `json:"sha_user"`
	UserContractID         string    `json:"userContractID"`
	CPContractID           string    `json:"CPContractID"`
	UserReturnID           string    `json:"userReturnID"`
	// random int
	Random63               int64     `json:"random63"`
	//use for state
	TimestampUser          int64     `json:"timestampUser"`
	TimestampUserNano      int64     `json:"timestampUserNano"`
	TimestampBrokering     int64     `json:"timestampBrokering"`
	TimestampBrokeringNano int64     `json:"timestampBrokeringNano"`
	TimestampCP            int64     `json:"timestampCP"`
	TimestampCPNano        int64     `json:"timestampCPNano"`
	TimestampLicencing     int64     `json:"timestampLicencing"`
	TimestampLicencingNano int64     `json:"timestampLicencingNano"`

	TimestampTE            int64     `json:"timestampTE"`
	TimestampTENano        int64     `json:"timestampTENano"`
	TimestampFinal         int64     `json:"timestampFinal"`
	TimestampFinalNano     int64     `json:"timestampFinalNano"`
}

type UserContract struct {
	UserId            string    `json:"userID"`
	ContentId         string    `json:"contentID"`
	//time max after the request is deleted
	TimestampMax      int64     `json:"timestampMax"`
	//use for stat
	TimestampUser     int64     `json:"timestampUser"`
	TimestampUserNano int64     `json:"timestampUserNano"`
}

type UserContractForCP struct {
	UserId                 string    `json:"userID"`
	ContentId              string    `json:"contentID"`
	//time max after the request is deleted
	TimestampMax           int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser                string    `json:"sha_user"`
	// random int
	Random63               int64     `json:"random63"`
	//use for state
	TimestampUser          int64     `json:"timestampUser"`
	TimestampUserNano      int64     `json:"timestampUserNano"`
	TimestampBrokering     int64     `json:"timestampBrokering"`
	TimestampBrokeringNano int64     `json:"timestampBrokeringNano"`
}
type CPContract struct {
	CPId                   string    `json:"cPId"`
	ContentId              string    `json:"contentID"`
	LicencingId            string    `json:"licencingID"`
	Price                  int   `json:"price"`
	PriceMax               int   `json:"priceMax"`
	//time max after the request is deleted
	TimestampMax           int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser                string    `json:"sha_user"`
	UserContractID         string    `json:"userContractID"`
	UserReturnID           string    `json:"userReturnID"`
	// random int
	Random63               int64     `json:"random63"`
	//use for state
	TimestampUser          int64     `json:"timestampUser"`
	TimestampUserNano      int64     `json:"timestampUserNano"`
	TimestampBrokering     int64     `json:"timestampBrokering"`
	TimestampBrokeringNano int64     `json:"timestampBrokeringNano"`
	TimestampCP            int64     `json:"TimestampCP"`
	TimestampCPNano        int64     `json:"TimestampCPNano"`
}

type CPContractForTE struct {
	CPId                   string    `json:"cPId"`
	ContentId              string    `json:"contentID"`
	LicencingId            string    `json:"licencingID"`
	Price                  int   `json:"price"`
	PriceMax               int   `json:"priceMax"`
	//time max after the request is deleted
	TimestampMax           int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser                string    `json:"sha_user"`
	UserContractID         string    `json:"userContractID"`
	UserReturnID           string    `json:"userReturnID"`
	// random int
	Random63               int64     `json:"random63"`
	//use for state
	TimestampUser          int64     `json:"timestampUser"`
	TimestampUserNano      int64     `json:"timestampUserNano"`
	TimestampBrokering     int64     `json:"timestampBrokering"`
	TimestampBrokeringNano int64     `json:"timestampBrokeringNano"`
	TimestampCP            int64     `json:"timestampCP"`
	TimestampCPNano        int64     `json:"timestampCPNano"`
	TimestampLicencing     int64     `json:"timestampLicencing"`
	TimestampLicencingNano int64     `json:"timestampLicencingNano"`
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

