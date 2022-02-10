package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
)

type HelperMock struct {
	mock.Mock
}

func (h *HelperMock) StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	args := h.Called(obj)
	if args.Get(0) == nil {
		return nil, errors.New("failed to struct map")
	}
	return args.Get(0).(map[string]interface{}), nil
}

func (h *HelperMock) GenerateRandomNumber() uint {
	args := h.Called()
	if args.Get(0) == nil {
		return 0
	}
	return args.Get(0).(uint)
}
