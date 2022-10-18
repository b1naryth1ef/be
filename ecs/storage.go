package ecs

import (
	"reflect"
)

type EntityId = uint32

var entityIdType reflect.Type

func init() {
	var empty EntityId
	entityIdType = reflect.TypeOf(empty)
}

// / EntityStorage manages storing, mutating, and querying a set of entities, described
// /  by their components.
type EntityStorage interface {
	Add(EntityId, ...interface{})
	Delete(EntityId)
	Get(EntityId) []interface{}
	GetComponent(EntityId, reflect.Type) interface{}
	RemoveComponent(EntityId, reflect.Type)
	AddComponent(EntityId, interface{})
	FindOne(reflect.Type) (interface{}, bool)
	FindAll([]reflect.Type, []reflect.Type) []EntityId
}

type EntityIterator interface {
	Next() bool
	Current() interface{}
}
