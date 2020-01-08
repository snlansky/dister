package runner

type IRunner interface {
	Call (baseUrl, script string) (string, error)
}
