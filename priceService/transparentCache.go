package price_service

import (
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
	prices             sync.Map
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             sync.Map{},
	}
}

func (cache *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	var entry *cacheEntry
	v, ok := cache.prices.Load(itemCode)
	if ok {
		entry = v.(*cacheEntry)
		timePassed := time.Since(entry.creationTime)
		if timePassed > cache.maxAge {
			cache.prices.Delete(itemCode)
		} else {
			return entry.value, nil
		}
	}

	price, err := cache.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, err
	}
	cache.prices.Store(itemCode, &cacheEntry{value: price, creationTime: time.Now()})
	return price, nil
}

func (cache *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}

	var wg sync.WaitGroup
	wg.Add(len(itemCodes))

	out := make(chan float64, len(itemCodes))
	errs := make(chan error, len(itemCodes))

	getPriceAsync := func(itemCode string) {
		defer wg.Done()
		price, err := cache.GetPriceFor(itemCode)
		if err != nil {
			errs <- err
		} else {
			out <- price
		}
	}

	for _, itemCode := range itemCodes {
		go getPriceAsync(itemCode)
	}
	wg.Wait()
	close(out)
	close(errs)
	if len(errs) > 0 {
		return nil, <-errs
	}
	for val := range out {
		results = append(results, val)
	}
	return results, nil
}
