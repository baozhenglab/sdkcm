package sdkcm

import "encoding/json"

func StructToJson(v interface{}) map[string]interface{} {
	var dataRes map[string]interface{}
	inrec, _ := json.Marshal(v)
	json.Unmarshal(inrec, &dataRes)
	return dataRes
}
