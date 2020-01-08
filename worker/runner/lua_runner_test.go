package runner

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func init() {

}
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
	resp := lr.Call("http://127.0.0.1:5000", script)

	assert.True(t, strings.Contains(resp, "pong"))
}
