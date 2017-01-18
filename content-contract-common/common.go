package content_contract_common



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
	TimestampCP        int64     `json:"TimestampCP"`
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

type EventContract struct {
	TypeContract string    `json:"typeContract"`
	Sha          string    `json:"sha"`
	Id           string    `json:"ID"`
}