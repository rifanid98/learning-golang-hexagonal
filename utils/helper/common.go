package helper

import (
	"crypto/rand"
	"encoding/json"
)

type Helper interface {
	StructToMap(obj interface{}) (newMap map[string]interface{}, err error)
	GenerateRandomNumber() uint
}

type helper struct{}

func New() *helper {
	return &helper{}
}

func (h *helper) StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

func (h *helper) GenerateRandomNumber() uint {
	p, _ := rand.Prime(rand.Reader, 32)
	return uint(p.Uint64())
}
