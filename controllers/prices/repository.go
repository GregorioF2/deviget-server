package prices

import (
	"fmt"
	"time"

	configs "github.com/gregorioF2/deviget/configs"
	priceService "github.com/gregorioF2/deviget/priceService"

	errors "github.com/gregorioF2/deviget/lib/errors"
)

type MockService struct {
	Results map[string]float64
	Delay   time.Duration
}

func (m *MockService) GetPriceFor(itemCode string) (float64, error) {
	time.Sleep(m.Delay) // sleep to simulate expensive call

	result, ok := m.Results[itemCode]
	if !ok {
		return 0, &errors.NotFoundError{Err: fmt.Sprintf("Item code '%s' not found", itemCode)}
	}
	return result, nil
}

var service *priceService.TransparentCache

func init() {
	var mock *MockService = &MockService{
		Results: map[string]float64{
			"p1": 5,
			"p2": 7,
		},
		Delay: time.Duration(2) * time.Second,
	}

	service = priceService.NewTransparentCache(mock, time.Duration(configs.CACHE_MAX_TIME)*time.Second)
}

func GetPricesOfItem(itemCodes []string) ([]float64, error) {
	res, err := service.GetPricesFor(itemCodes...)

	if err != nil {
		return nil, err
	}
	return res, nil
}
