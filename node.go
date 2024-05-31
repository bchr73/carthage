package carthage

import (
	"context"
	"encoding/json"
	"fmt"
)

type RPCCall struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

func (r *RPCCall) JsonUnmarshal(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	return nil
}

func (c *RPCCall) Validate() (bool, error) {
	if c.Method == "" {
		return false, fmt.Errorf("missing method name")
	}
	if c.Params == nil {
		return false, fmt.Errorf("missing params")
	}
	return true, nil
}

type RPCResult struct {
	Method string      `json:"method"`
	Result interface{} `json:"result"`
}

func (r *RPCResult) JsonMarshal() ([]byte, error) {
	return json.Marshal(r)
}

type NodeService interface {
	Start(ctx context.Context, result chan *RPCResult)
	Close()
	Call(call *RPCCall) error
}
