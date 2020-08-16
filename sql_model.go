package sdkcm

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SQLModel struct {
	// Real id in db, we would't show it
	ID uint32 `json:"-" gorm:"id,PRIMARY_KEY"`
	// Fake id, we will public it
	FakeID    UID       `json:"id" gorm:"-"`
	Status    *int      `json:"status,omitempty"`
	CreatedAt *JSONTime `json:"created_at,omitempty"`
	UpdatedAt *JSONTime `json:"updated_at,omitempty"`
}

func (sm *SQLModel) GenUID(objType int, shardID uint32) *SQLModel {
	sm.FakeID = NewUID(sm.ID, objType, shardID)
	return sm
}

func NewSQLModelWithStatus(status int) *SQLModel {
	return &SQLModel{
		Status: &status,
	}
}

func (sm *SQLModel) ToID() *SQLModel {
	sm.ID = sm.FakeID.localID
	return sm
}

// Set time format layout. Default: 2006-01-02
func SetDateFormat(layout string) {
	dateFmt = layout
}

type JSON []byte

// This method for mapping JSON to json data type in sql
func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

// This method for scanning JSON from json data type in sql
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	if s, ok := value.([]byte); ok {
		*j = append((*j)[0:0], s...)
		return nil
	}

	return errors.New("invalid Scan Source")
}

func (j *JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return nil, errors.New("object json is nil")
	}

	return *j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("object json is nil")
	}

	*j = JSON(data)
	return nil
}

type Image struct {
	Url           string `json:"url" bson:"url"`
	OriginWidth   int    `json:"org_width" bson:"org_width"`
	OriginHeight  int    `json:"org_height" bson:"org_height"`
	OriginUrl     string `json:"org_url" bson:"org_url"`
	CloudName     string `json:"cloud_name,omitempty" bson:"cloud_name"`
	CloudId       string `json:"cloud_id,omitempty" bson:"cloud_id"`
	DominantColor string `json:"dominant_color" bson:"dominant_color"`
	RequestId     string `json:"request_id,omitempty" bson:"-"`
}

func (i *Image) HideSomeInfo() *Image {
	if i != nil {
		//i.CloudID = ""
		i.CloudId = ""
	}
	return i
}

// This method for mapping Image to json data type in sql
func (i *Image) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}

	b, err := json.Marshal(i)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Image from date data type in sql
func (i *Image) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	if err := json.Unmarshal(v, i); err != nil {
		return err
	}
	return nil
}

type Images []Image

// This method for mapping Images to json array data type in sql
func (is *Images) Value() (driver.Value, error) {
	if is == nil {
		return nil, nil
	}

	b, err := json.Marshal(is)

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// This method for scanning Images from json array type in sql
func (is *Images) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.([]byte)
	if !ok {
		return errors.New("invalid Scan Source")
	}

	var imgs []Image

	if err := json.Unmarshal(v, &imgs); err != nil {
		return err
	}

	*is = Images(imgs)
	return nil
}
