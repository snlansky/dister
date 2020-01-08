package runner

import (
	"dister/model"
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
	ut := &model.UnitTest{
		BaseUrl:     "http://127.0.0.1:5000",
		Script:      script,
		Validator:   model.EqualValidatorType,
		ExceptValue: "pong",
	}
	resp, err := lr.Call(ut)
	assert.NoError(t, err)
	validator := ut.GetValidator()
	assert.True(t, validator.Valid(resp))
}
