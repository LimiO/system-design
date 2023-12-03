package web

import (
	"fmt"
)

func ValidateBuyRequest(req *BuyRequest) error {
	if req.Count <= 0 {
		return fmt.Errorf("count can't be less than 0")
	}
	return nil
}
