package runner

import (
	"strings"
	"testing"
)

func init()  {

}
func TestLuaRunner_Call(t *testing.T) {
	script := getLuaScript()
	lr := NewLuaRunner()
	respStr := lr.Call(script)

	expectStr := "pong"
	if !strings.Contains(respStr,expectStr){
		t.Fatalf("测试失败,期望包含:【%s】,实际:【%s】",expectStr,respStr)
	}

}
