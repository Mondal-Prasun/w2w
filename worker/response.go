package worker

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
}

func (r *Response) success(statusCode int, payload any) error {

	marshaledData, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	r.w.WriteHeader(statusCode)

	r.w.Write(marshaledData)

	return nil

}

func (r *Response) error(statusCode int, errorMsg string) {
	http.Error(r.w, errorMsg, statusCode)
}
