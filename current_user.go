package sdkcm

import (
	"encoding/base64"
	"encoding/json"
)

type Requester struct {
	ID uint32 `json:"id"`
	OAuthID string `json:"oauth_id"`
}

func (requester *Requester) EncodeString() string {
	str ,_ := json.Marshal(requester)
	return base64.StdEncoding.EncodeToString(str)
}

func (requester *Requester) GetSystemRole() string {
	return ""
}

func DecodeRequester(requesterStr string) (*Requester,error){
	var requester Requester
	base,_ := base64.StdEncoding.DecodeString(requesterStr)
	if err := json.Unmarshal(base,&requester); err != nil {
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
