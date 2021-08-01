package price_service

import (
	"fmt"
	"time"

	. "github.com/gregorioF2/deviget/lib/errors"
)

type MockService struct {
	Results map[string]float64
	Delay   time.Duration
}

func (m *MockService) GetPriceFor(itemCode string) (float64, error) {
	time.Sleep(m.Delay) // sleep to simulate expensive call

	result, ok := m.Results[itemCode]
	if !ok {
		return 0, &NotFoundError{Err: fmt.Sprintf("Item code '%s' not found", itemCode)}
	}
	return result, nil
}
