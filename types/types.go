package types

//
type HTTPTriggerRequest struct {
	FuncName int64       `json:"funcName"`
	Payload  interface{} `json:"payload"`
}

//
type FuncTrigger struct {
	FuncName    string `json:"funcName"`
	TriggerType string `json:"triggerType"`
	ImageName   string `json:"imageName"`
	FuncPort    int32  `json:"funcPort"`
}
