package price_service

import (
	"fmt"
	"time"
)

type mockResult struct {
	price float64
	err   error
}

type mockPriceService struct {
	numCalls    int
	mockResults map[string]mockResult // what price and err to return for a particular itemCode
	callDelay   time.Duration         // how long to sleep on each call so that we can simulate calls to be expensive
}

func (m *mockPriceService) GetPriceFor(itemCode string) (float64, error) {
	m.numCalls++            // increase the number of calls
	time.Sleep(m.callDelay) // sleep to simulate expensive call

	result, ok := m.mockResults[itemCode]
	if !ok {
		panic(fmt.Errorf("bug in the tests, we didn't have a mock result for [%v]", itemCode))
	}
	return result.price, result.err
}

func (m *mockPriceService) getNumCalls() int {
	return m.numCalls
}
