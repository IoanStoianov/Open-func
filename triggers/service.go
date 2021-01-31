package triggers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"open-func/k8s"
	"open-func/k8s/client"
	"open-func/types"
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
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	// Unmarshal
	var req *types.FuncTrigger
	err = json.Unmarshal(b, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

var portNum int32 = 0

func deployFunc(w http.ResponseWriter, req *http.Request) {
	client := client.OutCluster()
	dummy := types.FuncTrigger{
		FuncName:    "node-docker",
		ImageName:   "node-docker",
		TriggerType: "HttpTrigger",
		FuncPort:    portNum,
	}
	k8s.CreateDeployment(client, dummy)
	k8s.CreateService(client, dummy)
	portNum++
}

//
func HTTPTriggerRedirect(w http.ResponseWriter, r *http.Request) {
	payload, err := ReadTriggerRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.Post("http://192.168.49.2:32041/triggerHttp", contentType, bytes.NewReader(payload))
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

// curl localhost:8000 -d '{"name":"Hello"}'
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
