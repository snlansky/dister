package runner

type LuaRunner struct {
}

func NewLuaRunner() IRunner {
	return &LuaRunner{}
}

func (r *LuaRunner) Call(script string) error {
	return nil
}
