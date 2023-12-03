package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"onlinestore/pkg/web"
)

func (m *MiddlewareManager) RecoverRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			err, ok := v.(error)
			if !ok {
				err = fmt.Errorf("%+v", v)
			}

			m.HandleError(w, r, err)
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (m *MiddlewareManager) HandleError(w http.ResponseWriter, _ *http.Request, err error) {
	w.WriteHeader(500)
	if err = json.NewEncoder(w).Encode(web.ResponseMetadata{
		Code:  1,
		Error: err.Error(),
	}); err != nil {
		log.Println(err, "failed to send response")
	}
}
