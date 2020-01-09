package runner

import "dister/protos"

type IRunner interface {
	Call(task *protos.Task) (string, error)
}
