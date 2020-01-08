package runner

import (
	"dister/model"
)

type IRunner interface {
	Call(*model.UnitTest) (string, error)
}
