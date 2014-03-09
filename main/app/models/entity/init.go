package entity

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itang/gotang"
)

var rd regionData = make(map[string]string)

func init() {
	file, err := os.Open("public/libs/city/city.json")
	gotang.AssertNoError(err)
	var data map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	gotang.AssertNoError(err)

	for pid, v := range data {
		pData := v.(map[string]interface{})
		rd[pid] = pData["name"].(string)
		for cid, v := range pData["data"].(map[string]interface{}) {
			cData := v.(map[string]interface{})
			rd[pid+cid] = cData["name"].(string)
			if _, ok := cData["data"].(map[string]interface{}); ok {
				for did, v := range cData["data"].(map[string]interface{}) {
					dData := v.(map[string]interface{})
					rd[pid+cid+did] = dData["name"].(string)
				}
			} else {
				for i, v := range cData["data"].([]interface{}) {
					fmt.Printf("dd\t\t%s:%s\n", i, v)
				}
			}
		}
	}
}
