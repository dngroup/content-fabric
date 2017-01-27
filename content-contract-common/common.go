package content_contract_common

type TEContract struct {
	TEId               string    `json:"tEId"`
	Price              int   `json:"price"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID"`
	CPContractID       string    `json:"CPContractID"`
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

type FinalContract struct {
	TEId               string    `json:"tEId"`
	Price              int   `json:"price"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID"`
	CPContractID       string    `json:"CPContractID"`
	UserReturnID       string    `json:"userReturnID"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"timestampCP"`
	TimestampLicencing int64     `json:"timestampLicencing"`

	TimestampTE        int64     `json:"timestampTE"`
	TimestampFinal     int64     `json:"timestampFinal"`
}

type UserContract struct {
	UserId        string    `json:"userID"`
	ContentId     string    `json:"contentID"`
	//time max after the request is deleted
	TimestampMax  int64     `json:"timestampMax"`
	//use for stat
	TimestampUser int64     `json:"timestampUser"`
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
	CPId               string    `json:"cPId"`
	ContentId          string    `json:"contentID"`
	LicencingId        string    `json:"licencingID"`
	Price              int   `json:"price"`
	PriceMax           int   `json:"priceMax"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID"`
	UserReturnID       string    `json:"userReturnID"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"TimestampCP"`
}

type CPContractForTE struct {
	CPId               string    `json:"cPId"`
	ContentId          string    `json:"contentID"`
	LicencingId        string    `json:"licencingID"`
	Price              int   `json:"price"`
	PriceMax           int   `json:"priceMax"`
	//time max after the request is deleted
	TimestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	ShaUser            string    `json:"sha_user"`
	UserContractID     string    `json:"userContractID"`
	UserReturnID       string    `json:"userReturnID"`
	// random int
	Random63           int64     `json:"random63"`
	//use for state
	TimestampUser      int64     `json:"timestampUser"`
	TimestampBrokering int64     `json:"timestampBrokering"`
	TimestampCP        int64     `json:"timestampCP"`
	TimestampLicencing int64     `json:"timestampLicencing"`
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

