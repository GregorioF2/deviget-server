package prices

import (
	"time"

	configs "github.com/gregorioF2/deviget/configs"
	priceService "github.com/gregorioF2/deviget/priceService"
)

var mock *priceService.MockService = &priceService.MockService{
	Results: map[string]float64{
		"p1": 5,
		"p2": 7,
	},
	Delay: time.Duration(2) * time.Second,
}
var service *priceService.TransparentCache = priceService.NewTransparentCache(mock, time.Duration(configs.CACHE_MAX_TIME)*time.Second)

func GetPricesOfItem(itemCodes []string) ([]float64, error) {
	res, err := service.GetPricesFor(itemCodes...)

	if err != nil {
		return nil, err
	}
	return res, nil
}
