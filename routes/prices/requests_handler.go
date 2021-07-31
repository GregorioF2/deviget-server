package prices

import (
	"fmt"
	"net/http"

	pricesController "github.com/gregorioF2/deviget/controllers/prices"

	bytesEncoding "github.com/gregorioF2/deviget/lib/utils/bytesEncoding"

	. "github.com/gregorioF2/deviget/lib/consts"
	. "github.com/gregorioF2/deviget/lib/errors"
)

func getAndValidateGetPriceQueryParams(queryParams map[string][]string) (string, error) {
	itemCodeParam, ok := queryParams["code"]
	if !ok {
		return "", &InvalidParametersError{Err: fmt.Sprintf("query param '%s' is required.", "code")}
	}
	return itemCodeParam[0], nil
}

func GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	itemCode, err := getAndValidateGetPriceQueryParams(r.URL.Query())
	if err != nil {
		var responseError *ResponseError
		switch e := err.(type) {
		case *InvalidParametersError:
			responseError = &ResponseError{Err: e.Error(), StatusCode: HttpStatusCode["ClientError"]["BadRequest"]}
		default:
			responseError = &ResponseError{Err: e.Error(), StatusCode: HttpStatusCode["ServerError"]["InternalServerError"]}
		}
		http.Error(w, responseError.Error(), responseError.StatusCode)
		return
	}

	price, err := pricesController.GetPricesOfItem(itemCode)
	if err != nil {
		var responseError *ResponseError
		switch e := err.(type) {
		case *NotFoundError:
			responseError = &ResponseError{Err: e.Error(), StatusCode: HttpStatusCode["ClientError"]["NotFound"]}
		case *InvalidParametersError:
			responseError = &ResponseError{Err: e.Error(), StatusCode: HttpStatusCode["ClientError"]["BadRequest"]}
		default:
			responseError = &ResponseError{Err: e.Error(), StatusCode: HttpStatusCode["ServerError"]["InternalServerError"]}
		}
		http.Error(w, responseError.Error(), responseError.StatusCode)
		return
	}

	response := bytesEncoding.Float64bytes(price)

	w.Write(response)

}
