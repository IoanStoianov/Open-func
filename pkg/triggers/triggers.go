package triggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/IoanStoianov/Open-func/pkg/types"
)

const contentType = "application/json"

var triggers map[string]*types.FuncSpecs = make(map[string]*types.FuncSpecs)

func registerNewFuncSpecs(w http.ResponseWriter, r *http.Request) {
	req, err := funcReadFuncSpecs(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	triggers[req.FuncName] = req
}

func funcReadFuncSpecs(r *http.Request) (*types.FuncSpecs, error) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Unmarshal
	var req *types.FuncSpecs
	err = json.Unmarshal(b, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// HTTPTriggerRedirect sends http request and handles response for http trigger
func HTTPTriggerRedirect(w http.ResponseWriter, r *http.Request) {
	var trigger types.HTTPTriggerRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&trigger); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	payload, err := json.Marshal(trigger.Payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	url := fmt.Sprintf("http://%s-service/triggerHttp", trigger.FuncName)

	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}
