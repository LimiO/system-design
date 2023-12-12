package web

import (
	"fmt"

	"onlinestore/services/purchaseservice/types"
)

func ValidateBuyRequest(req *types.BuyRequest) error {
	if req.Count <= 0 {
		return fmt.Errorf("count can't be less than 0")
	}
	return nil
}
