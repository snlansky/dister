package runner

import (
	"dister/utils"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"log"
)

type LuaRunner struct {
	state *lua.LState
}

func NewLuaRunner() IRunner {
	return &LuaRunner{
		state: lua.NewState(
			lua.Options{
				MinimizeStackMemory: true,
				IncludeGoStackTrace: true,
			}),
	}
}

func (r *LuaRunner) Call(baseUrl, script string) (string, error) {
	dict := make(map[string]string)
	dict["root"] = baseUrl

	reqDict, err := r.getOutDictByLua(script, dict)
	if err != nil {
		return "", nil
	}
	log.Println("task.id", reqDict["task.id"])
	log.Println("task.name", reqDict["task.name"])

	urlStr := reqDict["req.url"]
	log.Println("req.url", reqDict["req.url"])

	return utils.HttpGet(urlStr)
}

func (r *LuaRunner) getOutDictByLua(luaCode string, dict map[string]string) (map[string]string, error) {
	// 创建一个字典类型变量（在调用时传入)
	luaTable := r.state.NewTable()
	for k, v := range dict {
		r.state.SetTable(luaTable, lua.LString(k), lua.LString(v))
	}
	var err error

	err = r.state.DoString(luaCode)
	if err != nil {
		return nil, errors.WithMessage(err, "lua do string error")
	}

	//函数名称
	funcName := "main"
	//执行
	p := lua.P{
		Fn:      r.state.GetGlobal(funcName),
		NRet:    1,
		Protect: true,
	}
	err = r.state.CallByParam(p, luaTable)

	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("解析脚本出错(可能函数名称不能识别或脚本写错):\n----\n%s\n----\n", luaCode))
	}

	ret := r.state.Get(-1) // returned value
	r.state.Pop(1)         // remove received value
	obj := gluamapper.ToGoValue(ret, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	//类型转换:luaMap -> map
	mapTemp := obj.(map[interface{}]interface{})
	result := make(map[string]string, len(mapTemp))
	for k, v := range mapTemp {
		newKey := k.(string)
		result[newKey] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
