package lib

import (
	"encoding/json"
	"net/http"
)

func ResponseWritter(w http.ResponseWriter, data interface{}, status int) {
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
