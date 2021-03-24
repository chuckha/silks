package core

import (
	"github.com/pkg/errors"
)

const (
	tableNameAttribute = "tablename"
)

type attribute struct {
	key, value string
}

func NewAttribute(key, value string) (*attribute, error) {
	switch key {
	case tableNameAttribute:
	default:
		return nil, errors.Errorf("%q is not a known attribute", key)
	}
	return &attribute{key, value}, nil
}

// attribute are custom comments that affect the output of the system in the models file
type Attributes map[string]string

func (a Attributes) AddAttribute(attr *attribute) {
	a[attr.key] = attr.value
}

func NewAttributes() Attributes {
	return make(map[string]string)
}
