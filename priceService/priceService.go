package price_service

type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}
