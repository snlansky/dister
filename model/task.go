package model

import "strings"

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
