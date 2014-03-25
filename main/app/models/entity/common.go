package entity

import (
	"time"

	gtime "github.com/itang/gotang/time"
	//	"encoding/json"
)

type regionData map[string]string

func (e regionData) GetById(id string) string {
	v, ok := e[id]
	if !ok {
		return ""
	}
	return v
}

type JsonTime time.Time

func JsonTimeNow() JsonTime {
	return JsonTime(time.Now())
}

func (j JsonTime) format() string {
	t := time.Time(j)
	if t.IsZero() {
		return ""
	}

	return t.Format(gtime.ChinaDefaultDateTime)
}

func (j JsonTime) MarshalText() ([]byte, error) {
	return []byte(j.format()), nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

type ParamsForNewOrder struct {
	CartId    int64
	ProductId int64
	PrefPrice float64
	Quantity  int
}
