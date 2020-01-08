package runner

type IRunner interface {
	Call (script string) error
}
