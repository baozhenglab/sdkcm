package sdkcm

import "encoding/json"

type Requester struct {
	ID uint32 `json:"id"`
	OAuthID string `json:"oauth_id"`
}

func (requester *Requester) EncodeString() string {
	str ,_ := json.Marshal(requester)
	return string(str)
}

func (requester *Requester) GetSystemRole() string {
	return ""
}

func DecodeRequester(requesterStr string) (*Requester,error){
	var requester Requester
	if err := json.Unmarshal([]byte(requesterStr),&requester); err != nil {
		return nil,err
	}
	return &requester,nil
}

type User interface {
	UserID() uint32
	GetSystemRole() string
	GetUser() interface{}
}

type OAuth interface {
	OAuthID() string
}

func CurrentUser(t OAuth, u User) *Requester {
	return &Requester{u.UserID(), t.OAuthID()}
}
