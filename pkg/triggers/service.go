package triggers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/IoanStoianov/Open-func/pkg/k8s"
	"github.com/IoanStoianov/Open-func/pkg/k8s/client"
	"github.com/IoanStoianov/Open-func/pkg/types"
)

const contentType = "application/json"

var triggers map[string]*types.FuncTrigger = make(map[string]*types.FuncTrigger)

//
func RegisterNewFuncTrigger(w http.ResponseWriter, r *http.Request) {
	req, err := FuncReadFuncTrigger(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	triggers[req.FuncName] = req
}

//
func FuncReadFuncTrigger(r *http.Request) (*types.FuncTrigger, error) {
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

//
func DeployFunc(w http.ResponseWriter, req *http.Request) {
	var payload types.FuncTrigger

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	defer req.Body.Close()

	client := client.InCluster()

	k8s.CreateDeployment(client, payload)
	k8s.CreateService(client, payload)

	w.WriteHeader(204)
}

//
func HTTPTriggerRedirect(w http.ResponseWriter, r *http.Request) {
	payload, err := ReadTriggerRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.Post("http://nod4o-service/triggerHttp", contentType, bytes.NewReader(payload))
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

// ReadTriggerRequest - curl localhost:8000 -d '{"name":"Hello"}'
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
