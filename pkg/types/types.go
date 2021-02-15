package types

// HTTPTriggerRequest model
type HTTPTriggerRequest struct {
	FuncName string      `json:"funcName"`
	Payload  interface{} `json:"payload"`
}

// ColdTriggerEvent model
type ColdTriggerEvent struct {
	FuncName  string      `json:"funcName"`
	ImageName string      `json:"imageName"`
	Payload   interface{} `json:"payload"`
}

// FuncTrigger model
type FuncTrigger struct {
	FuncName    string `json:"funcName"`
	TriggerType string `json:"triggerType"`
	ImageName   string `json:"imageName"`
	FuncPort    int32  `json:"funcPort"`
	Instances   int32  `json:"instances"`
}
