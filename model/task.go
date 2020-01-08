package model

import "strings"

type UnitTest struct {
	BaseUrl     string
	Script      string
	Validator   string
	ExceptValue string
}

func (ut *UnitTest) GetValidator() Validator {
	switch ut.Validator {
	case EqualValidatorType:
		return &EqualValidate{Value: ut.ExceptValue}
	}
	return nil
}

const (
	EqualValidatorType = "eq"
)

type Validator interface {
	Valid(v string) bool
}

type EqualValidate struct {
	Value string
}

func (eq *EqualValidate) Valid(v string) bool {
	return eq.Value == v
}

type InValidateString struct {
	Value string
}

func (i *InValidateString) Valid(v string) bool {
	return strings.Contains(i.Value, v)
}
