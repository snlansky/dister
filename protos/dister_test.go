package protos

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskData_String(t *testing.T) {
	task := &TaskData{
		Threads:    0,
		Delay:      0,
		CreateTime: 0,
		CountLimit: 0,
		TimeLimit:  0,
		Task: &Task{
			Url:    "http://127.0.0.1:8080",
			Path:   "/api/ping",
			Method: Task_GET,
			Body:   nil,
			Script: "",
			Vs: []*Validator{
				{
					Name:  "code",
					Vt:    Validator_EQ,
					Value: "200",
				},
			},
		},
		Result: nil,
	}

	bytes, err := json.Marshal(task)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}
