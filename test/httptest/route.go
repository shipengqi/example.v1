package route

import (
	"encoding/json"
	"net/http"
)

func Routes() {
	http.HandleFunc("/sendJson", SendJson)
}

func SendJson(rw http.ResponseWriter, req *http.Request)  {
	u := struct {
		Name string
	}{
		Name:"pooky",
	}
	rw.Header().Set("Content-Type","application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(u)
}
