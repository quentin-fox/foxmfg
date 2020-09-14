package chi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encoder(w http.ResponseWriter) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}

func ok(w http.ResponseWriter, result interface{}) {
	r := Response{
		Status: http.StatusOK,
		Result: result,
	}

	if err := encoder(w).Encode(r); err != nil {
		fmt.Println(err)
	}
}

func fail(w http.ResponseWriter, err error, status int) {
	r := Response{
		Status: status,
	}

	if msg := err.Error(); msg != "" {
		r.Message = msg
	}

	if err := encoder(w).Encode(r); err != nil {
		fmt.Println(err)
	}
}

func decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
