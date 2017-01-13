package content_contract_common


type UserContract struct {
	userId             string    `json:"userID"`
	contentId          string    `json:"contentID"`
	//time max after the request is deleted
	timestampMax       int64     `json:"timestampMax"`
	//sha of user massage
	sha_user           string    `json:"sha_user"`
	// random int
	random63           int64     `json:"random63"`
	//use for state
	timestampUser      int64     `json:"timestampUser"`
	timestampBrokering int64     `json:"timestampBrokering"`
}
