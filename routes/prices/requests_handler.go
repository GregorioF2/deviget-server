package prices

import (
	"fmt"
	"net/http"

	pricesController "github.com/gregorioF2/deviget/controllers/prices"

	"encoding/json"
	"strings"

	. "github.com/gregorioF2/deviget/lib/consts"
	. "github.com/gregorioF2/deviget/lib/errors"
)

func getAndValidateGetPriceQueryParams(queryParams map[string][]string) ([]string, error) {
	itemCodesParam, ok := queryParams["codes"]
	if !ok {
		return nil, &InvalidParametersError{Err: fmt.Sprintf("query param '%s' is required.", "codes")}
	}
	if len(itemCodesParam[0]) == 0 {
		return nil, &InvalidParametersError{Err: "no codes were sent"}
	}
	res := strings.Split(itemCodesParam[0], ",")
	return res, nil
}

func GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	itemCodes, err := getAndValidateGetPriceQueryParams(r.URL.Query())
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

	price, err := pricesController.GetPricesOfItem(itemCodes)
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

	response, err := json.Marshal(price)
	if err != nil {
		responseError := &ResponseError{Err: "failed parse response to []byte", StatusCode: HttpStatusCode["ServerError"]["InternalServerError"]}
		http.Error(w, responseError.Error(), responseError.StatusCode)
		return
	}

	w.Write(response)

}
