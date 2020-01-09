package runner

import (
	"dister/protos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLuaRunner_Call(t *testing.T) {
	script := `
function main(dict)
  rootUrl = dict["root"]
  outDict = {}

  outDict["task.id"] = "9"
  outDict["task.name"] = "simple ping"

  outDict["req.url"] = string.format("%s/ping",rootUrl)
  return outDict
end`
	lr := NewLuaRunner()
	var vs []*protos.Validator
	vs = append(vs, &protos.Validator{
		Name:  "body",
		Vt:    protos.Validator_EQ,
		Value: "pong",
	})
	task := &protos.Task{
		Url:    "http://127.0.0.1:5000",
		Path:   "/ping",
		Method: protos.Task_GET,
		Body:   nil,
		Script: script,
		Vs:     vs,
	}
	resp, err := lr.Call(task)
	assert.NoError(t, err)

	for _, v := range task.Vs {
		if v.Name == "body" {
			if v.Vt == protos.Validator_EQ {
				assert.Equal(t, v.Value, resp)
			}
		}
	}
}
