package runner

import (
	"fmt"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"log"
	"net/http"
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

func (r *LuaRunner) Call(baseUrl , script string) string {
	dict := make(map[string]string)
	dict["root"] = baseUrl

	reqDict := getOutDictByLua(r.state, script, dict)
	log.Println("task.id", reqDict["task.id"])
	log.Println("task.name", reqDict["task.name"])

	urlStr := reqDict["req.url"]
	log.Println("req.url", reqDict["req.url"])

	respStr := httpGet(urlStr)

	return respStr
}

func httpGet(urlStr string) string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func getOutDictByLua(luaState *lua.LState, luaCode string, dict map[string]string) (result map[string]string) {
	// 创建一个字典类型变量（在调用时传入)
	luaTable := luaState.NewTable()
	for k, v := range dict {
		luaState.SetTable(luaTable, lua.LString(k), lua.LString(v))
	}
	var err error
	err = luaState.DoString(luaCode)
	if err != nil {
		panic(fmt.Sprintf("DoString err:%s", err.Error()))
	}

	//函数名称
	funcName := "main"
	//执行
	err = luaState.CallByParam(lua.P{
		Fn:      luaState.GetGlobal(funcName),
		NRet:    1,
		Protect: true,},
		//下面是对应main函数传入参数的部分
		luaTable)

	if err != nil {
		panic(fmt.Sprintf("%s:解析脚本出错(可能函数名称不能识别或脚本写错):%s:\n----\n%s\n----\n", "", err.Error(), luaCode))
	}

	ret := luaState.Get(-1) // returned value
	luaState.Pop(1)         // remove received value
	obj := gluamapper.ToGoValue(ret, gluamapper.Option{NameFunc: func(s string) string {
		return s
	}})
	//类型转换:luaMap -> map
	mapTemp := obj.(map[interface{}]interface{})
	result = make(map[string]string)
	for k, v := range mapTemp {
		newKey := k.(string)
		result[newKey] = fmt.Sprintf("%v", v)
	}
	return result
}
