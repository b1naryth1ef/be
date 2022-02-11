package ecs

import (
	"log"
	"reflect"
)

type componentMap = map[reflect.Type]interface{}

/// EntitySimpleStorage stores entities within a id-keyed map.
type EntitySimpleStorage struct {
	id   EntityId
	data map[EntityId]componentMap
}

func NewEntitySimpleStorage() *EntitySimpleStorage {
	return &EntitySimpleStorage{
		data: map[EntityId]componentMap{},
	}
}

func (e *EntitySimpleStorage) Add(id EntityId, components ...interface{}) {
	if _, exists := e.data[id]; exists {
		log.Panicf("duplicate entity was added to EntitySimpleStorage: %v", id)
	}
	e.data[id] = make(componentMap)
	for _, component := range components {
		e.data[id][reflect.TypeOf(component)] = component
	}
}

func (e *EntitySimpleStorage) Delete(id EntityId) {
	delete(e.data, id)
}

func (e *EntitySimpleStorage) FindAll(componentTypes []reflect.Type) []EntityId {
	result := []EntityId{}
	for entityId, components := range e.data {
		matches := true
		for _, componentType := range componentTypes {
			if _, exists := components[componentType]; !exists {
				matches = false
				continue
			}
		}
		if !matches {
			continue
		}
		result = append(result, entityId)
	}
	return result
}

func (e *EntitySimpleStorage) GetComponent(id EntityId, componentType reflect.Type) interface{} {
	return e.data[id][componentType]
}

func (e *EntitySimpleStorage) Get(id EntityId) []interface{} {
	components := e.data[id]
	result := make([]interface{}, len(components))
	idx := 0
	for _, component := range components {
		result[idx] = component
		idx++
	}
	return result
}

func (e *EntitySimpleStorage) GetComponentMap(id EntityId) componentMap {
	return e.data[id]
}

func (e *EntitySimpleStorage) RemoveComponent(id EntityId, componentType reflect.Type) {
	delete(e.data[id], componentType)
}

func (e *EntitySimpleStorage) AddComponent(id EntityId, component interface{}) {
	e.data[id][reflect.TypeOf(component)] = component
}
