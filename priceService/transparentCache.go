package price_service

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	creationTime time.Time
	value        float64
}

type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]cacheEntry

	pricesMutex sync.Mutex
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]cacheEntry{},
	}
}

func (cache *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	cache.pricesMutex.Lock()
	entry, ok := cache.prices[itemCode]
	if ok {
		timePassed := time.Since(entry.creationTime)
		if timePassed > cache.maxAge {
			delete(cache.prices, itemCode)
			cache.pricesMutex.Unlock()
		} else {
			cache.pricesMutex.Unlock()
			return entry.value, nil
		}
	}

	price, err := cache.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}
	cache.pricesMutex.Lock()
	cache.prices[itemCode] = cacheEntry{value: price, creationTime: time.Now()}
	cache.pricesMutex.Unlock()

	return price, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}
	for _, itemCode := range itemCodes {
		// TODO: parallelize this, it can be optimized to not make the calls to the external service sequentially
		price, err := c.GetPriceFor(itemCode)
		if err != nil {
			return []float64{}, err
		}
		results = append(results, price)
	}
	return results, nil
}
