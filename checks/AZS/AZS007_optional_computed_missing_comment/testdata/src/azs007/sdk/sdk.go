// Package sdk is a minimal stub of the azurerm typed resource SDK interfaces.
package sdk

import "azs007/schema"

type ResourceFunc func() error

type Resource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema
	ResourceType() string
	ModelObject() interface{}
	IDValidationFunc() func(interface{}, string) ([]string, []error)
	Create() ResourceFunc
	Read() ResourceFunc
	Delete() ResourceFunc
}

type ResourceWithUpdate interface {
	Resource
	Update() ResourceFunc
}
