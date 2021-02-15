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

var triggers map[string]*types.FuncTrigger = make(map[string]*types.FuncTrigger)

func registerNewFuncTrigger(w http.ResponseWriter, r *http.Request) {
	req, err := funcReadFuncTrigger(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	triggers[req.FuncName] = req
}

func funcReadFuncTrigger(r *http.Request) (*types.FuncTrigger, error) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// Unmarshal
	var req *types.FuncTrigger
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

// ReadTriggerRequest UNUSED - curl localhost:8000 -d '{"name":"Hello"}'
func ReadTriggerRequest(r *http.Request) ([]byte, error) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var req types.HTTPTriggerRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
